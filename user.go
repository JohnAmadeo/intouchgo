package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

const (
	Domain   = "intouch-android.auth0.com"
	ClientId = "UiMO3i34HawDk03M2D7hpu4A2fhJoIoh"
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

type User struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"placeholderPassword"`
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

func createUser(accessToken string, user User) error {
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

func createUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write(messageToBytes("Only POST requests are allowed at this route"))
		return
	}

	var user User
	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(messageToBytes("Malformed body."))
		return
	}
	defer r.Body.Close()

	err = json.Unmarshal(bytes, &user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(messageToBytes("Request body must be a user."))
		return
	}

	accessToken, err := getManagementAcessToken()
	if err != nil {
		printErr(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(messageToBytes("Failed to get Auth0 Management API access token to create user."))
		return
	}

	err = createUser(accessToken, user)
	if err != nil {
		printErr(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(messageToBytes("Failed to create user: " + err.Error()))
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(messageToBytes("Successfully created user."))
}
