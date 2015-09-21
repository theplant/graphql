package examples

import (
	"fmt"

	"github.com/theplant/graphql/graphql"
)

func ExampleQuery() {

	q := graphql.Query{
		Name:      "User",
		Arguments: graphql.Arguments{"id": graphql.String("3500401")},
		Fields: graphql.Fields{
			graphql.Query{Name: "Id"},
			graphql.Query{Name: "Name"},
			//			graphql.Query{Name: "Isviewfriend"},
			graphql.Query{
				Name:      "ProfilePicture",
				Arguments: graphql.Arguments{"size": graphql.Int(50)},
				Fields: graphql.Fields{
					graphql.Query{Name: "Uri"},
					graphql.Query{Name: "Width"},
					// 	graphql.Query{Name: "Height"},
				},
			},
		},
	}

	fmt.Println(q)

	k := Query{}
	c, err := graphql.Transform(q, k)
	fmt.Println(c)
	fmt.Println(err)

	q = graphql.Query{
		Name: "Pictures",
		Fields: graphql.Fields{
			graphql.Query{Name: "Uri"},
			graphql.Query{Name: "Width"},
		},
	}
	c, err = graphql.Transform(q, Album{})
	fmt.Println(c)
	fmt.Println(err)

	// Output: {User  map[id:3500401] [{Id  map[] []} {Name  map[] []} {ProfilePicture  map[size:50] [{Uri  map[] []} {Width  map[] []}]}]}
	//map[Id:3500401 Name:Mr. Ed ProfilePicture:map[Width:50 Uri:http://50]]
	//<nil>
	//map[Pictures:{[map[Uri:http://10 Width:10] map[Width:20 Uri:http://20]]}]
	//<nil>

}
