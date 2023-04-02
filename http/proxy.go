package http

import (
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"io"
	"net/http"
	"net/url"
)

var hopHeaders = []string{
	"Connection",
	"Keep-Alive",
	"Proxy-Authenticate",
	"Proxy-Authorization",
	"Te",
	"Trailers",
	"Transfer-Encoding",
	"Upgrade",
}

func delHopHeaders(header http.Header) {
	for _, h := range hopHeaders {
		header.Del(h)
	}
}

func copyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}

func Proxy(c echo.Context) error {
	proxyDest := c.QueryParam("url")
	client := &http.Client{}

	u, err := url.Parse(proxyDest)
	if err != nil {
		return err
	}

	if u.Host != "2ch.hk" && u.Host != "i.4cdn.org" {
		return errors.New("not allowed host")
	}

	delHopHeaders(c.Request().Header)

	resp, err := client.Get(proxyDest)
	if err != nil {
		return err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	delHopHeaders(resp.Header)

	copyHeader(c.Response().Writer.Header(), resp.Header)
	c.Response().Writer.WriteHeader(resp.StatusCode)
	_, err = io.Copy(c.Response().Writer, resp.Body)
	if err != nil {
		return err
	}

	return nil
}
