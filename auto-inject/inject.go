package auto_inject

import (
	"errors"
	"fmt"
	"reflect"
	"unicode"
)

var (
	ErrCantSet   = errors.New("cannot set")
	ErrNotStruct = errors.New("not struct")
)

type ErrConflict struct {
	A, B       string
	DomainName string
	T          reflect.Type
}

func (e ErrConflict) Error() string {
	return fmt.Sprintf(`<Field <%s,%s>,Conflict <"%s",%v>>`, e.A, e.B, e.DomainName, e.T)
}

// Injector can bind dependencies to target injectTarget
type Injector interface {
	// Bind to target
	Bind(injectTarget interface{}) (notFoundList []string, err error)
}

// Source struct describes a dependency of targets collected from injector
type Source struct {
	S reflect.StructField
	V reflect.Value
}

// FlatSource records the result of parsed data, and implements Injector
type FlatSource map[reflect.Type]map[string]Source

// Bind to target
func (f FlatSource) Bind(i interface{}) ([]string, error) {
	v, err := getElement(i)
	if err != nil {
		return nil, err
	}
	t := v.Type()
	var notFoundList []string
	for i := 0; i < v.NumField(); i++ {
		field, fieldType := v.Field(i), t.Field(i)
		if !unicode.IsUpper(rune(fieldType.Name[0])) {
			continue
		}

		name := fieldType.Tag.Get("injector")
		if !field.CanSet() {
			return nil, ErrCantSet
		}
		tSet := f[fieldType.Type]
		if len(name) == 0 {
			if s, ok := tSet[fieldType.Name]; ok {
				field.Set(s.V)
			} else if s, ok := tSet[""]; ok {
				field.Set(s.V)
			} else {
				notFoundList = append(notFoundList, fieldType.Name)
			}
		} else {
			if s, ok := tSet[name]; ok {
				field.Set(s.V)
			} else {
				notFoundList = append(notFoundList, fieldType.Name)
			}
		}
	}
	return notFoundList, nil
}

type AnyStruct interface{}
// ParseFlatSource can parse any struct and return a dependency-injector
func ParseFlatSource(source AnyStruct) (FlatSource, error) {
	v, err := getElement(source)
	if err != nil {
		return nil, err
	}
	s := make(FlatSource)
	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field, fieldType := v.Field(i), t.Field(i)
		name := fieldType.Tag.Get("injector")

		if s[fieldType.Type] == nil {
			s[fieldType.Type] = make(map[string]Source)
		}
		x := Source{
			S: fieldType,
			V: field,
		}

		if oldField, ok := s[fieldType.Type][name]; ok {
			return nil, ErrConflict{oldField.S.Name, fieldType.Name, name, fieldType.Type}
		}
		s[fieldType.Type][name] = x

		if oldField, ok := s[fieldType.Type][fieldType.Name]; ok {
			return nil, ErrConflict{oldField.S.Name, fieldType.Name, fieldType.Name, fieldType.Type}
		}
		s[fieldType.Type][fieldType.Name] = x

		if len(name) != 0 {
			if _, ok := s[fieldType.Type][""]; !ok {
				s[fieldType.Type][""] = x
			}
		}
	}

	return s, nil
}

func getElement(i interface{}) (reflect.Value, error) {
	v := reflect.ValueOf(i)
	for v.Type().Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Type().Kind() != reflect.Struct {
		return reflect.Value{}, ErrNotStruct
	}
	return v, nil
}
