package gov

import "net/http"

func CrateTestCtx(w http.ResponseWriter) (c *Context, r *Gov) {
	r = New()
	c = r.allocateContext()
	c.reset()
	c.resetWriter(w)

	return
}
