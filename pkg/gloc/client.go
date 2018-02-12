package gloc

import (
	"errors"
	"fmt"
	"golang.org/x/net/html"
	"golang.org/x/net/publicsuffix"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"
)

const loginURL = "https://accounts.google.com/ServiceLogin"
const emailURL = "https://accounts.google.com/signin/v1/lookup"
const passwordURL = "https://accounts.google.com/signin/challenge/sl/password"

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
