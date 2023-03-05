package main

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"
)

var tCount = 1

func t(name string, test func() error) {
	fmt.Println("#", tCount, "-", name)
	err := test()
	i := "ok"
	if err != nil {
		i = fmt.Sprint(err)
		fmt.Println("#", tCount, "-", name, "error:", i)
	}
	if err != nil {
		os.Exit(1)
	}
	tCount++
}

func httpctx(method, rurl string, values url.Values) (http.ResponseWriter, *http.Request) {
	purl, _ := url.ParseRequestURI(rurl)
	w := resp{}
	r := &http.Request{
		Method:     method,
		URL:        purl,
		Form:       values,
		PostForm:   values,
		RequestURI: rurl,
	}
	return w, r
}

type resp struct{}

func (r resp) Header() http.Header {
	return http.Header{}
}
func (r resp) Write([]byte) (int, error) {
	return 0, nil
}
func (r resp) WriteHeader(int) {
	return
}

func tDate(month, day, hour, min int) time.Time {
	location, _ := time.LoadLocation("Europe/Moscow")
	return time.Date(2021, time.Month(month), day, hour, min, 0, 0, location)
}
