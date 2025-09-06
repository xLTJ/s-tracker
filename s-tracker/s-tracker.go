package s_tracker

import (
	"fmt"
	"github.com/gocolly/colly"
	"github.com/spf13/viper"
	"net/http"
)

const (
	SBaseUrl    = "https://mit.s.dk"
	profilePath = "/studiebolig/home/profile/"
)

type Client struct {
	Collector *colly.Collector
}

//func (c Client) GetAuthCookies() (csrfToken, sessionId string) {
//	return c.csrfToken, c.sessionId
//}

func NewClient() (Client, error) {
	collector := colly.NewCollector()
	err := setupCollector(collector)
	if err != nil {
		return Client{}, fmt.Errorf("error setting up collector: %v", err)
	}

	return Client{collector}, nil
}

func setupCollector(collector *colly.Collector) error {
	csrfToken := viper.GetString("auth.csrfToken")
	sessionId := viper.GetString("auth.sessionId")

	// for every response, save status code
	var resStatus int
	collector.OnResponse(func(r *colly.Response) {
		resStatus = r.StatusCode
	})

	// if existing token and sessionId exists, try using those first
	if csrfToken != "" && sessionId != "" {
		fmt.Println("Attempting to use saved tokens...")
		err := collector.SetCookies(SBaseUrl, []*http.Cookie{
			{Name: "csrftoken", Value: csrfToken},
			{Name: "sessionid", Value: sessionId},
		})
		if err != nil {
			return fmt.Errorf("error setting cookies: %v", err)
		}

		err = collector.Visit(SBaseUrl + profilePath)
		if err == nil {
			return nil
		}
		fmt.Println("Saved token invalid, you need to log in again")
	}

	for resStatus != 200 {
		err := SLogin(collector)
		if err != nil {
			fmt.Printf("Login attempt failed: %v\n", err)
			continue
		}

		err = collector.Visit(SBaseUrl + profilePath)
		if resStatus != 200 {
			fmt.Printf("Login success, but token is still invalid for some reason :/")
		}
	}

	fmt.Println("Log in success")

	return nil
}
