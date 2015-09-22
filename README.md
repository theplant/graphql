# Experimental Go GraphQL library

Components:

* Query parser (incomplete)
* Query "transformer" (incomplete, could use a better name)
* Schema parser (imaginary)
* Schema type generator (imaginary)

## How to run examples

Run all examples as tests:

```
go test ./examples/

```

Run a single example:

```
go test ./examples/ -run=ExampleQuery
```

## Todo

- [ ] unify input types in schema and parser

### Query Parser

- [ ] support more input types
- [ ] support for fragments, unions, and other GraphQL features

### Query Transformer

- [ ] support named parameters
- [ ] other stuff!

### Schema Parser

- [ ] write it

### Schema Type Generator

- [ ] write it

# Query Parser

Parser implemented with [`goparsec`](https://github.com/prataprc/goparsec).

`graphql.Parse` will decode this:

```
{
  User(id: "3500401") {
    Id
    Name
    pic: ProfilePicture(size:50) {
      Uri
      Width
    }
  }
}
```

into this:

```go
graphql.Query{
	Name:   "query",
	Fields: graphql.Fields{
		graphql.Query{
			Name:	  "User",
			Arguments: graphql.Arguments{"id": graphql.String("3500401")},
			Fields:    graphql.Fields{
				graphql.Query{Name: "Id"},
				graphql.Query{Name: "Name"},
				graphql.Query{
					Name:	   "ProfilePicture",
					Alias:	   "pic",
					Arguments: graphql.Arguments{"size": graphql.Int(50)},
					Fields:    graphql.Fields{
						graphql.Query{Name: "Uri"},
						graphql.Query{Name: "Width"},
					},
				},
			},
		},
	},
}
```

# Query transformer

Given a schema like this:

```
type Query {
     User(id: String!): User
}

type User {
     Id: String
     Name: String
     ProfilePicture(size: Int!): ProfilePicture
}

type ProfilePicture {
     Uri: String
     Width: String
     Height: String
}
```

The transformer will expect a set of interfaces like so:

```go
package schema

type Query interface {
	User(id string) User
}

type User interface {
	Id() string
	Name() string
	ProfilePicture(size int) ProfilePicture
}

type ProfilePicture interface {
	Uri() string
	Width() int
	Height() int
}
```

And `graphql.Transform(query graphql.Query, context schema.Query)` will transform a query like:

```
{
  User(id: "3500401") {
    Id
    Name
    pic: ProfilePicture(size:50) {
      Uri
      Width
    }
  }
}
```

into a basic data structures like:

```go
map[string]interface{}{
  "User": map[string]interface{}{
    "Id": "3500401",
    "Name": "Mr. Ed",
    "pic": map[string]interface{}{
      "Uri": "http://host/path/to/img/at/50px",
      "Width": 50
    },
  },
}
```

suitable for serialisation to JSON.

# Schema parser (imaginary)

WRITE/IMPLEMENTME

# Schema type generator (imaginary)

Given a schema like this:

```
type Query {
     User(id: String!): User
}

type User {
     Id: String
     Name: String
     ProfilePicture(size: Int!): ProfilePicture
}

type ProfilePicture {
     Uri: String
     Width: String
     Height: String
}
```

The generator should generate a set of interfaces like this

```go
package schema

type Query interface {
	User(id string) User
}

type User interface {
	Id() string
	Name() string
	ProfilePicture(size int) ProfilePicture
}

type ProfilePicture interface {
	Uri() string
	Width() int
	Height() int
}
```

that the client application can then implement.

The reason for generating the interfaces is to reintroduce some semblance of type safety (because otherwise everything will have to pass through the `interface{}` vortex).
