package plugin

import (
	"fmt"
	"github.com/edunx/lua"
	"github.com/edunx/lua/parse"
	pub "github.com/edunx/rock-public-go"
	"os"
	"sync"
	"time"
)

func (p *Plugin) Start() error {
	p.Cache = sync.Map{ }

	go p.sync()
	return nil
}

func (p *Plugin) compile( name string ) *lua.LFunction {
	pub.Out.Debug("start compiler %s plugin" , name)
	filename := fmt.Sprintf("%s/%s.lua" , p.C.path , name)
	stat , err := os.Stat( filename )
	if os.IsNotExist( err) {
		pub.Out.Debug("plugin %s not found" , name)
		return nil
	}

	file , err := os.Open( filename )
	if err != nil {
		pub.Out.Debug("plugin %s %v" , name , err)
		return nil
	}
	defer file.Close()

	chunk , err := parse.Parse(file , filename)
	if err != nil {
		pub.Out.Debug("plugin parse fail , %s %v" , name , err)
		return nil
	}

	proto , err := lua.Compile(chunk , filename)
	if err != nil {
		pub.Out.Debug("plugin compile fail , %s %v" , name , err)
		return nil
	}

	fn := pub.VM.NewFunctionFromProto( proto )
	p.Cache.Store(name , PluginFunction{
		fn: fn,
		modTime: stat.ModTime().Unix(),
	})

	return fn
}

func (p *Plugin) load( name string ) *lua.LFunction {
	pub.Out.Debug("start load %s plugin" , name)
	v , ok := p.Cache.Load(name)
	if ok { return v.(PluginFunction).fn }

	return p.compile( name )
}

func (p *Plugin) sync() {
	tk := time.NewTicker( time.Duration( p.C.interval) * time.Millisecond )
	defer tk.Stop()

	for range tk.C {
		p.Cache.Range( func(k interface{}, v interface{}) bool {
			name := k.(string)
			pl  := v.(PluginFunction)

			file := fmt.Sprintf("%s/%s.lua" , p.C.path , name)
			stat , err := os.Stat( file )
			if os.IsNotExist( err ) {
				p.Cache.Delete( name )
				return false //next
			}

			if stat.ModTime().Unix() != pl.modTime {
				p.compile( name )
				return false //next
			}

			return false
		})
	}
}

func (p *Plugin) pcall(L *lua.LState , name string) {
	fn := p.load( name )
	if fn == nil { return }

	L.Push(fn)
	err := L.PCall(0 , 0 , nil)
	if err != nil {
		pub.Out.Err("plugin pcall fail , err: %v" , err)

	}
}
