package binding

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
)

type xmlParser struct {
	BodyParser
}

func (p xmlParser) Bind(r *http.Request, obj interface{}) error {
	if r == nil || r.Body == nil {
		return fmt.Errorf("invalid request")
	}

	return p.decode(r.Body, obj)
}
func (p xmlParser) BindBody(data []byte, obj interface{}) error {
	return p.decode(bytes.NewReader(data), obj)
}

func (xmlParser) decode(r io.Reader, obj interface{}) error {
	decoder := xml.NewDecoder(r)

	if err := decoder.Decode(obj); err != nil {
		return err
	}

	return nil
}
