package convert

import "strconv"

type Str string

func (s Str) String() string {
	return string(s)
}

func (s Str) Int() (int, error) {
	return strconv.Atoi(s.String())
}

func (s Str) MustInt() int {
	i, _ := s.Int()
	return i
}

func (s Str) UInt() (uint, error) {
	v, err := strconv.Atoi(s.String())
	return uint(v), err
}

func (s Str) MustUInt() uint {
	v, _ := s.UInt()
	return v
}

func (s Str) UInt32() (uint32, error) {
	v, err := strconv.Atoi(s.String())
	return uint32(v), err
}

func (s Str) MustUInt32() uint32 {
	v, _ := s.UInt32()
	return v
}
