package graphql

import (
	"errors"
	"fmt"
	"strconv"

	parsec "github.com/prataprc/goparsec"
)

// {
//   user(id: 3500401) {
//     id,
//     name,
//     isViewerFriend,
//     pic: profilePicture(size: 50)  {
//       uri,
//       width,
//       height
//     }
//   }
// }

var y, selectionSet, name parsec.Parser

func Parse(input string) (query []Query, err error) {
	s := parsec.NewScanner([]byte(input))
	result, next := y(s)
	if next.Endof() {
		query = result.([]Query)
	} else {
		err = errors.New("Failed to parse input")
	}

	return
}

func init() {
	name = parsec.Maybe(identFn, parsec.Ident())

	// missing `fragment spread`, `inline fragment`
	selection := parsec.OrdChoice(extract(0), field())
	selectionSet = parsec.And(extract(1), parsec.Token(`{`, ""), parsec.Many(fieldArrayify, selection), parsec.Token(`}`, ""))

	// missing `fragment definition`, `type definition`
	definition := operationDefinition()
	document := parsec.OrdChoice(extract(0), parsec.Maybe(shorthandQueryify, selectionSet), parsec.Many(fieldArrayify, definition))

	y = document
}

func operationDefinition() parsec.Parser {
	operationType := parsec.OrdChoice(extract(0), parsec.Token(`query`, "OP_TYPE_QUERY"), parsec.Token(`mutation`, "OP_TYPE_MUTATION"))

	// missing `variable definitions`, `directives`
	return sequence(operationDefinitionify, &operationType, &name, &selectionSet)
}

func identFn(nodes []parsec.ParsecNode) parsec.ParsecNode {
	return nodes[0].(*parsec.Terminal).Value
}

// operationType name selectionSet -> query
func operationDefinitionify(nodes []parsec.ParsecNode) parsec.ParsecNode {
	if nodes[0] == nil || nodes[2] == nil {
		return nil
	}

	query := Query{
		Name:   nodes[0].(*parsec.Terminal).Value,
		Fields: nodes[2].([]Query),
	}
	if nodes[1] != nil {
		query.Alias = nodes[1].(*parsec.Terminal).Value
	}
	return query
}

// `{ fields }` -> []query
func shorthandQueryify(nodes []parsec.ParsecNode) parsec.ParsecNode {
	return []Query{
		Query{
			Name:   "query",
			Fields: nodes[0].([]Query),
		},
	}
}

func field() parsec.Parser {

	alias := parsec.And(identFn, parsec.Ident(), parsec.Token(`:`, ""))

	arguments := parsec.And(
		extract(1),
		parsec.Token(`\(`, ""),
		parsec.Many(
			argify,
			// `Argument` list
			parsec.And(nil, parsec.Maybe(identFn, parsec.Ident()), parsec.Token(`:`, ""), value())),
		parsec.Token(`\)`, ""))

	// Missing `directives`
	return sequence(queryify, &alias, &name, &arguments, &selectionSet)
}

func value() parsec.Parser {
	//   Variable (not implemented)
	//   IntValue
	//   FloatValue (not implemented)
	//   StringValue
	//   BooleanValue (not implemented)
	//   EnumValue (not implemented)
	//   ListValueConst (not implemented)
	//   ObjectValueConst (not implemented)
	return parsec.OrdChoice(valueify, parsec.Int(), parsec.String())
}

// { fields } -> []query
func fieldArrayify(nodes []parsec.ParsecNode) parsec.ParsecNode {
	return fieldArrayify_(nodes)
}

func fieldArrayify_(nodes []parsec.ParsecNode) []Query {
	fields := []Query{}
	for _, qr := range nodes {
		fields = append(fields, qr.(Query))
	}
	return fields
}

// alias field argumetns selectionset
func queryify(nodes []parsec.ParsecNode) parsec.ParsecNode {
	alias, field, arguments, selectionSet := nodes[0], nodes[1], nodes[2], nodes[3]

	query := Query{}

	if field != nil {
		query.Name = field.(string)
	} else {
		return nil
	}

	if alias != nil {
		query.Alias = alias.(string)
	}

	if arguments != nil {
		query.Arguments = arguments.(Arguments)
	}

	if selectionSet != nil {
		subQueries := selectionSet.([]Query)
		if len(subQueries) > 0 {
			query.Fields = subQueries
		}
	}

	return query
}

func dump(p string, nodes []parsec.ParsecNode) {
	fmt.Printf("%s<", p)
	for _, n := range nodes {
		fmt.Printf("%v,", n)
	}
	fmt.Println(">")
}

func dm(p string, cb parsec.Nodify) parsec.Nodify {
	return func(nodes []parsec.ParsecNode) parsec.ParsecNode {
		dump(p, nodes)
		return cb(nodes)
	}
}

func extract(n int) func(nodes []parsec.ParsecNode) parsec.ParsecNode {
	return func(nodes []parsec.ParsecNode) parsec.ParsecNode {
		return nodes[n]
	}
}

func isEmpty(nodes []parsec.ParsecNode) bool {
	for _, n := range nodes {
		if n != nil {
			return false
		}
	}
	return true
}

func argify(nodes []parsec.ParsecNode) parsec.ParsecNode {
	m := Arguments{}
	for _, node := range nodes {
		// dump("argify", nodes)
		nn := node.([]parsec.ParsecNode)
		key := nn[0].(string)
		value := nn[2].(Result)

		m[key] = value
	}
	return m
}

func valueify(nodes []parsec.ParsecNode) parsec.ParsecNode {
	node := nodes[0]
	switch node.(type) {
	case string:
		return String(node.(string))
	case *parsec.Terminal:
		t := node.(*parsec.Terminal)
		switch t.Name {
		case "INT":
			val, _ := strconv.ParseInt(t.Value, 10, 16)
			return Int(val)
		}
	}
	return node
}

func sequence(cb parsec.Nodify, parsers ...*parsec.Parser) parsec.Parser {
	return func(s parsec.Scanner) (node parsec.ParsecNode, next parsec.Scanner) {
		nodes := []parsec.ParsecNode{}
		next = s.Clone()
		for _, parser := range parsers {
			node, next = (*parser)(next)
			nodes = append(nodes, node)
			// dump(fmt.Sprintf("sequence %d", i), nodes)
		}

		if isEmpty(nodes) {
			//			dump("sequence empty", nodes)
			node, next = nil, s
		} else if cb != nil {
			//			dump("sequence cb", nodes)
			node = cb(nodes)
			if node == nil {
				next = s
			}
		} else {
			//			dump("sequence plain", nodes)
			node = nodes
		}

		//		dump("sequence result", nodes)
		return
	}
}
