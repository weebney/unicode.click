package main

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/julienschmidt/httprouter"
)

var ucVersion = "0.2.1"

func init() {
	f, err := os.OpenFile("./log/"+fmt.Sprint(time.Now())+".log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(io.MultiWriter(f, os.Stdout))
	rand.Seed(int64(f.Fd()))
}

func main() {
	router := httprouter.New()
	router.GET("/*filepath", serveUnicodeClick)

	fmt.Println("unicode.click listening on 443")

	go func() {
		if err := http.ListenAndServe(":80", http.HandlerFunc(redirectToTLS)); err != nil {
			log.Fatalf("ListenAndServe error: %v", err)
		}
	}()

	log.Fatal(http.ListenAndServeTLS(":443", "./ssl/domain.cert.pem", "./ssl/private.key.pem", router))
}
