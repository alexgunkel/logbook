package logbook

import "net/http"

func DisplayLogBook(w http.ResponseWriter, r *http.Request)  {
	w.Header().Set("Content-type", "application/json")
	cookie := &http.Cookie{Name: "name", Value: "value"}
	http.SetCookie(w, cookie)
}
