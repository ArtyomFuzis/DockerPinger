package transfer

import "time"

type MessageAddService struct {
	Address string
}
type MessageAddPing struct {
	Address string
	Date    time.Time
	State   bool
}
