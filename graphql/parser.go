package graphql

import (
	"errors"
	"fmt"

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

var y, fields parsec.Parser

func Parse(input string) (query Query, err error) {
	s := parsec.NewScanner([]byte(input))
	result, next := y(s)
	if next.Endof() {
		query = result.([]parsec.ParsecNode)[0].(Query)
	} else {
		err = errors.New("Failed to parse input")
	}

	return
}

func init() {
	fields = parsec.Many(nil, field(&y))
	y = parsec.And(extract(1), parsec.Token(`{`, ""), fields, parsec.Token(`}`, ""))
}

func field(query *parsec.Parser) parsec.Parser {
	identFn := func(nodes []parsec.ParsecNode) parsec.ParsecNode {
		return nodes[0].(*parsec.Terminal).Value
	}

	field := parsec.Maybe(identFn, parsec.Ident())

	alias := parsec.And(identFn, parsec.Ident(), parsec.Token(`:`, ""))

	arguments := parsec.And(
		extract(1),
		parsec.Token(`\(`, ""),
		parsec.Many(argify, parsec.And(nil, parsec.Ident(), parsec.Token(`:`, ""), parsec.Int())),
		parsec.Token(`\)`, ""))

	return sequence(queryify, &alias, &field, &arguments, query)
}

func queryify(nodes []parsec.ParsecNode) parsec.ParsecNode {
	alias, field, arguments, subquery := nodes[0], nodes[1], nodes[2], nodes[3]

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

	if subquery != nil {
		subQueries := subquery.([]parsec.ParsecNode)
		if len(subQueries) > 0 {
			query.Fields = []Query{}
			for _, qr := range subQueries {
				query.Fields = append(query.Fields, qr.(Query))
			}
		}
	}

	return query
}

func dump(p string, nodes []parsec.ParsecNode) {
	fmt.Printf("%s[", p)
	for _, n := range nodes {
		fmt.Printf("%v,", n)
	}
	fmt.Println("]")
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
		nn := node.([]parsec.ParsecNode)
		m[nn[0].(*parsec.Terminal).Value] = String(nn[2].(*parsec.Terminal).Value)
	}
	return m
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
