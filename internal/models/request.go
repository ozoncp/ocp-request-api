package models

import "fmt"

// Request student's request information
type Request struct {
	Id     uint64
	UserId uint64
	Type   uint64
	Text   string
}

// NewRequest create new Request instance
func NewRequest(id, userId, typeOfRequest uint64, text string) Request {
	return Request{Id: id, UserId: userId, Type: typeOfRequest, Text: text}
}

func (r Request) String() string {
	return fmt.Sprintf("Request{%v, %v, %v, %v}", r.Id, r.UserId, r.Type, r.Text)
}
