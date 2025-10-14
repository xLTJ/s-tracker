package s_tracker

import (
	"encoding/json"
	"fmt"
	"github.com/gocolly/colly"
)

const (
	applicantInfoApiPath = "api/applicant"
)

type Applicant struct {
	UserInfo User `json:"user"`
}

type User struct {
	Username  string `json:"username"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

func (c Client) GetUserInfo() (User, error) {
	var err error
	var user User
	c.Collector.OnResponse(func(r *colly.Response) {
		var applicant Applicant
		err = json.Unmarshal(r.Body, &applicant)
		user = applicant.UserInfo
	})

	_ = c.Collector.Visit(fmt.Sprintf("%s/%s/%s", SBaseUrl, applicantInfoApiPath, c.applicantId))
	if err != nil {
		return User{}, fmt.Errorf("error decoding response: %v", err)
	}

	return user, nil
}
