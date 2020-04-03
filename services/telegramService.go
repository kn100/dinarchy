package services

import (
	"strconv"
)

type Recipient struct {
	TGID int
}

func (r *Recipient) Recipient() string {
	return strconv.Itoa(r.TGID)
}
