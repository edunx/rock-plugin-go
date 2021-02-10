package plugin

import (
	"fmt"
	"github.com/edunx/lua"
	"github.com/edunx/lua/parse"
	pub "github.com/edunx/rock-public-go"
	"os"
	"time"
)

func (p *Plugin) Start() error {
	p.Scripts = make(map[string]PluginFunction , p.C.buffer)

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
	p.Scripts[name] = PluginFunction{
		fn: fn,
		modTime: stat.ModTime().Unix(),
	}

	return fn
}

func (p *Plugin) load( name string ) *lua.LFunction {
	pub.Out.Debug("start load %s plugin" , name)
	plg , ok := p.Scripts[name]
	if ok { return plg.fn }

	return p.compile( name )
}

func (p *Plugin) Check() {
	for name , plg := range p.Scripts {
		filename := fmt.Sprintf("%s/%s.lua" ,p.C.path , name)
		stat , err := os.Stat( filename )
		if os.IsNotExist(err) {
			delete(p.Scripts , name)
			continue
		}

		if stat.ModTime().Unix() != plg.modTime {
			p.compile( name )
			continue
		}
	}
}

func (p *Plugin) Remove( name string ) {
	delete(p.Scripts , name)
}

func (p *Plugin) sync() {
	tk := time.NewTicker( time.Duration( p.C.interval) * time.Millisecond )
	defer tk.Stop()

	for range tk.C { for name , pl := range p.Scripts {
		file := fmt.Sprintf("%s/%s.lua" , p.C.path , name)
		stat , err := os.Stat( file )
		if os.IsNotExist( err ) { p.Remove( name ) ; continue }
		if stat.ModTime().Unix() != pl.modTime { p.compile( name )}
	}}
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
