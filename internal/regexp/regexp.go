package regexp

import (
	"encoding/gob"
	re "regexp"
)

func init() {
	gob.Register(&Regexp{})
}

type Regexp struct {
	*re.Regexp
}

func (r *Regexp) MarshalText() (text []byte, err error) {
	return []byte(r.String()), nil
}

func (r *Regexp) UnmarshalText(text []byte) (err error) {
	s := string(text)
	r.Regexp, err = re.Compile(s)
	return
}

func (r *Regexp) MarshalBinary() (data []byte, err error) {
	return r.MarshalText()
}

func (r *Regexp) UnmarshalBinary(data []byte) (err error) {
	return r.UnmarshalText(data)
}

func Compile(expr string) (r *Regexp, err error) {
	rxp, err := re.Compile(expr)
	if err != nil {
		return
	}
	r = &Regexp{rxp}
	return
}

func MustCompile(str string) (r *Regexp) {
	r = &Regexp{re.MustCompile(str)}
	return
}

func MatchString(pattern string, s string) (matched bool, err error) {
	return re.MatchString(pattern, s)
}
