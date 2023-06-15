package main

import (
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"
	"unicode"

	"github.com/julienschmidt/httprouter"
)

func logNow(request *http.Request, timer time.Time) {
	logString := fmt.Sprintf("%s, %s, %s, %s, %s, took %s", request.UserAgent(), request.RemoteAddr, request.Method, request.Proto, request.URL, time.Since(timer))

	if strings.Contains(strings.ToLower(logString), "bot") || strings.Contains(strings.ToLower(logString), "spider") {
		return
	}

	log.Println(logString)
}

func setHeaders(writer http.ResponseWriter) http.ResponseWriter {
	writer.Header().Set("Content-Type", "text/html; charset=utf-8")
	writer.Header().Set("Cache-Control", "public, max-age=3600")
	writer.Header().Set("X-Powered-By", "Diet Coke")
	writer.Header().Set("ETag", ucVersion)

	return writer
}

func serveFilesFromTemplate(writer http.ResponseWriter, request *http.Request, templates []string, data interface{}, timer time.Time) {
	tmpl, err := template.New("").ParseFiles(templates...)
	if err != nil {
		log.Print(err.Error())
		http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	err = tmpl.ExecuteTemplate(writer, "base", data)
	if err != nil {
		log.Print(err.Error())
		http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	logNow(request, timer)
}

func redirectToTLS(writer http.ResponseWriter, request *http.Request) {
	log.Printf("Redirecting, %v, %v\n", request.RemoteAddr, request.URL)
	http.Redirect(writer, request, "https://unicode.click:443"+request.RequestURI, http.StatusMovedPermanently)
}

func serveIndex(writer http.ResponseWriter, request *http.Request, timer time.Time) {
	templateFiles := []string{
		"./template/base.template.html",
		"./template/index.template.html",
	}
	serveFilesFromTemplate(writer, request, templateFiles, struct{ UnicodeVersion string }{UnicodeVersion: unicode.Version}, timer)
}

func getRandomRune(maximum int) rune {
	// a much better way to do this would
	// be to randomly pick a table, then
	// randomly pick a rune within it...
	// I'm way too lazy to write that tho

	if maximum == 0 {
		// highest assigned
		maximum = 1114112
	}
	randy := rune(rand.Intn(maximum))

	// filter out hanja
	// otherwise its literally always hanja
	if unicode.In(randy, unicode.Han) {
		if randy%6 == 0 {
			return randy
		}
		return getRandomRune(maximum / 2)
	}

	if unicode.IsPrint(randy) {
		return randy
	}
	return getRandomRune(maximum / 2)
}

func serveRandom(writer http.ResponseWriter, request *http.Request, timer time.Time) {
	random := getRandomRune(0)
	http.Redirect(writer, request, "https://unicode.click:443/cp/"+string(random), http.StatusSeeOther)
	logNow(request, timer)
}

func serveUnicodeClick(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	timer := time.Now()
	route := params.ByName("filepath")

	// INFO: routes contain leading slash

	switch {
	case route == "/":
		serveIndex(writer, request, timer)
		return
	case route == "/random":
		serveRandom(writer, request, timer)
		return
	case len(route) >= 8 && route[0:6] == "/range":
		route = route[7:]
		serveRange(writer, request, route, timer)
		return
	case len(route) >= 5 && route[0:3] == "/cp":
		route = route[4:]
		serveCodepoint(writer, request, route, timer)
		return
	}

	route = "public" + route

	http.ServeFile(writer, request, route)
	logNow(request, timer)
}
