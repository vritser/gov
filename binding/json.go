package binding

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type jsonParser struct {
	BodyParser
}

func (p jsonParser) Bind(r *http.Request, obj interface{}) error {
	if r == nil || r.Body == nil {
		return fmt.Errorf("invalid request")
	}

	return p.decode(r.Body, obj)
}
func (p jsonParser) BindBody(data []byte, obj interface{}) error {
	return p.decode(bytes.NewReader(data), obj)
}

func (jsonParser) decode(r io.Reader, obj interface{}) error {
	decoder := json.NewDecoder(r)

	if err := decoder.Decode(obj); err != nil {
		return err
	}

	return nil
}
