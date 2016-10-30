package oghttp

import (
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"ogbot/helpers"
	"ogbot/ogdata"
	"strings"
)

func MakeHttpClient(logger helpers.Logger, dump bool) *http.Client {
	// set basic cookie jar
	jar, _ := cookiejar.New(nil)
	var cookies []*http.Cookie
	firstCookie := &http.Cookie{
		Name:  "OG_lastServer",
		Value: "s131-en.ogame.gameforge.com",
	}
	cookies = append(cookies, firstCookie)
	cookieURL, _ := url.Parse("http://en.ogame.gameforge.com")
	jar.SetCookies(cookieURL, cookies)
	// set client to make the login request
	client := &http.Client{
		Jar: jar,
	}
	// set up redirection policy and break when encouters second redirection
	// in order to get the response header that contains cookie information
	// vital to make further requests
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		helpers.LogMark("Redirection...", logger)
		if helpers.CountWords("PHPSESSID", req.URL.String()) > 0 {
			helpers.LogMark("Stop redirection...", logger)
			return fmt.Errorf("stop redirect")
		}
		if dump {
			helpers.DumpRequest(req, logger)
		}
		return nil
	}
	return client
}

func MakeLoginRequest(meta ogdata.MetaData, logger helpers.Logger, dump bool) *http.Request {
	// fill in login form
	data := &url.Values{}
	data.Add("kid", "")
	data.Add("uni", meta.Uni+"-"+meta.Lang+".ogame.gameforge.com")
	data.Add("login", meta.Login)
	data.Add("pass", meta.Pass)
	req, _ := http.NewRequest("POST", "https://"+meta.Lang+".ogame.gameforge.com/main/login",
		strings.NewReader(data.Encode()))
	// set up basic header
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Host", "en.ogame.gameforge.com")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("Content-Length", "77")
	req.Header.Add("Pragma", "no-cache")
	req.Header.Add("Cache-Control", "no-cache")
	req.Proto = "HTTP/1.0"
	if dump {
		helpers.DumpRequest(req, logger)
	}
	return req
}

func MakePageRequest(page, args string, meta ogdata.MetaData, cookieHeader []string) *http.Request {
	req, _ := http.NewRequest("GET", "http://"+meta.Uni+"-"+meta.Lang+".ogame.gameforge.com/game/index.php?page="+page+args, nil)
	// set up basic header
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header["Cookie"] = cookieHeader
	return req
}
