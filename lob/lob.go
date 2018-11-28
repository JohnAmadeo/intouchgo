package lob

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
)

// https://lob.com/docs#letters_create
type LobSendLetterRequest struct {
	Color          bool              `json:"color"`
	MailType       string            `json:"mail_type"`
	From           string            `json:"from"` // Temporarily set to my address
	To             string            `json:"to"`   // Lob Address object ID
	File           string            `json:"file"` // HTML string of letter's layout
	MergeVariables map[string]string `json:"merge_variables"`
}

type LobLetterThumbnail struct {
	Small  string `json:"small"`
	Medium string `json:"medium"`
	Large  string `json:"large"`
}

// https://lob.com/docs#letters_object
// Selectively chose whoch fields from the letter object we care about
type LobSendLetterResponse struct {
	Id                   string               `json:"id"`
	Thumbnails           []LobLetterThumbnail `json:"thumbnails"`
	ExpectedDeliveryDate string               `json:"expected_delivery_date"`
}

const (
	USPSStandard               = "usps_standard"
	LobTestEnvironment         = "test"
	LobLiveEnvironment         = "live"
	LobInvalidEnvironmentError = "Lob API environment must either be specified as test or live"
	LobCreateLetterRoute       = "https://api.lob.com/v1/letters"
	LobTestAPIKey              = "LOB_TEST_API_KEY"
	LobLiveAPIKey              = "LOB_LIVE_API_KEY"
	LobInTouchTestAddressId    = "adr_4cd2f1452346231d"
	LobInTouchLiveAddressId    = "adr_abd50ec552ac4841"
	LobBaseAPI                 = "https://api.lob.com/v1/"
	LobUserAgent               = "intouch/1.0"
)

func GetInTouchAddress(lobEnvironment string) (string, error) {
	switch lobEnvironment {
	case LobTestEnvironment:
		return LobInTouchTestAddressId, nil
	case LobLiveEnvironment:
		return LobInTouchLiveAddressId, nil
	default:
		return "", errors.New(LobInvalidEnvironmentError)
	}
}

func GetLetterHTMLTemplate(letterText string) (string, error) {
	bytes, err := ioutil.ReadFile("./templates/letter.html") // just pass the file name
	if err != nil {
		return "", err
	}

	htmlString := string(bytes) // convert content to a 'string'

	// NOTE: THIS IS A HACK to get around Lob's restriction on merge variables
	// needing to be <500 characters long
	parts := strings.Split(htmlString, "{{text}}")
	htmlString = parts[0] + letterText + parts[1]

	fmt.Println(htmlString)

	return htmlString, nil
}

func LobDateToDBDate(date string) string {
	return strings.Join(strings.Split(date, "-"), "/")
}

// Post performs a POST request to the Lob API.
func Post(endpoint string, params map[string]string, returnValue interface{}, environment string) error {
	fullURL := LobBaseAPI + endpoint
	fmt.Println("Lob POST ", fullURL)

	fmt.Println(params)

	var body io.Reader
	if params != nil {
		form := url.Values(make(map[string][]string))
		for k, v := range params {
			form.Add(k, v)
		}
		bodyString := form.Encode()
		body = bytes.NewBuffer([]byte(bodyString))
	}

	req, err := http.NewRequest("POST", fullURL, body)
	if err != nil {
		return err
	}

	if body != nil {
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	}

	key, err := getAPIKey(environment)
	if err != nil {
		return err
	}

	req.SetBasicAuth(key, "")
	// req.Header.Add("Lob-Version", APIVersion)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("User-Agent", LobUserAgent)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		err = fmt.Errorf("Non-200 status code %d returned from %s with body %s", resp.StatusCode, fullURL, data)
		json.Unmarshal(data, returnValue) // try, anyway -- in case the caller wants error info
		return err
	}

	return json.Unmarshal(data, returnValue)
}

func getAPIKey(lobEnvironment string) (string, error) {
	switch lobEnvironment {
	case LobTestEnvironment:
		apiKey, ok := os.LookupEnv(LobTestAPIKey)
		if !ok {
			return "", errors.New("Test API Key doesn't exist as an environment variable")
		}
		return apiKey, nil
	case LobLiveEnvironment:
		apiKey, ok := os.LookupEnv(LobLiveAPIKey)
		if !ok {
			return "", errors.New("Live API Key doesn't exist as an environment variable")
		}
		return apiKey, nil
	default:
		return "", errors.New(LobInvalidEnvironmentError)
	}
}
