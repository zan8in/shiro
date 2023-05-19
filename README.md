# shiro
Usage of the [afrog](https://github.com/zan8in/afrog) calling library.

# Examples

## Single Target

```
s, err := shiro.NewShiro()
if err != nil {
    panic(err)
}

result, err := s.Run(shiro.Options{
    Target: "http://x.x.x.x:8090",
})
if err != nil {
    panic(err)
}

if result == nil {
    fmt.Println("result is nil")
    return
}

fmt.Printf("target: %s\n", result.Target)
fmt.Printf("shiroKey: %s\n", result.ShiroKey)
fmt.Printf("rememberMe: %s\n", result.RememberMe)
```


# Forked and adapted from
https://github.com/Peony2022/shiro_killer