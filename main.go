package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

var proto = "http"

func main() {
	http.HandleFunc("/", proxyHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func proxyHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if rProto := r.URL.Query().Get("proto"); rProto != "" {
		proto = rProto
	}
	host := r.URL.Query().Get("cors-host")
	url := fmt.Sprintf("%s://%s%s", proto, host, r.URL)

	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := client.Get(url)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Error")
		return
	}
	defer resp.Body.Close()

	byteBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Error")
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(byteBody))
}
