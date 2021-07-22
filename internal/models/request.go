package models

import "fmt"

type Request struct {
	Id     uint64
	UserId uint64
	Type   uint64
	Text   string
}

// NewRequest create new Request instance
func NewRequest(Id uint64, UserId uint64, Type uint64, Text string) *Request {
	r := new(Request)
	r.Id = Id
	r.UserId = UserId
	r.Type = Type
	r.Text = Text
	return r
}

func (r Request) String() string {
	return fmt.Sprintf("Request{%v, %v, %v, %v}", r.Id, r.UserId, r.Type, r.Text)
}
