package transfer

import "time"

type PingServiceTransferObject struct {
	Address     string
	LastPing    PingTransferObject
	LastSuccess PingTransferObject
	State       bool
}
type PingTransferObject struct {
	Date time.Time
}
