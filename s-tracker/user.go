package s_tracker

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	applicantInfoApiPath = "/api/applicant/"
)

type applicant struct {
	UserInfo User `json:"user"`
}

type User struct {
	Username  string `json:"username"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

func (c Client) GetUserInfo() (User, error) {
	client := &http.Client{}

	cookies := c.Collector.Cookies(SBaseUrl)
	req, _ := http.NewRequest("GET", SBaseUrl+applicantInfoApiPath+"/"+c.applicantId, nil)

	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}

	res, err := client.Do(req)
	if err != nil {
		return User{}, fmt.Errorf("error getting response: %v", err)
	}

	defer func() { _ = res.Body.Close() }()
	if res.StatusCode != 200 {
		return User{}, fmt.Errorf("unable to call api: %v", err)
	}

	var applicantInfo applicant
	err = json.NewDecoder(res.Body).Decode(&applicantInfo)
	if err != nil {
		return User{}, fmt.Errorf("error decoding response: %v", err)
	}

	return applicantInfo.UserInfo, nil
}
