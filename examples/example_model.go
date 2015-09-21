package examples

import (
	"fmt"

	"github.com/theplant/graphql/local"
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
