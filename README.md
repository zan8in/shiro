# shiro
Usage of the [afrog](https://github.com/zan8in/afrog) calling library.

# Running Shiro

## Single Target
```sh
shiro -t https://example.com
```

## Multiple Target
```sh
shiro -T urls.txt
```

## Custom Keys
```sh
shiro -T urls.txt -K keys.txt
```

## More Command
```sh
shiro -h
```

# Examples

## Single Target

```go
package main

import (
	"fmt"
	"time"

	"github.com/zan8in/shiro"
)

func main() {
	s, err := shiro.NewShiro()
	if err != nil {
		panic(err)
	}

	result, err := s.Run(shiro.Options{
		Target:            "http://x.x.x.x:8090",
		RateLimitKey:      time.Duration(60),
		ConcurrencyKey:    25,
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

}
```

## Multi Target

```go
package main

import (
	"fmt"
	"time"

	"github.com/zan8in/shiro"
)

func main() {
	s, err := shiro.NewShiro()
	if err != nil {
		panic(err)
	}

	result, err := s.RunMulti(shiro.Options{
		TargetFile:        "./1.txt",
		ShiroKeysFile:     "./keys.txt",
		RateLimitKey:      time.Duration(60),
		RateLimitTarget:   time.Duration(5),
		ConcurrencyKey:    25,
		ConcurrencyTarget: 2,
	})
	if err != nil {
		panic(err)
	}

	for v := range result {
		fmt.Printf("target: %s\n", v.Target)
		fmt.Printf("shiroKey: %s\n", v.ShiroKey)
		fmt.Printf("rememberMe: %s\n", v.RememberMe)
		fmt.Println("-----------------------")
	}

}

```


# Forked and adapted from
https://github.com/Peony2022/shiro_killer