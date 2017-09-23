package main

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

var proto = "http"

type response struct {
	Resp *http.Response
	Err  error
}

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

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	body, err := corsGet(ctx, url)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Error")
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, body)
}

func corsGet(ctx context.Context, url string) (string, error) {
	tr := &http.Transport{}
	client := &http.Client{Transport: tr}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	ch := make(chan response, 1)
	go func() {
		resp, err := client.Do(req)
		ch <- response{Resp: resp, Err: err}
	}()

	select {
	case <-ctx.Done():
		tr.CancelRequest(req)
		<-ch
		log.Printf("cancel request %s\n", req.URL.RequestURI())
		return "", ctx.Err()
	case res := <-ch:
		if res.Err != nil {
			return "", res.Err
		}
		defer res.Resp.Body.Close()

		byteBody, err := ioutil.ReadAll(res.Resp.Body)
		if err != nil {
			return "", err
		}
		return string(byteBody), nil
	}
	return "", errors.New("not expected reached")
}
