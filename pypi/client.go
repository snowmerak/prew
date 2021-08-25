package pypi

import (
	"encoding/json"
	"net"
	"net/http"
	"runtime"
	"time"
)

var c = func() *http.Client {
	tmp := http.Client{}
	tmp.Transport = &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
	tmp.Timeout = time.Second * 30
	return &tmp
}()

var UserAgent string = "Prew/1.0 (" + runtime.GOOS + "; " + runtime.GOARCH + "; " + runtime.Version() + "; https://github.com/snowmerak/prew)"

func GetPackageInfo(Name string, Version string) (Package, error) {
	var URL = "https://pypi.org/pypi/" + Name + "/" + Version + "/json"
	if URL == "" {
		URL = "https://pypi.org/pypi/" + Name + "/json"
	}

	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		return Package{}, err
	}
	req.Header.Set("User-Agent", UserAgent)

	resp, err := c.Do(req)
	if err != nil {
		return Package{}, err
	}
	defer resp.Body.Close()

	var p Package
	err = json.NewDecoder(resp.Body).Decode(&p)
	if err != nil {
		return Package{}, err
	}

	return p, nil
}
