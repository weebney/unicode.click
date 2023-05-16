package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/julienschmidt/httprouter"
)

func main() {
	f, err := os.OpenFile("./log/"+fmt.Sprint(time.Now())+".log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	logWriter := io.MultiWriter(f, os.Stdout)
	log.SetOutput(logWriter)

	router := httprouter.New()
	router.GET("/favicon.ico", serveFavicon)
	router.GET("/", serveIndex)
	router.GET("/cp/:cpRoute", codepointFromRoute)
	router.GET("/range/:rangeRoute", rangeFromRoute)
	router.GET("/sitemaps/:route", serveSitemap)

	fmt.Println("unicode.click listening on 443")

	go func() {
		if err := http.ListenAndServe(":80", http.HandlerFunc(redirectToTls)); err != nil {
			log.Fatalf("ListenAndServe error: %v", err)
		}
	}()

	log.Fatal(http.ListenAndServeTLS(":443", "./ssl/domain.cert.pem", "./ssl/private.key.pem", router))
}
