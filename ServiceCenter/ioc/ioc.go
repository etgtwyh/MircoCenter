package ioc

type Ioc interface {
	Init() error
	Destroy() error
}
