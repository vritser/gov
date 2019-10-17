package gov

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Param struct {
	Key   string
	Value string
}

type Params []Param

type Context struct {
	Request    *http.Request
	Response   http.ResponseWriter
	params     Params
	queryCache url.Values
	formCache  url.Values
	Storage    map[string]interface{}
}

func (c *Context) reset() {
	c.params = c.params[0:0]
	c.Storage = nil
	c.formCache = nil
	c.queryCache = nil
}

func (c *Context) resetWriter(w http.ResponseWriter) {
	c.Response = w
}

func (c *Context) Path() string {
	return c.Request.URL.Path
}

func (c *Context) Method() string {
	return c.Request.Method
}

func (c *Context) Json(resp_body interface{}) {
	writeContentType(c.Response, []string{"application/json; charset=utf-8"})
	r, err := json.Marshal(resp_body)
	if err != nil {
		panic(err)
	}
	fmt.Fprintln(c.Response, string(r))
}

func (c *Context) String(resp_body string) {
	writeContentType(c.Response, []string{"text/plain; charset=utf-8"})
	fmt.Fprintln(c.Response, resp_body)
}

func (c *Context) Status(code int) {
	c.Response.WriteHeader(code)
}

func writeContentType(w http.ResponseWriter, value []string) {
	header := w.Header()

	if value := header["Content-Type"]; len(value) == 0 {
		header["Content-Type"] = value
	}
}

// context storage methods

func (c *Context) Set(key string, value interface{}) {
	if c.Storage == nil {
		c.Storage = make(map[string]interface{})
	}

	c.Storage[key] = value
}

func (c *Context) Get(key string) (value interface{}, exists bool) {
	value, exists = c.Storage[key]
	return
}

func (c *Context) MustGet(key string) interface{} {
	if value, exists := c.Get(key); exists {
		return value
	}

	panic("Key \"" + key + "\" does not exists")
}

func (c *Context) GetString(key string) (s string) {
	if val, ok := c.Get(key); ok && val != nil {
		s, _ = val.(string)
	}
	return
}

func (c *Context) GetBool(key string) (b bool) {
	if val, ok := c.Get(key); ok && val != nil {
		b, _ = val.(bool)
	}
	return
}

func (c *Context) GetInt(key string) (i int) {
	if val, ok := c.Get(key); ok && val != nil {
		i, _ = val.(int)
	}
	return

}

func (c *Context) GetInt64(key string) (i64 int64) {
	if val, ok := c.Get(key); ok && val != nil {
		i64, _ = val.(int64)
	}
	return
}

func (c *Context) GetFloat64(key string) (f64 float64) {
	if val, ok := c.Get(key); ok && val != nil {
		f64, _ = val.(float64)
	}
	return
}

func (c *Context) GetTime(key string) (t time.Time) {
	if val, ok := c.Get(key); ok && val != nil {
		t, _ = val.(time.Time)
	}
	return
}

func (c *Context) GetDuration(key string) (d time.Duration) {
	if val, ok := c.Get(key); ok && val != nil {
		d, _ = val.(time.Duration)
	}
	return
}

func (c *Context) GetStringSlice(key string) (ss []string) {
	if val, ok := c.Get(key); ok && val != nil {
		ss, _ = val.([]string)
	}
	return
}

func (c *Context) GetStringMap(key string) (sm map[string]interface{}) {
	if val, ok := c.Get(key); ok && val != nil {
		sm, _ = val.(map[string]interface{})
	}
	return
}

// Parse params、query、form

func (c *Context) Param(key string) interface{} {
	for _, p := range c.params {
		if p.Key == key {
			return p.Value
		}
	}

	return nil
}

func (c *Context) QueryOr(key string, defaultValue string) string {
	if value, ok := c.GetQuery(key); ok {
		return value
	}

	return defaultValue
}

func (c *Context) Query(key string) string {
	value, _ := c.GetQuery(key)
	return value
}

func (c *Context) GetQuery(key string) (string, bool) {
	if values, ok := c.GetQueryArray(key); ok {
		return values[0], ok
	}
	return "", false
}

func (c *Context) QueryArray(key string) []string {
	values, _ := c.GetQueryArray(key)
	return values
}

func (c *Context) QueryMap(key string) map[string]string {
	m, _ := c.GetQueryMap(key)
	return m
}

func (c *Context) GetQueryMap(key string) (map[string]string, bool) {
	c.getQueryCache()
	return c.get(c.queryCache, key)
}

func (c *Context) GetQueryArray(key string) ([]string, bool) {
	c.getQueryCache()
	if values, ok := c.queryCache[key]; ok && len(values) > 0 {
		return values, true
	}
	return []string{}, false
}

func (c *Context) Form(key string) string {
	value, _ := c.GetForm(key)
	return value
}

func (c *Context) GetForm(key string) (string, bool) {
	if values, ok := c.GetFormArray(key); ok {
		return values[0], ok
	}
	return "", false
}

func (c *Context) GetFormMap(key string) (map[string]string, bool) {
	c.getFormCache()
	return c.get(c.formCache, key)
}

func (c *Context) FormAry(key string) []string {
	values, _ := c.GetFormArray(key)
	return values
}

func (c *Context) GetFormArray(key string) ([]string, bool) {
	c.getFormCache()
	if values, ok := c.formCache[key]; ok && len(values) > 0 {
		return values, ok
	}
	return []string{}, false
}

func (c *Context) get(store url.Values, key string) (map[string]string, bool) {
	dicts := make(map[string]string)
	exist := false

	for k, v := range store {
		if i := strings.IndexByte(k, '['); i >= 1 && k[0:i] == key {
			if j := strings.IndexByte(k, ']'); j > i {
				dicts[k[i+1:j]] = v[0]
			}
		}
	}

	return dicts, exist
}

func (c *Context) getQueryCache() {
	if c.queryCache == nil {
		c.queryCache = c.Request.URL.Query()
	}
}

func (c *Context) getFormCache() {
	if c.formCache == nil {
		c.Request.ParseForm()
		c.formCache = c.Request.PostForm
	}
}
