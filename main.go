package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	host = "https://api.github.com"
)

func main() {
	http.HandleFunc("/", proxyHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func proxyHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	url := fmt.Sprintf("%s%s", host, r.URL)
	resp, err := http.Get(url)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	byteBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(byteBody))
}
