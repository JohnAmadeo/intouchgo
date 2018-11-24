package routes

import (
	"encoding/json"
	"net/http"

	"github.com/johnamadeo/intouchgo/models"
	"github.com/johnamadeo/intouchgo/utils"
)

/*
curl -X GET -H "Content-Type: application/json" -H "Authorization: Bearer eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiIsImtpZCI6IlJqVTJRVEZFTTBVeU5EazFRalk1TmtFM05UazVOak0xUVVJeVFUWTRPVEZGUVVJeFJEY3lOdyJ9.eyJpc3MiOiJodHRwczovL2ludG91Y2gtYW5kcm9pZC5hdXRoMC5jb20vIiwic3ViIjoiYXV0aDB8NWJkODk1MWEzZGRjYjQwNWQzYWM2Y2RjIiwiYXVkIjpbImh0dHBzOi8vaW50b3VjaC1hbmRyb2lkLWJhY2tlbmQuaGVyb2t1YXBwLmNvbS8iLCJodHRwczovL2ludG91Y2gtYW5kcm9pZC5hdXRoMC5jb20vdXNlcmluZm8iXSwiaWF0IjoxNTQzMDA1MTU5LCJleHAiOjE1NDMwOTE1NTksImF6cCI6InpjVU54OGxQQVE2djhVSXB4OTIwVkdvVTVnMmplNXl6Iiwic2NvcGUiOiJvcGVuaWQgcHJvZmlsZSBlbWFpbCBvZmZsaW5lX2FjY2VzcyJ9.nwRklfXGeQzRBBbh20Qh2MTBYJWPIn4fPVv8BRqetPN5Gn1X2eVXL3t9XgzzHst7Y50R7temxmMFubFYwNc-OvKWKPlrjVa8sylORwMzN04MzLYlr8LDwYGP3uUBGpys_1LpdBGUXxP1JVMGleK24TtXoYyaVPyBAi-Qr9BFFNUm8Kw2jBnDaj-lFuwu8IIaKAYOwKVRAEcR1fe62shum5LOxc-NpsNNp4JQZp1ajlxjAVaZifd0BuFTLoGB13L2sIpp7AGxVpSX_YNed4LoGE2xU_UQxmfl6_hZypdzTnp1nkDLrXEtjGzAbDDoXboDfDUawEStNlfUSdciJgegLg" http://localhost:8080/inmates
*/
func InmatesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" && r.Method != "" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write(utils.MessageToBytes("Only GET requests are allowed at this route"))
		return
	}

	queries, ok := r.URL.Query()["query"]
	if !ok || len(queries) > 1 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(utils.MessageToBytes("Request query parameters must contain a single username"))
		return
	}

	inmates, err := models.GetInmatesFromDB(queries[0])
	if err != nil {
		utils.PrintErr(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(utils.MessageToBytes(err.Error()))
		return
	}

	bytes, err := json.Marshal(inmates)
	if err != nil {
		utils.PrintErr(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(utils.MessageToBytes(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(bytes)
}
