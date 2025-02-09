package transfer

import "time"

const (
	Add int = iota
	Delete
)

type MessageAddService struct {
	Address string
	Action  int
}
type MessageAddPing struct {
	Address string
	Date    time.Time
	State   bool
}
