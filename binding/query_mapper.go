package binding

import (
	"net/url"
	"reflect"
)

type queryMapper struct {
	mapper
	data url.Values
}

func (m queryMapper) bind(obj interface{}) (bool, error) {
	return m.mapping(reflect.ValueOf(obj), emptyField, "query")
}

func (m queryMapper) mapping(v reflect.Value, field reflect.StructField, tag string) (bool, error) {
	formapper.data = m.data
	return formapper.mapping(v, field, tag)
}
