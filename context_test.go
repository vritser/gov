package gov

import (
	"bytes"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestContextStorage(t *testing.T) {
	c, _ := CrateTestCtx(nil)

	c.Set("foo", "bar")

	value, ok := c.Get("foo")
	assert.True(t, ok)
	assert.Equal(t, "bar", value)

	value, ok = c.Get("bar")
	assert.False(t, ok)
	assert.Nil(t, value)
}

func TestContextSetAndGetTypedValues(t *testing.T) {
	c, _ := CrateTestCtx(nil)

	c.Set("f64", float64(3.14159))
	c.Set("i64", int64(40))
	c.Set("bool", true)
	c.Set("string", "string")
	c.Set("int", int(235))
	date, _ := time.Parse("1/2/2006 15:04:05", "01/01/2017 12:00:00")
	c.Set("time", date)
	c.Set("duration", time.Second)
	c.Set("string_slice", []string{"hello", "world"})
	c.Set("string_map", map[string]interface{}{
		"key": "value",
	})

	assert.Equal(t, float64(3.14159), c.GetFloat64("f64"))
	assert.Equal(t, int64(40), c.GetInt64("i64"))
	assert.Equal(t, true, c.GetBool("bool"))
	assert.Equal(t, "string", c.GetString("string"))
	assert.Equal(t, 235, c.GetInt("int"))
	assert.Equal(t, date, c.GetTime("time"))
	assert.Equal(t, time.Second, c.GetDuration("duration"))
	assert.Equal(t, []string{"hello", "world"}, c.GetStringSlice("string_slice"))
	assert.Equal(t, map[string]interface{}{"key": "value"}, c.GetStringMap("string_map"))
}

func TestContextQuery(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := CrateTestCtx(w)

	c.Request, _ = http.NewRequest("GET", "http://example.com/?id=10&name=vritser&q=", nil)
	value, ok := c.GetQuery("id")
	assert.True(t, ok)
	assert.Equal(t, "10", value)
	assert.Equal(t, "10", c.QueryOr("id", "20"))
	assert.Equal(t, "10", c.Query("id"))

	name, ok := c.GetQuery("name")
	assert.True(t, ok)
	assert.Equal(t, "vritser", name)
	assert.Equal(t, "vritser", c.QueryOr("name", "alice"))
	assert.Equal(t, "vritser", c.Query("name"))

	q, ok := c.GetQuery("q")
	assert.True(t, ok)
	assert.Empty(t, q)
	assert.Empty(t, c.QueryOr("q", "something"))
	assert.Empty(t, c.Query("q"))

	value, ok = c.GetQuery("key")
	assert.False(t, ok)
	assert.Empty(t, value)
	assert.Equal(t, "value", c.QueryOr("key", "value"))
	assert.Empty(t, c.Query("key"))

	// post form should be empty
	value, ok = c.GetForm("id")
	assert.False(t, ok)
	assert.Empty(t, value)
	assert.Empty(t, c.Form("id"))
}

func TestContextPostForm(t *testing.T) {
	c, _ := CrateTestCtx(nil)

	body := bytes.NewBufferString("name=vritser&age=20")
	c.Request, _ = http.NewRequest("POST", "http://example.com/usr?id=123", body)
	c.Request.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	id, ok := c.GetQuery("id")
	assert.True(t, ok)
	assert.Equal(t, "123", id)

	name, ok := c.GetForm("name")
	assert.True(t, ok)
	assert.Equal(t, "vritser", name)

	age := c.Form("age")
	assert.Equal(t, "20", age)
}

func TestContextBindJSON(t *testing.T) {
	c, _ := CrateTestCtx(nil)

	body := bytes.NewBufferString("{\"id\": \"1024\", \"name\": \"vritser\", \"foo\": {\"bar\":\"bar\"}}")
	c.Request, _ = http.NewRequest("POST", "/", body)
	c.Request.Header.Add("Content-Type", "application/json")

	var obj struct {
		Id   string `json:"id"`
		Name string `json:"name"`
		Foo  Foo    `json:"foo"`
	}

	err := c.Bind(&obj)
	assert.Nil(t, err)
	assert.Equal(t, "vritser", obj.Name)
	assert.Equal(t, "1024", obj.Id)
	assert.Equal(t, "bar", obj.Foo.Bar)
}

func TestContextBindForm(t *testing.T) {
	c, _ := CrateTestCtx(nil)
	body := bytes.NewBufferString("name=vritser&id=20&age=10&ary=520&ary=521&slc=747&slc=749")
	c.Request, _ = http.NewRequest("POST", "http://example.com/?t=123", body)
	c.Request.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	var obj struct {
		Id   string `form:"id"`
		Name string `form:"name"`
		Age  int
		sex  int
		I32  int32  `form:"age"`
		Ary  [2]int `form:"ary"`
		Slc  []int  `form:"slc"`
		int  `form:"age"`
	}

	err := c.Bind(&obj)

	assert.Nil(t, err)
	assert.Equal(t, []int{747, 749}, obj.Slc)
	assert.Equal(t, [2]int{520, 521}, obj.Ary)
	assert.Equal(t, "20", obj.Id)
	assert.Equal(t, "vritser", obj.Name)
	assert.Equal(t, reflect.Int32, reflect.ValueOf(obj.I32).Kind())
	assert.Zero(t, obj.int)
	assert.Zero(t, obj.Age)
	assert.Zero(t, obj.sex)
}

func TestContextBindMultipartForm(t *testing.T) {
	c, _ := CrateTestCtx(nil)

	buf := new(bytes.Buffer)
	mw := multipart.NewWriter(buf)
	assert.NoError(t, mw.WriteField("id", "10"))
	assert.NoError(t, mw.WriteField("name", "vritser"))
	mw.Close()

	c.Request, _ = http.NewRequest("POST", "/", buf)
	c.Request.Header.Set("Content-Type", mw.FormDataContentType())

	var obj struct {
		Id   int    `form:"id"`
		Name string `form:"name"`
	}

	err := c.Bind(&obj)
	assert.Nil(t, err)
	assert.Equal(t, 10, obj.Id)
	assert.Equal(t, "vritser", obj.Name)
}

func TestContextWriteFile(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := CrateTestCtx(w)
	c.Request, _ = http.NewRequest("GET", "/", nil)

	c.WriteFile("./gov.go")

	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "func New() *Gov {")
	assert.Equal(t, "text/plain; charset=utf-8", w.Header().Get("Content-Type"))
}

func TestContextWriteAttachement(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := CrateTestCtx(w)

	c.Request, _ = http.NewRequest("GET", "/", nil)
	filename := "SpecifiedFilieName.go"
	c.WriteAttachement("./gov.go", filename)

	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "func New() *Gov {")
	assert.Equal(t, fmt.Sprintf("attachment; filename=\"%s\"", filename), w.Header().Get("Content-Disponsion"))
}

type Foo struct {
	Bar string
}
