package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

const (
	chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	page = `
    <!DOCTYPE html>
    <html>
        <head><title>url shortener</title></head>
        <body>
			<h2>url shortener</h2>
			<form method = "POST" action = "/register">
				<input type = "url" name = "url" placeholder = "enter your url" required>
				<button type = "submit">submit</button>
			</form>
			<p>%s</p>
        </body>
    </html>`
)

var (
	cnt uint64 = 0
	ctu        = make(map[string]string)
	utc        = make(map[string]string)
)

func insert(url string) string {
	cur := cnt
	cnt++
	var code string
	for i := 0; i < 5; i++ {
		rem := cur % 62
		code = code + string(chars[rem])
		cur /= 62
	}
	ctu[code] = url
	utc[url] = code
	return code
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	var code string
	if r.Method == http.MethodPost {
		url := r.FormValue("url")
		val, found := utc[url]
		if found {
			code = val
		} else {
			code = insert(url)
		}
		code = "your short link: http://" + r.Host + "/" + code
	}
	fmt.Fprintf(w, page, code)
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	code := ""
	for i := 1; i < len(path); i++ {
		code = code + string(path[i])
	}
	if code == "" {
		http.Redirect(w, r, "/register", http.StatusFound)
		return
	}
	val, found := ctu[code]
	if found {
		http.Redirect(w, r, val, http.StatusFound)
		return
	}
	http.Error(w, "URL not registered", http.StatusBadRequest)
}

func main() {
	http.HandleFunc("/register", registerHandler)
	http.HandleFunc("/", redirectHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Printf("Server running on port %s", port)
	http.ListenAndServe(":"+port, nil)
}
