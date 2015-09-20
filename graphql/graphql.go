package graphql

import (
	"errors"
	"fmt"
	"reflect"
)

type QueryContext interface{}

type Result interface {
	unwrap() interface{}
}

type Int int

func (i Int) unwrap() interface{} {
	return int(i)
}

type String string

func (s String) unwrap() interface{} {
	return string(s)
}

type resultMap map[string]Result

func (m resultMap) unwrap() interface{} {
	return m
}

type resultArray struct {
	arr []Result
}

func (a resultArray) unwrap() interface{} {
	return a.arr
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
		data[field.Name] = val
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
