package main

// Input
import (
	"fmt"

	"./graphql"
	"./local"
)

// {
//   user(id: 3500401) {
//     id,
//     name,
//     isViewerFriend,
//     profilePicture(size: 50)  {
//       uri,
//       width,
//       height
//     }
//   }
// }

// Output
// {
//   "user" : {
//     "id": 3500401,
//     "name": "Jing Chen",
//     "isViewerFriend": true,
//     "profilePicture": {
//       "uri": "http://someurl.cdn/pic.jpg",
//       "width": 50,
//       "height": 50
//     }
//   }
// }

// use args to find context
// use context + field list to find fields

type Query struct{}

func (Query) User(id string) local.User {
	// find user
	return User{id: id, name: "Mr. Ed"}
}

type User struct {
	id   string
	name string
}

func (u User) Id() string {
	return u.id
}

func (u User) Name() string {
	return u.name
}

func (u User) ProfilePicture(size int) local.ProfilePicture {
	return ProfilePicture{size: size}
}

type ProfilePicture struct {
	size int
}

func (p ProfilePicture) Uri() string {
	return fmt.Sprintf("http://%d", p.size)
}

func (p ProfilePicture) Width() int {
	return p.size
}

func (p ProfilePicture) Height() int {
	return p.size
}

type Album struct{}

func (Album) Pictures() []ProfilePicture {
	return []ProfilePicture{ProfilePicture{size: 10}, ProfilePicture{size: 20}}
}

func main() {
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

	graphql.Printq(q)

	var v = make(map[string]interface{})

	fmt.Println("res:", v)

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
}
