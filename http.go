package goauth

import "net/http"

func HttpGen(r *http.Request, w http.ResponseWriter) {
	if r.Method == "POST" {
		r.ParseForm()
		user, qr, err := NewUser(r.Form.Get("username"))
		if err != nil {
			http.Error(w, "error: "+err.Error(), 500)
			return
		}
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte("<h1>Success!</h1>"))
		w.Write([]byte("<p>QR Code:</p>"))
		w.Write([]byte(qr.HTML()))
		w.Write([]byte("<a href='/validate'>validate</a>"))
	}

}
