# MircoCenter
微服务中心
├─MircoCenter--------------------->web前端
│  ├─axios
│  ├─public
│  ├─router
│  ├─src
│  │  ├─components
│  │  │  ├─container
│  │  │  │  ├─content
│  │  │  │  └─sider
│  │  │  ├─cpu
│  │  │  ├─etree
│  │  │  ├─icmp
│  │  │  ├─memory
│  │  │  ├─metrics
│  │  │  └─tree
│  │  └─svg
│  ├─store
│  │  └─modules
│  └─utils
│      └─api
└─ServiceCenter---------------------------------->服务中心后端
    ├─apps
    │  ├─MetricManage
    │  │  ├─controller
    │  │  └─pb
    │  └─ServiceManage
    │      ├─api
    │      ├─controller
    │      └─pb
    ├─conf
    ├─docs
    │  └─images
    ├─enum
    ├─etc
    │  └─log
    ├─execption
    ├─ioc
    ├─middlewares
    │  └─http
    ├─package
    │  ├─MetricMonitor
    │  └─ServiceUtils
    ├─reponse
    ├─test
    └─utils
        └─Log