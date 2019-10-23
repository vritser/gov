package binding

import (
	"io"
	"net/http"
)

const (
	MIMEJSON              = "application/json"
	MIMEHTML              = "text/html"
	MIMEXML               = "application/xml"
	MIMEXML2              = "text/xml"
	MIMEPlain             = "text/plain"
	MIMEPOSTForm          = "application/x-www-form-urlencoded"
	MIMEMultipartPOSTForm = "multipart/form-data"
	MIMEPROTOBUF          = "application/x-protobuf"
	MIMEMSGPACK           = "application/x-msgpack"
	MIMEMSGPACK2          = "application/msgpack"
	MIMEYAML              = "application/x-yaml"
)

type BodyParser interface {
	Bind(*http.Request, interface{}) error
	BindBody([]byte, interface{}) error
	decode(io.Reader, interface{}) error
}

var (
	JSON          = jsonParser{}
	XML           = xmlParser{}
	Form          = formParser{}
	MultipartForm = multipartFormParser{}
	Query         = queryParser{}
)

func Default(method, contentType string) BodyParser {

	if method == "GET" {
		return Query
	}

	switch contentType {
	case MIMEJSON:
		return JSON
	case MIMEXML, MIMEXML2:
		return XML
	case MIMEMultipartPOSTForm:
		return MultipartForm
	default:
		return Form
	}
}
