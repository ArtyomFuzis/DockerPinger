package transfer

import "time"

type PingServiceTransferObject struct {
	Address  string
	LastPing PingTransferObject
}
type PingTransferObject struct {
	Date time.Time
}
