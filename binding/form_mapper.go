package binding

import (
	"errors"
	"net/url"
	"reflect"
	"strconv"
)

type mapper interface {
	mapping(reflect.Value, reflect.StructField, string) (bool, error)
	bind(interface{}) (bool, error)
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

func (m formMapper) bind(ptr interface{}) (bool, error) {
	return m.mapping(reflect.ValueOf(ptr), emptyField, "form")
}

func (m formMapper) mapping(v reflect.Value, field reflect.StructField, tag string) (bool, error) {
	if v.Kind() == reflect.Ptr {
		vptr := v
		var isNew bool
		if v.IsNil() {
			isNew = true
			vptr = reflect.New(v.Type().Elem())
		}
		ok, err := m.mapping(vptr.Elem(), field, tag)

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
			m.mapping(v.Field(i), sf, tag)
		}
	}

	// native type
	if v.Kind() != reflect.Struct || !field.Anonymous {
		tryToSet(v, field, tag, m)
	}

	return false, nil
}

func tryToSet(v reflect.Value, field reflect.StructField, tag string, m formMapper) {
	fieldName := field.Tag.Get(tag)
	vs, ok := m.data[fieldName]

	if !ok {
		return
	}

	switch v.Kind() {
	case reflect.Array:
		if v.Len() < len(vs) {
			panic("")
		}
		setAry(vs, v, field)
	case reflect.Slice:
		setSlice(vs, v, field)
	default:
		mapValueType(v, vs[0])
	}
}

func mapValueType(field reflect.Value, val string) error {
	switch field.Kind() {
	case reflect.Int:
		return setInt(val, 0, field)
	case reflect.Int8:
		return setInt(val, 8, field)
	case reflect.Int16:
		return setInt(val, 16, field)
	case reflect.Int32:
		return setInt(val, 32, field)
	case reflect.Int64:
		return setInt(val, 64, field)
	case reflect.Uint:
		return setUInt(val, 0, field)
	case reflect.Uint8:
		return setUInt(val, 8, field)
	case reflect.Uint16:
		return setUInt(val, 16, field)
	case reflect.Uint32:
		return setUInt(val, 32, field)
	case reflect.Uint64:
		return setUInt(val, 64, field)
	case reflect.String:
		return setString(val, field)
	case reflect.Float32:
		return setFloat(val, 32, field)
	case reflect.Float64:
		return setFloat(val, 64, field)
	case reflect.Bool:
		return setBool(val, field)
	default:
		return errors.New("Unknown type")
	}
	return nil
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

func setAry(xs []string, v reflect.Value, field reflect.StructField) error {
	for i, x := range xs {
		err := mapValueType(v.Index(i), x)
		if err != nil {
			return err
		}
	}
	return nil
}

func setSlice(xs []string, v reflect.Value, field reflect.StructField) error {
	s := reflect.MakeSlice(v.Type(), len(xs), len(xs))
	err := setAry(xs, s, field)
	if err == nil {
		v.Set(s)
	}
	return err
}

func setFloat(s string, bitSize int, field reflect.Value) error {
	f, err := strconv.ParseFloat(s, bitSize)
	if err == nil {
		field.SetFloat(f)
	}
	return err
}

func setBool(s string, field reflect.Value) error {
	b, err := strconv.ParseBool(s)
	if err == nil {
		field.SetBool(b)
	}
	return err
}
