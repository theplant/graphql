package graphql

import "fmt"

type QuerySchema struct {
	QueryType *Type
	Types     map[string]*Type
}

type Type struct {
	Name   string
	Fields map[string]Field
}

type Field struct {
	Type      *Type
	Arguments TypeArguments
}

type TypeArguments map[string]*Type

// type query {
// 	user(id: String): user
// }
//
// type user {
// 	id: String
// 	name: String
// 	profilePicture(size: Int): ProfilePicture
// }
//
// type profilePicture {
// 	uri: String
// 	width: Int
//     height: Int
// }

var string_, int_ *Type

func (t *Type) Scalar() bool {
	return len(t.Fields) == 0
}

func validateArgs(fieldName string, fieldArgs Arguments, field Field) error {
	// validate all type's args are present
	for argName, argType := range field.Arguments {
		queryArg, ok := fieldArgs[argName]
		if !ok {
			return fmt.Errorf("Missing argument %s for %s", argName, fieldName)
		} else if !queryArg.OfType(argType) {
			return fmt.Errorf("Wrong type %v for argument `%s` on %s, expected %s", queryArg, argName, fieldName, argType.Name)
		}
	}

	for argName, _ := range fieldArgs {
		if _, ok := field.Arguments[argName]; !ok {
			return fmt.Errorf("Unknown argument %s for %s", argName, fieldName)
		}
	}
	return nil
}

func validate(q Query, type_ *Type, schema QuerySchema) (err error) {
	if !type_.Scalar() && len(q.Fields) == 0 {
		return fmt.Errorf("Requested no fields from non-scalar type %s", q.Name)
	}

	for _, queryField := range q.Fields {
		if typeField, ok := type_.Fields[queryField.Name]; ok {
			if err = validateArgs(queryField.Name, queryField.Arguments, typeField); err != nil {
				return
			} else if err = validate(queryField, typeField.Type, schema); err != nil {
				return
			}
		} else {
			return fmt.Errorf("Unknown field %s on %s\n", queryField.Name, q.Name)
		}
	}
	return
}

func (schema QuerySchema) Validate(query Query) error {
	queryOrMutation := schema.Types[query.Name]
	if queryOrMutation == nil {
		return fmt.Errorf("unknown query entry point `%s`", query.Name)
	}
	return validate(query, queryOrMutation, schema)
}

func Schema() QuerySchema {

	string_ = &Type{Name: "String"}
	int_ = &Type{Name: "Int"}

	profilePicture := Type{
		Fields: map[string]Field{
			"Uri":    Field{Type: string_},
			"Width":  Field{Type: int_},
			"Height": Field{Type: int_},
		},
	}

	user := Type{
		Fields: map[string]Field{
			"Id":   Field{Type: string_},
			"Name": Field{Type: string_},
			"ProfilePicture": Field{
				Type: &profilePicture,
				Arguments: TypeArguments{
					"size": int_,
				},
			},
		},
	}

	query := Type{
		Fields: map[string]Field{
			"User": Field{
				Type: &user,
				Arguments: TypeArguments{
					"id": string_,
				},
			},
		},
	}

	return QuerySchema{
		Types: map[string]*Type{
			"query": &query,
			//			"User":           &user,
			//			"ProfilePicture": &profilePicture,
		},
	}
}
