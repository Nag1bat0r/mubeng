package mubeng

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"

	"golang.org/x/net/proxy"
)

// Transport to auto-switch transport between HTTP/S or SOCKSv5 proxies.
// Depending on the protocol scheme, returning value of http.Transport with Dialer or Proxy.
func Transport(p string) (tr *http.Transport, err error) {
	proxyURL, err := url.Parse(p)
	if err != nil {
		return nil, err
	}

	ProxyPassword, pp := proxyURL.User.Password()
        if pp = false {
		dialer, err := proxy.SOCKS5("tcp", proxyURL.Host, nil, proxy.Direct)
		if err != nil {
			return nil, err
		}
        }
	else
	{
	        auth :=proxy.Auth
		{
	        	User: proxyURL.User.Username(),
	        	Password: ProxyPassword,
	        }
		dialer, err := proxy.SOCKS5("tcp", proxyURL.Host, &auth, proxy.Direct)
		if err != nil {
			return nil, err
		}

	}

	switch proxyURL.Scheme {
	case "socks5":
		tr = &http.Transport{
			Dial: dialer.Dial,
		}
	case "http":
		tr = &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		}
	default:
		return nil, fmt.Errorf("unsupported proxy protocol scheme: %s", proxyURL.Scheme)
	}

	tr.DisableKeepAlives = true
	tr.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	return tr, nil
}
