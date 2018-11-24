package routes

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/johnamadeo/intouchgo/models"
	"github.com/johnamadeo/intouchgo/utils"
)

/*
curl -X POST -H "Content-Type: application/json" -H "Authorization: Bearer eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiIsImtpZCI6IlJqVTJRVEZFTTBVeU5EazFRalk1TmtFM05UazVOak0xUVVJeVFUWTRPVEZGUVVJeFJEY3lOdyJ9.eyJpc3MiOiJodHRwczovL2ludG91Y2gtYW5kcm9pZC5hdXRoMC5jb20vIiwic3ViIjoiYXV0aDB8NWJkODk1MWEzZGRjYjQwNWQzYWM2Y2RjIiwiYXVkIjpbImh0dHBzOi8vaW50b3VjaC1hbmRyb2lkLWJhY2tlbmQuaGVyb2t1YXBwLmNvbS8iLCJodHRwczovL2ludG91Y2gtYW5kcm9pZC5hdXRoMC5jb20vdXNlcmluZm8iXSwiaWF0IjoxNTQyNTA3MTM2LCJleHAiOjE1NDI1OTM1MzYsImF6cCI6InpjVU54OGxQQVE2djhVSXB4OTIwVkdvVTVnMmplNXl6Iiwic2NvcGUiOiJvcGVuaWQgcHJvZmlsZSBlbWFpbCBvZmZsaW5lX2FjY2VzcyJ9.kNFhzrMJ75jPjWufClOm8oXYTGIGEZHnmDGOV0Hu-T8r_ObtFhngsT1QoqI8-6Y2jIyKJyKUc4ySbfS9PRAeYaulntnFlvdcJXJK6M5_gNizSS84ROgb6vHVBxdZ6YfIIS3iQbf51g2xJs-jRWWwyb1_Shr8Fnt0fxOu_D_KBt82GiQWdJgePT-SVKMRrmZno_4aq07-YFxdJppLbfQIeMz4a4kApBHQ1ZKPZ6qrUcZ1xPKqd3dTMXrMtWAvkISGp31zTPEbm-40t4xIWgfHOXK20WJ-KZ4zV7nczmYHMNqdxx2ww5J9xc2vXah2ITPJagQUlQH6ZqrOsmGC4ZJcGA" -d '{"id":"eiwo-19da-p2gv", "author": "jadk157", "recipient": "John Grant", "recipientId": "asdf-123s-ddss", "subject": "Arden - 1st day at school", "text": "Arden went to his 1st day at Oakwood High this Monday! Mom and I dropped him off.", "timeSent": "10/08/18", "timeLastEdited": "10/08/18", "isDraft": false}' http://localhost:8080/letter
*/
func CreateLetterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write(utils.MessageToBytes("Only POST requests are allowed at this route"))
		return
	}

	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(utils.MessageToBytes("Malformed body."))
		return
	}
	defer r.Body.Close()

	var letter models.Letter
	err = json.Unmarshal(bytes, &letter)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(utils.MessageToBytes("Request body must be a letter"))
		return
	}

	err = models.CreateLetterInDB(letter)
	if err != nil {
		utils.PrintErr(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(utils.MessageToBytes("Error inserting to DB: " + err.Error()))
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(utils.MessageToBytes("Letter successfully created!"))
}

/*
curl -X GET -H "Content-Type: application/json" -H "Authorization: Bearer eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiIsImtpZCI6IlJqVTJRVEZFTTBVeU5EazFRalk1TmtFM05UazVOak0xUVVJeVFUWTRPVEZGUVVJeFJEY3lOdyJ9.eyJpc3MiOiJodHRwczovL2ludG91Y2gtYW5kcm9pZC5hdXRoMC5jb20vIiwic3ViIjoiYXV0aDB8NWJkODk1MWEzZGRjYjQwNWQzYWM2Y2RjIiwiYXVkIjpbImh0dHBzOi8vaW50b3VjaC1hbmRyb2lkLWJhY2tlbmQuaGVyb2t1YXBwLmNvbS8iLCJodHRwczovL2ludG91Y2gtYW5kcm9pZC5hdXRoMC5jb20vdXNlcmluZm8iXSwiaWF0IjoxNTQyNTA3MTM2LCJleHAiOjE1NDI1OTM1MzYsImF6cCI6InpjVU54OGxQQVE2djhVSXB4OTIwVkdvVTVnMmplNXl6Iiwic2NvcGUiOiJvcGVuaWQgcHJvZmlsZSBlbWFpbCBvZmZsaW5lX2FjY2VzcyJ9.kNFhzrMJ75jPjWufClOm8oXYTGIGEZHnmDGOV0Hu-T8r_ObtFhngsT1QoqI8-6Y2jIyKJyKUc4ySbfS9PRAeYaulntnFlvdcJXJK6M5_gNizSS84ROgb6vHVBxdZ6YfIIS3iQbf51g2xJs-jRWWwyb1_Shr8Fnt0fxOu_D_KBt82GiQWdJgePT-SVKMRrmZno_4aq07-YFxdJppLbfQIeMz4a4kApBHQ1ZKPZ6qrUcZ1xPKqd3dTMXrMtWAvkISGp31zTPEbm-40t4xIWgfHOXK20WJ-KZ4zV7nczmYHMNqdxx2ww5J9xc2vXah2ITPJagQUlQH6ZqrOsmGC4ZJcGA" http://localhost:8080/letters
*/
func LettersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" && r.Method != "" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write(utils.MessageToBytes("Only GET requests are allowed at this route"))
		return
	}

	usernames, ok := r.URL.Query()["username"]
	if !ok || len(usernames) > 1 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(utils.MessageToBytes("Request query parameters must contain a single username"))
		return
	}

	letters, err := models.GetLettersFromDB(usernames[0])
	if err != nil {
		utils.PrintErr(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(utils.MessageToBytes(err.Error()))
		return
	}

	bytes, err := json.Marshal(letters)
	if err != nil {
		utils.PrintErr(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(utils.MessageToBytes(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(bytes)
}
