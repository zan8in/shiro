package req

import (
	"github.com/zan8in/retryablehttp"
)

var (
	defaultTimeout = 6
)

type Options struct {
	Proxy   string
	Timeout int
	Retries int
}

type Req struct {
	*retryablehttp.Client
}

func New(opts *Options) (*Req, error) {
	if opts == nil {
		opts = &Options{}
	}

	if opts.Timeout == 0 {
		opts.Timeout = defaultTimeout
	}

	if opts.Retries == 0 {
		opts.Retries = 3
	}

	po := &retryablehttp.DefaultPoolOptions
	po.Proxy = opts.Proxy
	po.Timeout = opts.Timeout
	po.Retries = opts.Retries
	po.DisableRedirects = true

	err := retryablehttp.InitClientPool(po)
	if err != nil {
		return nil, err
	}

	c, err := retryablehttp.GetPool(po)
	if err != nil {
		return nil, err
	}

	return &Req{c}, nil
}
