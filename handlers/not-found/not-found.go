package notFound

import "net/http"

func NotFound(w http.ResponseWriter)  {
	w.Header().Set("Content-Type", "application/json")
	
	
}