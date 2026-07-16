package gender

import (
	"bytes"
	"errors"
)

var (
	ErrInvalidGender = errors.New("invalid gender value")
)

type Gender uint8

const (
	Unknown Gender = 0 + iota
	Male
	Female
)

var genderKeys = []string{"unknown", "male", "female"}

// fmt.Stringer
func (g Gender) String() string {
	if g >= Unknown && g <= Female {
		return genderKeys[g]
	}
	return "unknown"
}

func (g Gender) MarshalJSON() ([]byte, error) {
	switch g {
	case Male:
		return []byte{'"', 'm', '"'}, nil
	case Female:
		return []byte{'"', 'f', '"'}, nil
	}
	return []byte{'"', 'u', '"'}, nil
}

func (g *Gender) UnmarshalJSON(b []byte) (err error) {
	if len(b) == 0 {
		*g = Unknown
		return
	}
	r := bytes.Runes(b)
	if r[0] == '"' && r[len(r)-1] == '"' {
		r = r[1 : len(r)-1]
	}
	switch c := r[0]; c {
	case 'm', 'M', '1', '男':
		*g = Male
	case 'f', 'F', '2', '女':
		*g = Female
	case 'u', 'U', '0':
		*g = Unknown
	default:
		err = ErrInvalidGender
	}
	return
}
