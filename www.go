package main

import (
	"html/template"
	"log"
	"net/http"
	"time"
	"unicode"

	"github.com/julienschmidt/httprouter"
)

func serveFilesFromTemplate(writer http.ResponseWriter, request *http.Request, params httprouter.Params, templates []string, data interface{}, timer time.Time) {

	// Parse the templates and add the function map
	tmpl, err := template.New("").ParseFiles(templates...)
	if err != nil {
		log.Print(err.Error())
		http.Error(writer, "Internal Server Error", 500)
		return
	}

	err = tmpl.ExecuteTemplate(writer, "base", data)
	if err != nil {
		log.Print(err.Error())
		http.Error(writer, "Internal Server Error", 500)
	}

	log.Printf("%s, %s, %s, %s, %s, %s\n", request.UserAgent(), request.RemoteAddr, request.Method, request.Proto, request.URL, time.Since(timer))
}

func serveFavicon(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	http.ServeFile(writer, request, "./favicon.ico")
}

func serveIndex(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	timer := time.Now()
	templateFiles := []string{
		"./template/base.template.html",
		"./template/index.template.html",
	}

	serveFilesFromTemplate(writer, request, params, templateFiles, struct{ UnicodeVersion string }{UnicodeVersion: unicode.Version}, timer)
}

func redirectToTls(writer http.ResponseWriter, request *http.Request) {
	log.Printf("Redirecting, %v, %v\n", request.RemoteAddr, request.URL)
	http.Redirect(writer, request, "https://unicode.click:443"+request.RequestURI, http.StatusMovedPermanently)

}

func serveSitemap(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	timer := time.Now()
	route := params.ByName("route")
	http.ServeFile(writer, request, "./sitemaps/"+route)
	log.Printf("Sitemap, %v, %v, %v, %v\n", request.UserAgent(), request.RemoteAddr, request.URL, time.Since(timer))
}
