package plugin

import (
	"github.com/edunx/lua"
)

type Config struct {
	path                string
	buffer    int   //缓存大小
	interval  int
}

type PluginFunction struct {
	fn    *lua.LFunction
	modTime  int64
}

type Plugin struct {
	C  Config
	Scripts   map[string]PluginFunction
}