package binding

import (
	"bytes"
	"io"
	"net/http"
)

type formParser struct {
	BodyParser
}

var formapper = formMapper{}

func (p formParser) Bind(r *http.Request, obj interface{}) error {
	if err := r.ParseForm(); err != nil {
		return err
	}

	formapper.data = r.PostForm
	_, err := formMapping(obj, formapper)
	return err
}
func (p formParser) BindBody(data []byte, obj interface{}) error {
	return p.decode(bytes.NewReader(data), obj)
}

func (formParser) decode(r io.Reader, obj interface{}) error {
	return nil
}
