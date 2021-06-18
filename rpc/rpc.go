package rpc

type Query struct {
	IP      string
	Project string
	Ua      string
}

type Answer struct {
	Ok bool
}

type Error struct {
	Error error
}
