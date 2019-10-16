package gov

import (
	"strings"
)

type HandlerFunc func(*Context)

type HandlerChain []HandlerFunc

type IRouter interface {
	Use(...HandlerFunc)

	Get(string, ...HandlerFunc)
	Post(string, ...HandlerFunc)
	Put(string, ...HandlerFunc)
	Delete(string, ...HandlerFunc)
	Options(string, ...HandlerFunc)
	Patch(string, ...HandlerFunc)
	Any(string, ...HandlerFunc)
	Head(string, ...HandlerFunc)
}

type methodTree struct {
	method string
	root   *node
}

type nodeType uint8

const (
	static nodeType = iota
	root
	param
)

type node struct {
	path     string
	handlers HandlerChain
	children []*node
	fullpath string
	nType    nodeType
	name     string
}

type nodeValue struct {
	handlers HandlerChain
	params   Params
	fullPath string
}

type methodTrees []methodTree

type Router struct {
	trees       methodTrees
	Middlewares HandlerChain
}

type RouteInfo struct {
	Method  string
	Path    string
	Handler string
}

func (route *Router) Use(middlewares ...HandlerFunc) {
	route.Middlewares = append(route.Middlewares, middlewares...)
}

func insert2(root *node, m string, p []string, hs ...HandlerFunc) {
	for _, n := range root.children {

		if n.path == p[0] {

			if len(p[1:]) == 0 {
				n.handlers = hs
			} else {
				insert2(n, m, p[1:], hs...)
			}

			return
		}
	}

	var n *node
	if strings.HasPrefix(p[0], ":") {
		n = &node{
			path:     p[0],
			handlers: hs,
			children: []*node{},
			nType:    param,
			name:     p[0][1:],
		}
	} else {
		n = &node{
			path:     p[0],
			handlers: hs,
			children: []*node{},
		}
	}

	root.children = append(root.children, n)

	if len(p[1:]) > 0 {
		insert2(n, m, p[1:], hs...)
	}
}

func (r *Router) add(m, path string, hs ...HandlerFunc) {

	root := r.trees.get(m)

	if root == nil {
		root = &node{
			path:     "/",
			handlers: HandlerChain{},
			children: []*node{},
			fullpath: "/",
		}

		r.trees = append(r.trees, methodTree{method: m, root: root})
	}

	if root.path == path {
		root.handlers = hs
	} else {
		segs := strings.Split(path, "/")

		insert2(root, m, segs[1:], hs...)
	}
}

func (r *Router) Get(path string, handlers ...HandlerFunc) {
	r.add("GET", path, handlers...)
}

func (r *Router) Post(path string, handlers ...HandlerFunc) {
	r.add("POST", path, handlers...)
}

func (r *Router) Put(path string, handlers ...HandlerFunc) {
	r.add("PUT", path, handlers...)
}

func (r *Router) Delete(path string, handlers ...HandlerFunc) {
	r.add("DELETE", path, handlers...)
}

func (trees methodTrees) get(method string) *node {
	for _, tree := range trees {
		if tree.method == method {
			return tree.root
		}
	}
	return nil
}

func (r *Router) Routes() []RouteInfo {
	ret := []RouteInfo{}
	for _, t := range r.trees {
		ret = append(ret, iterate("", t.method, t.root)...)
	}

	return ret

}

func iterate(path, method string, root *node) []RouteInfo {

	ret := []RouteInfo{}
	path += root.path

	ret = append(ret, RouteInfo{
		Method:  method,
		Path:    path,
		Handler: string(len(root.handlers)),
	})

	for _, child := range root.children {
		ret = append(ret, iterate(path, method, child)...)
	}

	return ret
}
