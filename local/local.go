package local

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

type Album interface {
	Pictures() []ProfilePicture
}
