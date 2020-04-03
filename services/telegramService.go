package services

import (
	"strconv"
)

// TODO: Whyyyyyy
type Recipient struct {
	TGID int
}

func (r *Recipient) Recipient() string {
	return strconv.Itoa(r.TGID)
}
