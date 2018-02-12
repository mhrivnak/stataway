package gloc

import (
	"bytes"
	"errors"
	"fmt"
	"golang.org/x/net/html"
	"golang.org/x/net/publicsuffix"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"
)

const loginURL = "https://accounts.google.com/ServiceLogin"
const emailURL = "https://accounts.google.com/signin/v1/lookup"
const passwordURL = "https://accounts.google.com/signin/challenge/sl/password"
const locationURL = "https://www.google.com/maps/preview/locationsharing/read?gl=en&pb=%211m7%218m6%211m3%211i14%212i8413%213i5385%212i6%213x4095%212m3%211e0%212sm%213i407105169%213m7%212sen%215e1105%2112m4%211e68%212m2%211sset%212sRoadmap%214e1%215m4%211e4%218m2%211e0%211e1%216m9%211e12%212i2%2126m1%214b1%2130m1%211f1.3953487873077393%2139b1%2144e1%2150e0%2123i4111425&authuser=0&hl=en"

func Get(username, password string) error {
	client, err := login(username, password)
	if err != nil {
		return err
	}

	resp, err := client.Get(locationURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var hbuf bytes.Buffer
	_ = resp.Header.Write(&hbuf)
	fmt.Println(hbuf.String())

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	fmt.Println(string(b))

	return nil
}

func newClient() (*http.Client, error) {
	options := cookiejar.Options{
		PublicSuffixList: publicsuffix.List,
	}
	jar, err := cookiejar.New(&options)
	if err != nil {
		return nil, err
	}

	client := &http.Client{
		// The default of 10 isn't enough for Google's login flow.
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= 15 {
				return errors.New("stopped after 15 redirects")
			}
			return nil
		},
		Jar: jar,
		// The default is no timeout, which is too optimistic even for me.
		Timeout: time.Second * 30,
	}
	return client, nil
}

func login(username, password string) (*http.Client, error) {
	client, err := newClient()
	if err != nil {
		return nil, err
	}

	respInital, err := client.Get(loginURL)
	if err != nil {
		return client, err
	}
	defer respInital.Body.Close()

	emailValues := findFormInputs(respInital.Body)
	emailValues.Set("Email", username)

	fmt.Println("### submitting email")

	respEmail, err := client.PostForm(emailURL, emailValues)
	if err != nil {
		return client, err
	}
	defer respEmail.Body.Close()

	fmt.Println("### submitting password")

	passwordValues := findFormInputs(respEmail.Body)
	passwordValues.Set("Passwd", password)

	respPassword, err := client.PostForm(passwordURL, passwordValues)
	if err != nil {
		return client, err
	}
	defer respPassword.Body.Close()

	fmt.Printf("Got status %s\n", respPassword.Status)

	return client, nil
}

func findFormInputs(body io.ReadCloser) url.Values {
	ret := url.Values{}

	z := html.NewTokenizer(body)
	inForm := false

	for {
		tt := z.Next()

		switch tt {
		case html.ErrorToken:
			return ret
		case html.StartTagToken:
			token := z.Token()
			if token.Data == "form" {
				fmt.Println(token.String())
				inForm = true
			} else if token.Data == "input" && inForm == true {
				var name string
				var value string
				for _, attribute := range token.Attr {
					if attribute.Key == "name" {
						name = attribute.Val
					} else if attribute.Key == "value" {
						value = attribute.Val
					}
				}
				ret.Set(name, value)
				name = ""
				value = ""
			}
		case html.EndTagToken:
			token := z.Token()
			if token.Data == "form" {
				fmt.Println("end of form")
				return ret
			}
		}
	}
}
