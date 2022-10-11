package model

type GLock struct {
	Id uint64
	Rd uint32
	C  chan interface{}
}
