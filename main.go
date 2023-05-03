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

/* var allUnicodeRangeTable unicode.RangeTable = unicode.RangeTable{
	R16: []unicode.Range16{{
		Lo:     0,
		Hi:     0xFFFF,
		Stride: 1,
	}},

	R32: []unicode.Range32{{
		Lo:     0,
		Hi:     0xFFFFFFFF,
		Stride: 1,
	}},

	LatinOffset: 0,
} */

type table struct {
	name string
	rows []row
}

type row struct {
	name string
	row  []rune
}

type rangeData struct {
	RangeTableName string
	Tables         []table

	UnicodeVersion string

	NumberOfTables int
	TableLengths   []int

	TableLiteral string
}

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
