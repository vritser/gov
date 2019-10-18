package binding

import (
	"net/url"
	"reflect"
	"strconv"
)

type mapper interface {
	source() interface{}
	mapping()
	get(key string) interface{}
}

type formMapper struct {
	mapper
	data url.Values
}

func (m formMapper) get(key string) interface{} {
	if values, ok := m.data[key]; ok {
		return values
	}
	return []string{}
}

var emptyField = reflect.StructField{}

func formMapping(ptr interface{}, data formMapper) (bool, error) {
	return parse(reflect.ValueOf(ptr), emptyField, data, "form")
}

func parse(v reflect.Value, field reflect.StructField, source mapper, tag string) (bool, error) {
	if v.Kind() == reflect.Ptr {
		vptr := v
		var isNew bool
		if v.IsNil() {
			isNew = true
			vptr = reflect.New(v.Type().Elem())
		}
		ok, err := parse(vptr.Elem(), field, source, tag)

		if err != nil {
			return false, err
		}

		if isNew && ok {
			v.Set(vptr)
		}

		return true, nil
	}

	if v.Kind() == reflect.Struct {
		st := v.Type()
		for i := 0; i < v.NumField(); i++ {
			sf := st.Field(i)
			if sf.PkgPath != "" || sf.Anonymous {
				continue
			}
			parse(v.Field(i), sf, source, tag)
		}
	}

	// native type
	if v.Kind() != reflect.Struct || !field.Anonymous {
		tryToSet(v, field, tag, source)
	}

	return false, nil
}

func tryToSet(v reflect.Value, field reflect.StructField, tag string, m mapper) {
	fieldName := field.Tag.Get(tag)

	if fieldName == "" {
		return
	}

	vs := m.get(fieldName).([]string)

	switch v.Kind() {
	case reflect.Array:
		if v.Len() < len(vs) {
			panic("")
		}
		setAry(vs, v, field)
	case reflect.Slice:

	default:
		mapValueType(v, vs[0])
	}
}

func mapValueType(field reflect.Value, val string) {
	switch field.Kind() {
	case reflect.Int:
		setInt(val, 0, field)
	case reflect.Int8:
		setInt(val, 8, field)
	case reflect.Int16:
		setInt(val, 16, field)
	case reflect.Int32:
		setInt(val, 32, field)
	case reflect.Int64:
		setInt(val, 64, field)
	case reflect.Uint:
		setUInt(val, 0, field)
	case reflect.Uint8:
		setUInt(val, 8, field)
	case reflect.Uint16:
		setUInt(val, 16, field)
	case reflect.Uint32:
		setUInt(val, 32, field)
	case reflect.Uint64:
		setUInt(val, 64, field)
	case reflect.String:
		setString(val, field)
	}
}
func setInt(s string, bitSize int, field reflect.Value) error {
	i, err := strconv.ParseInt(s, 10, bitSize)
	if err == nil {
		field.SetInt(i)
	}
	return err
}
func setUInt(s string, bitSize int, field reflect.Value) error {
	val, err := strconv.ParseUint(s, 10, bitSize)
	if err == nil {
		field.SetUint(val)
	}
	return err
}

func setString(s string, field reflect.Value) error {
	field.SetString(s)
	return nil
}

func setAry(xs []string, v reflect.Value, field reflect.StructField) {
	for _, x := range xs {
		mapValueType(v, x)
	}
}
