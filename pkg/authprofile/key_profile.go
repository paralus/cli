package authprofile

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/rafaylabs/rcloud-cli/pkg/log"

	"github.com/levigross/grequests"
	"github.com/spacemonkeygo/httpsig"
)

type KeyProfile struct {
	Name   string
	URL    string
	Key    string
	Secret string
}

func (p *KeyProfile) Auth(s *grequests.Session) (map[string]string, error) {
	headers := make(map[string]string)

	headers["KEY"] = p.Key
	log.GetLogger().Infof("creating headers")

	return headers, nil
}

func getBodyCheckSum(body []byte) string {
	hash := md5.New()
	hash.Write(body)
	return base64.StdEncoding.EncodeToString(hash.Sum(nil))
}

func (p *KeyProfile) SendRequest(s *grequests.Session, uri, method string, ro *grequests.RequestOptions) (*grequests.Response, error) {
	var resp *grequests.Response
	var payload []byte
	var err error
	path := p.URL + uri

	if ro.RequestBody != nil {
		payload, err = ioutil.ReadAll(ro.RequestBody)
		if err != nil {
			return nil, err
		}
		ro.RequestBody = nil
	}
	if ro.JSON != nil {
		payload, _ = json.Marshal(ro.JSON)
		ro.JSON = nil
	}
	if ro.Files == nil {
		ro.RequestBody = bytes.NewReader(payload)
	}

	reqCallback := func(r *http.Request) error {
		rand.Seed(time.Now().UnixNano())

		r.Header.Add("date", strconv.FormatInt(time.Now().Unix(), 10))
		r.Header.Add("content-md5", getBodyCheckSum(payload))
		r.Header.Add("nonce", strconv.Itoa(rand.Int()))
		r.Header.Add("X-RAFAY-API-KEYID", p.Key)

		signer := httpsig.NewHMACSHA256Signer(p.Key, []byte(p.Secret),
			[]string{"content-md5", "date", "host", "nonce"})
		return signer.Sign(r)
	}

	ro.BeforeRequest = reqCallback

	switch method {
	case "GET":
		resp, err = s.Get(path, ro)
	case "PUT":
		resp, err = s.Put(path, ro)
	case "PATCH":
		resp, err = s.Patch(path, ro)
	case "DELETE":
		resp, err = s.Delete(path, ro)
	case "POST":
		resp, err = s.Post(path, ro)
	case "HEAD":
		resp, err = s.Head(path, ro)
	case "OPTIONS":
		resp, err = s.Options(path, ro)
	default:
		err = errors.New("Unknown HTTP method")
		resp = nil
	}

	log.GetLogger().Debugf("%s %s %s", method, path, string(payload))

	if err != nil {
		log.GetLogger().Debugf("http request error: %s. resp %s", err, resp)
	} else {
		log.GetLogger().Debugf("http response ok: %s",
			strings.TrimSuffix(resp.String(), "\n"))
	}

	return resp, err
}
