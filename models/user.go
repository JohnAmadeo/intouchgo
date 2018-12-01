package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/google/go-querystring/query"
)

const (
	AuthConnection = "Username-Password-Authentication"
	ClientId       = "UiMO3i34HawDk03M2D7hpu4A2fhJoIoh"
	Domain         = "intouch-android.auth0.com"
)

type GetAccessTokenRequest struct {
	GrantType    string `json:"grant_type"`
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Audience     string `json:"audience"`
}

type GetAccessTokenResponse struct {
	AccessToken string  `json:"access_token"`
	ExpiresIn   float64 `json:"expires_in"`
	Scope       string  `json:"scope"`
	TokenType   string  `json:"token_type"`
}

type CreateUserRequest struct {
	Connection  string `json:"connection"`
	Username    string `json:"username"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	VerifyEmail bool   `json:"verify_email"`
}

type GetUserRealNameParams struct {
	SearchEngine string `url:"search_engine"`
	Query        string `url:"q"`
}

type Auth0User struct {
	Email        string          `json:"email"`
	Username     string          `json:"username"`
	UserMetadata Auth0UserMedata `json:"user_metadata"`
}

type Auth0UserMedata struct {
	Name string `json:"name"`
}

type User struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"placeholderPassword"`
}

func GetUserRealName(username string) (string, error) {
	accessToken, err := getManagementAcessToken()
	if err != nil {
		return "", err
	}

	url := "https://" + Domain + "/api/v2/users"
	values, err := query.Values(GetUserRealNameParams{
		SearchEngine: "v3",
		Query:        `username:"` + username + `"`,
	})
	if err != nil {
		return "", err
	}

	params := values.Encode()

	request, err := http.NewRequest("GET", url+"?"+params, nil)
	if err != nil {
		return "", err
	}

	request.Header.Add("content-type", "application/json")
	request.Header.Add("authorization", "Bearer "+accessToken)

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return "", err
	}

	if response.StatusCode != 200 {
		bytes, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return "", err
		}

		fmt.Printf("%d, %s", response.StatusCode, string(bytes))
		return "", errors.New("Status Code: " + string(response.StatusCode))
	}

	bytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	var users []Auth0User
	err = json.Unmarshal(bytes, &users)
	if err != nil {
		return "", err
	}

	if len(users) == 0 {
		return "", errors.New("No users with given username found")
	}

	return users[0].UserMetadata.Name, nil
}

func CreateUser(user User) error {
	accessToken, err := getManagementAcessToken()
	if err != nil {
		return err
	}

	url := "https://" + Domain + "/api/v2/users"
	body := CreateUserRequest{
		Connection:  AuthConnection,
		Username:    user.Username,
		Email:       user.Email,
		Password:    user.Password,
		VerifyEmail: false,
	}

	bytes, err := json.Marshal(body)
	if err != nil {
		return err
	}
	payload := strings.NewReader(string(bytes))

	request, err := http.NewRequest("POST", url, payload)
	if err != nil {
		return err
	}
	request.Header.Add("content-type", "application/json")
	request.Header.Add("authorization", "Bearer "+accessToken)

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return err
	}

	if response.StatusCode != 201 {
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return err
		}

		return errors.New(
			"Status Code: " + string(response.StatusCode) +
				"\nBody: " + string(body),
		)
	}

	return nil
}

func getManagementAcessToken() (string, error) {
	url := "https://" + Domain + "/oauth/token"
	secret, ok := os.LookupEnv("AUTH0_INTOUCH_CLIENT_SECRET")
	if !ok {
		return "", errors.New("Client secret doesn't exist as an environment variable")
	}

	request := GetAccessTokenRequest{
		GrantType:    "client_credentials",
		ClientId:     ClientId,
		ClientSecret: secret,
		Audience:     "https://" + Domain + "/api/v2/",
	}

	bytes, err := json.Marshal(request)
	if err != nil {
		return "", err
	}

	payload := strings.NewReader(string(bytes))

	response, err := http.Post(url, "application/json", payload)
	if err != nil {
		return "", err
	}

	bytes, err = ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	var responseBody GetAccessTokenResponse
	err = json.Unmarshal(bytes, &responseBody)

	if err != nil {
		return "", err
	}

	return responseBody.AccessToken, nil
}
