package iping

type PingerInterface interface {
	DoPinging()
	AddService(address string)
	DeleteService(address string)
}
