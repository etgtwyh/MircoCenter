package ioc

import (
	"fmt"
	"gitee.com/King_of_Universe_Sailing/wcenter/ServiceCenter/utils/Log"
	"github.com/gin-gonic/gin"
)

var logger = Log.NewLogger("etc/log/ioc.log", 2, 7, 100)

var ns = &NameSpace{N_Storage: map[string]*Container{
	"Conf":       NewContainer(),
	"Controller": NewContainer(),
	"Api":        NewContainer(),
}}

type NameSpace struct {
	N_Storage map[string]*Container
}

type RegisterApi interface {
	Ioc
	RegisterApi(router gin.IRouter) error
}

func Conf() *Container {
	return ns.N_Storage["Conf"]
}

func Controller() *Container {
	return ns.N_Storage["Controller"]
}

func Api() *Container {
	return ns.N_Storage["Api"]
}

func RegisterGinApi(router gin.IRouter) error {
	if err := ns.RegisterApi(router); err != nil {
		return err
	}
	return nil
}

func Init() error {
	logger.Info().Msg("IOC Initializing(正在依赖注入)")
	return ns.Init()
}

func Destroy() error {
	return ns.Destroy()
}

func (n *NameSpace) RegisterApi(router gin.IRouter) error {
	for c_name := range n.N_Storage {
		container := n.N_Storage[c_name]
		for ioc := range container.C_Storage {
			i := container.C_Storage[ioc]
			if api, ok := i.(RegisterApi); ok {
				logger.Info().Msg(fmt.Sprintf("%s容器正在给%s注入api依赖", c_name, ioc))
				if err := api.RegisterApi(router); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (n *NameSpace) Init() error {
	for c_name := range n.N_Storage {
		container := n.N_Storage[c_name]
		for ioc := range container.C_Storage {
			logger.Info().Msg(fmt.Sprintf("%s容器正在给%s注入依赖", c_name, ioc))
			i := container.C_Storage[ioc]
			if err := i.Init(); err != nil {
				return err
			}
		}
	}
	return nil
}

func (n *NameSpace) Destroy() error {
	for c_name := range n.N_Storage {
		container := n.N_Storage[c_name]
		for ioc := range container.C_Storage {
			logger.Info().Msg(fmt.Sprintf("%s容器正在给%s注销依赖", c_name, ioc))
			i := container.C_Storage[ioc]
			if err := i.Destroy(); err != nil {
				return err
			}
		}
	}
	return nil
}
