package xml

import (
	"reflect"

	"github.com/clbanning/mxj"
	"github.com/pkg/errors"
)

func Marshal(v interface{}) ([]byte, error) {
	return mxj.AnyXml(v)
}

func Unmarshal(data []byte, v interface{}) error {
	rv := reflect.ValueOf(v)

	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return errors.Errorf("invalid value %s", rv.Type())
	}

	m, err := mxj.NewMapXml(data, true)
	if err != nil {
		return err
	}

	rv = indirect(rv)
	switch rv.Kind() {
	case reflect.Interface:
		if rv.NumMethod() > 0 {
			return errors.Errorf("unsupported type %s", rv.Type())
		}
	case reflect.Map:
		if _, ok := v.(*map[string]interface{}); !ok {
			return errors.Errorf("unsupported type %s", rv.Type())
		}
	default:
		return errors.Errorf("unsupported type %s", rv.Type())
	}

	rv.Set(reflect.ValueOf(m.Old()))
	return nil
}

// This was copied (and simplified) from github.com/BurntSushi/toml/decode.go
func indirect(v reflect.Value) reflect.Value {
	if v.Kind() != reflect.Ptr {
		return v
	}

	if v.IsNil() {
		v.Set(reflect.New(v.Type().Elem()))
	}

	return indirect(reflect.Indirect(v))
}
