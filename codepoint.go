package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/julienschmidt/httprouter"
	"golang.org/x/text/unicode/runenames"
)

func codepointFromRoute(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	timer := time.Now()
	route := params.ByName("cpRoute")

	writer.Header().Set("Content-Type", "text/html; charset=utf-8")
	writer.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	writer.Header().Set("Pragma", "no-cache")
	writer.Header().Set("Expires", "0")

	var codepoint rune

	if strings.ContainsRune(route, '+') {
		// convert UTF-16 w/ U+ prefix (i.e. +0061) to rune type
		tempRoute := strings.SplitAfter(route, "+")
		codepointInt64, _ := strconv.ParseInt(tempRoute[1], 16, 0)
		codepoint = int32(codepointInt64)
	} else {
		// convert literal character to rune
		runeArray := []rune(route)
		codepoint = int32(runeArray[0])
	}

	// check if codepoint exists, do something else if not
	if codepoint > unicode.MaxRune || codepoint < 0 || codepoint > 2147483647 {
		http.Error(writer, "Not Found", 404)
		return
	}

	majorCategoryLiteral, categoryLiteral, categories, majorCategories := getCategoryData(codepoint)

	var scripts []string
	for scriptName, scriptRangeTable := range unicode.Scripts {
		if unicode.Is(scriptRangeTable, codepoint) {
			scripts = append(scripts, scriptName)
		}
	}

	var properties []string
	for propertyName, propertyRangeTable := range unicode.Properties {
		if unicode.Is(propertyRangeTable, codepoint) {
			properties = append(properties, propertyName)
		}
	}

	var data = struct {
		CodepointHexAsString string
		LitRune              string
		RuneName             string
		UnicodeVersion       string

		Scripts    string
		Properties []string

		MajorCategories string
		Categories      string
		MajCatLiteral   string
		CatLiteral      string

		IsControl bool
		IsDigit   bool
		IsGraphic bool
		IsLetter  bool
		IsLower   bool
		IsMark    bool
		IsNumber  bool
		IsPrint   bool
		IsPunct   bool
		IsSpace   bool
		IsSymbol  bool
		IsTitle   bool
		IsUpper   bool

		AsUppercase string
		AsLowercase string
		AsTitlecase string

		HasDifferentCase bool
	}{
		CodepointHexAsString: fmt.Sprintf("%U", codepoint),
		LitRune:              string(codepoint),
		RuneName:             runenames.Name(codepoint),
		UnicodeVersion:       unicode.Version,

		Scripts:    strings.Join(scripts, ", "),
		Properties: properties,

		MajCatLiteral:   majorCategoryLiteral,
		CatLiteral:      categoryLiteral,
		MajorCategories: strings.Join(majorCategories, ", "),
		Categories:      strings.Join(categories, ", "),

		IsControl: unicode.IsControl(codepoint),
		IsDigit:   unicode.IsDigit(codepoint),
		IsGraphic: unicode.IsGraphic(codepoint),
		IsLetter:  unicode.IsLetter(codepoint),
		IsLower:   unicode.IsLower(codepoint),
		IsMark:    unicode.IsMark(codepoint),
		IsNumber:  unicode.IsNumber(codepoint),
		IsPrint:   unicode.IsPrint(codepoint),
		IsPunct:   unicode.IsPunct(codepoint),
		IsSpace:   unicode.IsSpace(codepoint),
		IsSymbol:  unicode.IsSymbol(codepoint),
		IsTitle:   unicode.IsTitle(codepoint),
		IsUpper:   unicode.IsUpper(codepoint),

		AsUppercase: string(unicode.ToUpper(codepoint)),
		AsLowercase: string(unicode.ToLower(codepoint)),
		AsTitlecase: string(unicode.ToTitle(codepoint)),

		HasDifferentCase: (unicode.IsUpper(codepoint) || unicode.IsLower(codepoint)) || unicode.IsTitle(codepoint),
	}

	templateFiles := []string{
		"./template/base.template.html",
		"./template/rune.template.html",
	}

	serveFilesFromTemplate(writer, request, params, templateFiles, data, timer)
}
