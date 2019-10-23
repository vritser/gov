package binding

import (
	"io"
	"net/http"
)

type formParser struct {
	BodyParser
}
type multipartFormParser struct {
	BodyParser
}
type queryParser struct {
	BodyParser
}

var (
	formapper = formMapper{}
	mulMapper = multipartFormMapper{}
	qmapper   = queryMapper{}
)

func (p formParser) Bind(r *http.Request, obj interface{}) (err error) {
	if err = r.ParseForm(); err != nil {
		return
	}

	formapper.data = r.PostForm
	_, err = formapper.bind(obj)
	return
}
func (p formParser) BindBody(data []byte, obj interface{}) error {
	return nil
}

func (formParser) decode(r io.Reader, obj interface{}) error {
	return nil
}

func (p multipartFormParser) Bind(r *http.Request, obj interface{}) (err error) {
	if err = r.ParseMultipartForm(256); err != nil {
		return
	}

	mulMapper.data = r.PostForm
	mulMapper.files = r.MultipartForm.File
	_, err = mulMapper.bind(obj)

	return err
}

func (p queryParser) Bind(r *http.Request, obj interface{}) error {
	qmapper.data = r.URL.Query()
	_, err := qmapper.bind(obj)
	return err
}
