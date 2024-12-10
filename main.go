package main

import (
	"crypto/sha1"
	"encoding/base64"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
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
		err := r.ParseForm()
		if err != nil {
			log.Fatal("Errore:", err)
		}

		url := r.PostForm.Get("url")
		if url == "" {
			log.Fatal("Il valore del parametro url non è corretto o è vuoto")
		}

		shortUrl := shortenURL(url)

		db.urls[shortUrl] = url

		fmt.Fprintf(w, "http://localhost%v/short/%v\n", port, shortUrl)

	}
}

func handleRedirect(db *ShortUrls) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		shortUrl := strings.SplitAfter(r.URL.Path, "/short/")[1]
		if shortUrl == "" {
			log.Fatal("Short url non presente")
		}

		http.Redirect(w, r, db.urls[shortUrl], http.StatusMovedPermanently)
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
	http.HandleFunc("/short/", handleRedirect(&db))

	fmt.Println("server start on ", PORT)
	http.ListenAndServe(PORT, nil)
}
