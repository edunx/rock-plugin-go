package plugin

import (
	"github.com/edunx/lua"
	pub "github.com/edunx/rock-public-go"
)

const (
	MT string = "ROCK_PLUGIN_GO_MT"
)

func CheckPluginUserData(L *lua.LState , idx int) *Plugin {
	ud := L.CheckUserData( idx )
	v , ok := ud.Value.(*Plugin)
	if ok { return v }

	L.RaiseError("#%d must be Plugin userdata , got fail" , idx)
	return nil
}

func CreatePluginUserData(L *lua.LState) int {
	opt := L.CheckTable(1)

	p := &Plugin{
		C: Config{
			path: opt.CheckString("path" , "./"),
			buffer: opt.CheckInt("buffer" , 1024),
			interval: opt.CheckInt("interval" , 1000),
		},
	}

	if e := p.Start(); e != nil {
		pub.Out.Debug("start plugin fail , err: %v" , e)
		return 0
	}

	ud := L.NewUserDataByInterface(p , MT)
	L.Push(ud)

	return 1

}


func LuaInjectApi(L *lua.LState , parent *lua.LTable) {
	mt := L.NewTypeMetatable( MT )

	L.SetField(mt , "__index" , L.NewFunction(Get))
	L.SetField(mt , "__newindex" , L.NewFunction(Set))

	L.SetField(parent , "plugin" , L.NewFunction(CreatePluginUserData))
}

func Get(L *lua.LState) int {
	self := CheckPluginUserData(L , 1)
	name := L.CheckString(2)
	switch name {
	case "pcall":
		L.Push(L.NewFunction( func (L *lua.LState) int {
			name := L.CheckString(1)
			self.pcall(L , name)
			return 0
		}))
	default:
		L.Push(lua.LNil)
	}

	return 1
}

func Set(L *lua.LState) int {
	return 0
}

func (p *Plugin) ToUserData(L *lua.LState) *lua.LUserData {
	return L.NewUserDataByInterface( p , MT )
}
