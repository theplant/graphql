package examples

import (
	"fmt"

	"github.com/theplant/graphql/graphql"
)

func ExampleParser() {
	qr, err := graphql.Parse("{ Pictures { Width u: Width } Pictures { Width }  }")
	fmt.Println(qr)
	fmt.Println(err)
	// Output: {Pictures  map[] [{Width  map[] []} {Width u map[] []}]}
	// <nil>

}

func ExampleParserWithParam() {

	qr, err := graphql.Parse("{ User(id: \"3500401\") { Id Name pic: ProfilePicture(size:50) { Uri Width } } }")
	fmt.Println(qr)
	fmt.Println(err)
	c, err := graphql.Transform(qr[0], Query{})
	fmt.Println(c)
	fmt.Println(err)
	// Output: {User  map[id:"3500401"] [{Id  map[] []} {Name  map[] []} {ProfilePicture pic map[size:50] [{Uri  map[] []} {Width  map[] []}]}]}
	//<nil>
	//map[Id:"3500401" Name:Mr. Ed pic:map[Uri:http://50 Width:50]]
	//<nil>

}
