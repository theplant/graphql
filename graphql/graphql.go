package graphql

import (
	"errors"
	"fmt"
	"reflect"
)

type QueryContext interface{}

type Result interface {
	unwrap() interface{}
	OfType(*Type) bool
}

type Int int

func (i Int) unwrap() interface{} {
	return int(i)
}

func (i Int) OfType(type_ *Type) bool {
	return type_ == int_
}

type String string

func (s String) unwrap() interface{} {
	return string(s)
}

func (s String) OfType(type_ *Type) bool {
	return type_ == string_
}

type resultMap map[string]Result

func (m resultMap) unwrap() interface{} {
	return m
}

func (m resultMap) OfType(type_ *Type) bool {
	return false
}

type resultArray struct {
	arr []Result
}

func (a resultArray) unwrap() interface{} {
	return a.arr
}

func (a resultArray) OfType(type_ *Type) bool {
	return false
}

// Fixme mapped by order, not by name :(
func values(args map[string]Result) (v []reflect.Value) {
	v = []reflect.Value{}
	for _, val := range args {
		v = append(v, reflect.ValueOf(val.unwrap()))
	}
	return
}

func transformArray(query Query, context QueryContext) (result Result, err error) {
	v := reflect.ValueOf(context)
	r := []Result{}
	for i := 0; i < v.Len(); i++ {
		item := v.Index(i)
		rx, e := transformValue(query, item)
		if e != nil {
			return nil, e
		}
		r = append(r, rx)
	}
	return resultMap{query.Name: resultArray{arr: r}}, nil
}

func transformScalar(value reflect.Value) (Result, error) {
	val := value.Interface()
	switch val.(type) {
	case string:
		return String(val.(string)), nil
	case int:
		return Int(val.(int)), nil
	default:
		return nil, errors.New(fmt.Sprintf("Unknown base type: %s", value))
	}
}

func transformValue(query Query, value reflect.Value) (Result, error) {
	data := resultMap{}
	for _, field := range query.Fields {
		val, err := Transform(field, value.Interface())
		if err != nil {
			return nil, err
		}
		name := field.Alias
		if len(name) == 0 {
			name = field.Name
		}
		data[name] = val
	}
	return data, nil
}

func Transform(query Query, context QueryContext) (Result, error) {
	fn := reflect.ValueOf(context).MethodByName(query.Name)

	if !fn.IsValid() {
		return nil, errors.New(fmt.Sprintf("Unknown query '%s' on '%v'", query.Name, context))
	}

	values := values(query.Arguments)
	for i, val := range values {
		if fn.Type().In(i) != reflect.TypeOf(val.Interface()) {
			return nil, errors.New("Invalid type for fn")
		}
	}

	// [0] erk!
	r := fn.Call(values)[0]

	if len(query.Fields) > 0 {
		switch r.Kind() {
		case reflect.Slice:
			return transformArray(query, r.Interface())
		default:
			return transformValue(query, r)
		}
	} else {
		return transformScalar(r)
	}
}

type Arguments map[string]Result
type Fields []Query

type Query struct {
	Name      string
	Alias     string
	Arguments Arguments
	Fields    []Query
}

func Printq(query Query) {
	{
		fmt.Print(query.Name)
		if len(query.Arguments) > 0 {
			fmt.Print("(")
			for k, v := range query.Arguments {
				fmt.Printf("%s: %s", k, v)
			}
			fmt.Print(")")
		}
		if len(query.Fields) > 0 {
			fmt.Println(" {")
			for _, sub := range query.Fields {
				Printq(sub)
			}
			fmt.Print("},")
		}
		fmt.Println()
	}
}
