package starwars

// enum Episode { NEWHOPE, EMPIRE, JEDI }

// interface Character {
//   id: String!
//   name: String
//   friends: [Character]
//   appearsIn: [Episode]
// }

// type Human : Character {
//   id: String!
//   name: String
//   friends: [Character]
//   appearsIn: [Episode]
//   homePlanet: String
// }

// type Droid : Character {
//   id: String!
//   name: String
//   friends: [Character]
//   appearsIn: [Episode]
//   primaryFunction: String
// }

// type Query {
//   hero(episode: Episode): Character
//   human(id: String!): Human
//   droid(id: String!): Droid
// }

type Episode string

const (
	NEWHOPE Episode = "NEWHOPE"
	EMPIRE  Episode = "EMPIRE"
	JEDI    Episode = "JEDI"
)

type Character struct {
	Id        string `gql:"required"`
	Name      string
	Friends   []*Character
	AppearsIn []Episode
}

type Human struct {
	Character
	HomePlanet string
}

type Droid struct {
	Character
	PrimaryFunction string
}

type Query interface {
	Hero(episode Episode) *Character
	Human(id string) *Human
	Droid(id string) *Droid
}
