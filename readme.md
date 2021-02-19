---
   磐石lua字节缓存插件
---

# 配置
```lua
    local p = rock.plugin{
        path = "resource/plugin",
        interval = 1000, --ms
    }

   --测试
   p.pcall("default")


  -- resource/plugin/default.lua 文件内容
  print( "hello ")
```

# 调用
```golang
    import (
        plugin "github.com/edunx/rock-plugin-go"
    )

    //注入代码
    plugin.LuaInjectApi(L , rock)

    //引用
    p := plugin.CheckPluginUserData(L , 1)
    p.pcall("default")
    
```