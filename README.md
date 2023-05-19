# shiro


# Examples

## Single Target

```
s, err := shiro.NewShiro()
if err != nil {
    panic(err)
}

err = s.Run(shiro.Options{
    Target: "http://x.x.x.x:8090",
})
if err != nil {
    panic(err)
}

fmt.Println("shiro scanner finished")
```


# Forked and adapted from
https://github.com/Peony2022/shiro_killer