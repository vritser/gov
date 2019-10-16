package gov

import (
	"net/http"
	"strings"
)

type Gov struct {
	Router
}

func New() *Gov {
	return &Gov{
		Router: Router{
			trees:       make(methodTrees, 0, 9),
			Middlewares: make(HandlerChain, 0),
		},
	}
}

func (v *Gov) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := &Context{
		Request:  r,
		Response: w,
	}

	v.handleHTTPRequest(c)
}

func (v *Gov) Run(s string) error {
	err := http.ListenAndServe(s, v)
	return err
}

func (v *Gov) handleHTTPRequest(c *Context) {

	for _, m := range v.Router.Middlewares {
		m(c)
	}

	path := c.Path()
	m := c.Method()

	n := v.handle(m, path)

	if n.handlers == nil {
		// fmt.Println("not found the route")
		return
	}

	c.params = n.params

	for _, h := range n.handlers {
		h(c)
		c.GetQueryArray("xs")
		c.GetQueryMap("ids")
	}
}

func (v *Gov) handle(m, path string) (ret nodeValue) {
	r := v.Router.trees.get(m)

	if len(path) > len(r.path) {
		// not this node

		p := splitBy('/', path)

		children := r.children
	GOV:
		for {
			for _, child := range children {
				if child.path == p[0] {
					p = p[1:]
					if len(p) > 0 {
						children = child.children
						ret.fullPath += "/" + child.path
						continue GOV
					}
					ret.handlers = child.handlers
					return
				}

				if child.nType == param {
					ret.params = append(ret.params, Param{
						Key:   child.name,
						Value: p[0],
					})
					p = p[1:]
					if len(p) > 0 {
						ret.fullPath += "/" + child.path
						children = child.children
						continue GOV
					}

					ret.handlers = child.handlers
					return
				}

			}

			return
		}

	}
	ret.handlers = r.handlers
	return
}

func (v *Gov) handleRequest(method, path string) *node {
	root := v.Router.trees.get(method)

	if len(path) == 1 {
		return root
	}

	// segs := splitBy('/', path)
	segs := strings.Split(path, "/")

	if root.path == segs[0] && len(segs[1:]) == 0 {
		return root
	} else {
		return matchUrl(root.children, segs[1:])
	}
}

func matchUrl(ns []*node, path []string) *node {
	for _, n := range ns {
		if n.path == path[0] {
			if len(path[1:]) == 0 {
				return n
			}

			return matchUrl(n.children, path[1:])
		}
	}

	return nil
}
func splitBy(sperator rune, str string) []string {
	return strings.FieldsFunc(str, func(r rune) bool {
		return r == sperator
	})
}
