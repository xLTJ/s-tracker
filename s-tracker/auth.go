package s_tracker

import (
	"fmt"
	"github.com/gocolly/colly"
	"github.com/pterm/pterm"
	"github.com/spf13/viper"
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

	err = collector.Post(SBaseUrl+loginUrl, map[string]string{
		"csrfmiddlewaretoken": csrfMiddlewareToken,
		"username":            username,
		"password":            password,
	})

	if err != nil {
		return fmt.Errorf("error loggin in: %v", err)
	}

	cookies := collector.Cookies(SBaseUrl)
	if len(cookies) != 2 {
		return fmt.Errorf("invalid credentials, try again")
	}

	for _, cookie := range cookies {
		viper.Set("auth."+cookie.Name, cookie.Value)
	}
	err = viper.WriteConfig()
	if err != nil {
		fmt.Println("Unable to write to config, tokens are not saved")
	}

	return nil
}

func inputCredentials() (username, password string, err error) {
	pterm.FgGreen.Println("Log in with your s.dk credentials")
	username, _ = pterm.DefaultInteractiveTextInput.Show("Enter Username")
	pterm.Println()

	password, _ = pterm.DefaultInteractiveTextInput.Show("Enter Password")
	return
}
