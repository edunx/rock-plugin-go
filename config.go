package plugin

import (
	"github.com/edunx/lua"
	"sync"
)

type Config struct {
	path                string
	interval  int
}

type PluginFunction struct {
	fn    *lua.LFunction
	modTime  int64
}

type Plugin struct {
	C  Config
	Cache sync.Map
}