package ioc

type Container struct {
	C_Storage map[string]Ioc
}

// 构造函数
func NewContainer() *Container {
	return &Container{C_Storage: map[string]Ioc{}}
}

func (c *Container) Register(AppName string, ioc Ioc) {
	c.C_Storage[AppName] = ioc
}

func (c *Container) Get(AppName string) any {
	return c.C_Storage[AppName]
}
