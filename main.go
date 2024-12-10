package main

import (
	"crypto/sha1"
	"encoding/base64"
	"flag"
	"fmt"
	"log"
	"net/http"
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

		// Here the url should be passed as a JSON object, for example:
		// curl -X POST -d '{"url": "http://www.google.it"}' http://localhost:8080/create
		// rather than via form data.

		err := r.ParseForm()
		if err != nil {
			// Doing log.Fatal means that any client can crash the server by
			// just sending a malformed request. You should log the error, send
			// a reasoannable http status to the client and return to end the
			// request.
			log.Fatal("Errore:", err)
		}

		url := r.PostForm.Get("url")
		if url == "" {
			// same here.
			log.Fatal("Il valore del parametro url non è corretto o è vuoto")
		}

		// before proceeding with storing the URL, we should check the URL is
		// valid, for example with:
		//
		// _, err := url.Parse(url)
		//
		// and return the appropriate status code if the URL is not valid.
		shortUrl := shortenURL(url)

		db.urls[shortUrl] = url

		fmt.Fprintf(w, "http://localhost%v/short/%v\n", port, shortUrl)

	}
}

func handleRedirect(db *ShortUrls) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		shortcode := r.PathValue("shortcode")

		// shortUrl := strings.SplitAfter(r.URL.Path, "/short/")[1]
		// if shortUrl == "" {
		// 	// same here about the use of log.Fatal
		// 	log.Fatal("Short url non presente")
		// }

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
	// Here you could use "/short/{code}", this has the advantage of only
	// entering in the handler if code is non-empty.
	//
	// If you do this, then in the handler you can directly use:
	//  r.PathValue("code")
	// to access the value of 'code'.
	http.HandleFunc("/short/{code}", handleRedirect(&db))

	fmt.Println("server start on ", PORT)
	http.ListenAndServe(PORT, nil)
}
