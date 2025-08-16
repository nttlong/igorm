package wx_container

import (
	"vdi"
)

type Server struct {
}
type WxContainer struct {
	Handler vdi.Singleton[WxContainer, Server] // Cho phép inject router
}
