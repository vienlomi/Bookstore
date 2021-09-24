package helperApis

import "net/http"

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "not found page", http.StatusNotFound)
}

func NotAllowHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "not allow request", http.StatusMethodNotAllowed)
}
