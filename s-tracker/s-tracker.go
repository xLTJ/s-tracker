package s_tracker

import (
	"fmt"
	"github.com/gocolly/colly"
	"github.com/spf13/viper"
	"net/http"
	"regexp"
	"strings"
)

const (
	SBaseUrl    = "https://mit.s.dk"
	profilePath = "/studiebolig/home/profile/"
)

type Client struct {
	Collector   *colly.Collector
	applicantId string
	username    string
}

func (c Client) GetApplicantId() string {
	return c.applicantId
}

func NewClient() (Client, error) {
	collector := colly.NewCollector()
	applicantId, err := clientSignIn(collector)

	if err != nil {
		return Client{}, fmt.Errorf("error setting up collector: %v", err)
	}

	return Client{Collector: collector, applicantId: applicantId}, nil
}

func clientSignIn(collector *colly.Collector) (applicantId string, err error) {
	csrfToken := viper.GetString("auth.csrfToken")
	sessionId := viper.GetString("auth.sessionId")

	// for every response, save status code
	var resStatus int
	collector.OnResponse(func(r *colly.Response) {
		resStatus = r.StatusCode
	})

	// try to get the Applicant id for every HTML response. should be inside a script tag after login.
	collector.OnHTML("script", func(e *colly.HTMLElement) {
		scriptContent := e.Text

		if strings.Contains(scriptContent, "applicant_pk") {
			re := regexp.MustCompile(`applicant_pk:\s*(\d+)`)
			applicantId = re.FindStringSubmatch(scriptContent)[1]
		}
	})

	defer collector.OnHTMLDetach("script")

	// if existing token and sessionId exists, try using those first
	if csrfToken != "" && sessionId != "" {
		fmt.Println("Attempting to use saved token...")
		err := collector.SetCookies(SBaseUrl, []*http.Cookie{
			{Name: "csrftoken", Value: csrfToken},
			{Name: "sessionid", Value: sessionId},
		})
		if err != nil {
			return "", fmt.Errorf("error setting cookies: %v", err)
		}

		err = collector.Visit(SBaseUrl + profilePath)
		if applicantId == "" {
			return "", fmt.Errorf("unable to get applicantId")
		}

		if err == nil {
			fmt.Println("Saved token is valid")
			return applicantId, nil
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

	return applicantId, nil
}
