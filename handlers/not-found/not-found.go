package notFound

import (
	"encoding/json"
	"ibn-salamat/simple-chat-api/types"
	"net/http"
)

func NotFound(w http.ResponseWriter, _ *http.Request) {
	jsonBody, _ := json.Marshal(types.ResponseMap{
		"errorMessage": "Not found",
	})

	w.WriteHeader(http.StatusNotFound)
	w.Write(jsonBody)
}
