package s_tracker

import (
	"bufio"
	"fmt"
	"github.com/gocolly/colly"
	"os"
	//"github.com/google/go-querystring/query"
	//"golang.org/x/net/html"
	//"net/http"
	//"net/http/cookiejar"
	//"net/url"
)

const loginUrl = "/studiebolig/login/"

func SLogin(collector *colly.Collector) error {
	username, password, err := inputCredentials()
	if err != nil {
		return err
	}

	// extract csrfMiddleWareToken from hidden inputs
	var csrfMiddlewareToken string
	collector.OnHTML("input[name='csrfmiddlewaretoken']", func(e *colly.HTMLElement) {
		csrfMiddlewareToken = e.Attr("value")
	})

	// visit login url and get csrf tokens. the main token is in cookies
	err = collector.Visit(SBaseUrl + loginUrl)
	if err != nil {
		return fmt.Errorf("error visiting login site: %v", err)
	}

	// set content type for login post-requests
	collector.OnRequest(func(r *colly.Request) {
		if r.URL.Path == loginUrl && r.Method == "POST" {
			r.Headers.Set("Content-Type", "application/x-www-form-urlencoded")
			r.Headers.Set("Referer", "https://mit.s.dk/studiebolig/login/")
		}
	})

	//collector.OnError(func(r *colly.Response, err error) {
	//	fmt.Printf("Error: %d %v URL: %s\n", r.StatusCode, err, r.Request.URL)
	//	fmt.Printf("Response Headers: %v\n", r.Headers)
	//	fmt.Printf("Response Body:\n%s\n", string(r.Body))
	//})

	err = collector.Post(SBaseUrl+loginUrl, map[string]string{
		"csrfmiddlewaretoken": csrfMiddlewareToken,
		"username":            username,
		"password":            password,
	})

	if err != nil {
		return fmt.Errorf("error loggin in: %v", err)
	}

	cookies := collector.Cookies(SBaseUrl)
	if len(cookies) < 2 {
		return fmt.Errorf("invalid credentials, try again")
	}

	return nil
}

func inputCredentials() (username, password string, err error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Log in with your s.dk credentials")
	fmt.Printf("Enter username: ")
	username, err = reader.ReadString('\n')

	if err != nil {
		return "", "", fmt.Errorf("error reading username: %v", err)
	}

	fmt.Printf("Enter password: ")
	password, err = reader.ReadString('\n')
	if err != nil {
		return "", "", fmt.Errorf("error reading username: %v", err)
	}

	return
}
