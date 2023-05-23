package main

import (
	"fmt"

	"github.com/zan8in/goflags"
	"github.com/zan8in/shiro"
)

func main() {
	ShowBanner()

	options, err := NewOptions()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	s, err := shiro.NewShiro()
	if err != nil {
		panic(err)
	}

	if len(options.TargetFile) > 0 {
		result, err := s.RunMulti(*options)
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

	if len(options.Target) > 0 {
		result, err := s.Run(*options)
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

}

func NewOptions() (*shiro.Options, error) {

	options := &shiro.Options{}
	flagSet := goflags.NewFlagSet()
	flagSet.SetDescription(`shiro`)

	flagSet.CreateGroup("input", "Target",
		flagSet.StringVarP(&options.Target, "target", "t", "", "target URLs/hosts to scan"),
		flagSet.StringVarP(&options.TargetFile, "target-file", "T", "", "list of target URLs/hosts to scan (one per line)"),
	)

	flagSet.CreateGroup("keys", "Keys",
		flagSet.StringVarP(&options.ShiroKeysFile, "key-file", "K", "", "PoC file or directory to scan"),
	)

	flagSet.CreateGroup("rate-limit", "Rate-Limit",
		flagSet.IntVarP(&options.RateLimitKey, "rate-limit-key", "rlk", 60, "maximum number of requests to send per second"),
		flagSet.IntVarP(&options.ConcurrencyKey, "concurrency-key", "ck", 25, "maximum number of afrog-pocs to be executed in parallel"),
		flagSet.IntVarP(&options.RateLimitTarget, "rate-limit-target", "rlt", 6, "maximum number of requests to send per second"),
		flagSet.IntVarP(&options.ConcurrencyTarget, "concurrency-target", "ct", 5, "maximum number of afrog-pocs to be executed in parallel"),
	)

	_ = flagSet.Parse()

	if err := verifyOptions(options); err != nil {
		return options, err
	}

	return options, nil
}

func verifyOptions(opt *shiro.Options) error {
	if len(opt.Target) == 0 && len(opt.TargetFile) == 0 {
		return fmt.Errorf("target or target file must be specified")
	}

	return nil
}

const Version = "1.0.1"

func ShowBanner() {
	fmt.Printf("\n|\tS H I R O\t%s\n\n", Version)
}
