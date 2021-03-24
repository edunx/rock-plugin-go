package plugin

import (
	"github.com/edunx/lua"
	pub "github.com/edunx/rock-public-go"
)

func (self *Plugin) LPcall(L *lua.LState , args *lua.Args) lua.LValue {
	self.Pcall(L , args.CheckString(L , 1) )
	return lua.LNil
}

func (self *Plugin) ToLightUserData(L *lua.LState) *lua.LightUserData {
	return L.NewLightUserData( self )
}

func (self *Plugin) Index(L *lua.LState , key string ) lua.LValue {
	if key == "pcall"  { return lua.NewGFunction( self.LPcall )}

	return lua.LNil
}

func createPluginLightUserData(L *lua.LState , args *lua.Args) lua.LValue {
	opt := args.CheckTable(L , 1)

	p := &Plugin{
		C: Config{
			path: opt.CheckString("path" , "./"),
			interval: opt.CheckInt("interval" , 1000),
		},
	}

	if e := p.Start(); e != nil {
		pub.Out.Debug("start plugin fail , err: %v" , e)
		return lua.LNil
	}

	return p.ToLightUserData(L)
}

func LuaInjectApi(L *lua.LState , parent *lua.LTable) {
	L.SetField(parent , "plugin" , lua.NewGFunction( createPluginLightUserData ) )
}
