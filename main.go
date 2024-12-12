package main

import (
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

type ShortUrls struct {
	urls map[string]string
}

func shortenURL(url string) string {
	h := sha1.New()
	h.Write([]byte(url))
	bs := h.Sum(nil)
	return base64.URLEncoding.EncodeToString(bs)[:6]
}

func handleGenerateShortUrl(db *ShortUrls, port string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var u struct {
			Url string `json:"url"`
		}

		err := json.NewDecoder(r.Body).Decode(&u)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if u.Url == "" {
			http.Error(w, "Il valore del parametro url non è corretto o è vuoto", http.StatusBadRequest)
			return
		}

		if _,err := url.Parse(u.Url); err != nil {
			http.Error(w,err.Error(),http.StatusBadRequest)
		}

		shortUrl := shortenURL(u.Url)

		db.urls[shortUrl] = u.Url

		fmt.Fprintf(w, "http://localhost%v/short/%v\n", port, shortUrl)

	}
}

func handleRedirect(db *ShortUrls) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		shortcode := r.PathValue("code")

		http.Redirect(w, r, db.urls[shortcode], http.StatusMovedPermanently)
	}
}

func main() {
	var portNum int
	flag.IntVar(&portNum, "p", 8080, "Definisce il numero di porta sulla quale dovrà girare il server, di default è 8080")
	flag.Parse()

	var db ShortUrls
	db.urls = make(map[string]string)

	PORT := ":" + strconv.Itoa(portNum)

	http.HandleFunc("POST /create", handleGenerateShortUrl(&db, PORT))
	http.HandleFunc("/short/{code}", handleRedirect(&db))

	fmt.Println("server start on ", PORT)
	http.ListenAndServe(PORT, nil)
}
