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

type boolR bool

func (i boolR) unwrap() interface{} {
	return i
}

type mapR map[string]Result

func (m mapR) unwrap() interface{} {
	return m
}

// Fixme mapped by order, not by name :(
func values(args map[string]Result) (v []reflect.Value) {
	v = []reflect.Value{}
	for _, val := range args {
		v = append(v, reflect.ValueOf(val.unwrap()))
	}
	return
}

func Transform(query Query, context QueryContext) (Result, error) {
	v := reflect.ValueOf(context)
	fn := v.MethodByName(query.Name)
	if !fn.IsValid() {
		return nil, errors.New("Unknown query")
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
		data := mapR{}
		for _, field := range query.Fields {
			val, err := Transform(field, r.Interface())
			if err != nil {
				return nil, err
			}
			data[field.Name] = val
		}
		return data, nil
	} else {
		val := r.Interface()
		switch val.(type) {
		case string:
			return String(val.(string)), nil
		case int:
			return Int(val.(int)), nil
		case bool:
			return boolR(val.(bool)), nil
		default:
			return nil, errors.New(fmt.Sprintf("Unknown base type: %s (%s:%s) for %s", val, r.Type(), r, query.Name))
		}
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
