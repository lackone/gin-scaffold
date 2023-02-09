package id

import "github.com/rs/xid"

type ID struct {
}

func NewIDService(params ...interface{}) (interface{}, error) {
	return &ID{}, nil
}

func (s *ID) NewID() string {
	return xid.New().String()
}
