package binding

import (
	"mime/multipart"
	"net/url"
	"reflect"
)

type multipartFormMapper struct {
	mapper
	data  url.Values
	files map[string][]*multipart.FileHeader
}

func (m multipartFormMapper) bind(obj interface{}) (bool, error) {
	return m.mapping(reflect.ValueOf(obj), emptyField, "form")
}

func (m multipartFormMapper) mapping(v reflect.Value, field reflect.StructField, tag string) (bool, error) {
	name := field.Tag.Get(tag)

	if files, _ := m.files[name]; len(files) <= 0 {
		formapper.data = m.data
		return formapper.mapping(v, field, tag)
	}
	return false, nil
}
