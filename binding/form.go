package binding

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

type formParser struct {
	BodyParser
}

func (p formParser) Bind(r *http.Request, obj interface{}) error {
	if r == nil || r.Body == nil {
		return fmt.Errorf("invalid request")
	}

	return p.decode(r.Body, obj)
}
func (p formParser) BindBody(data []byte, obj interface{}) error {
	return p.decode(bytes.NewReader(data), obj)
}

func (formParser) decode(r io.Reader, obj interface{}) error {
	return nil
}
