package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
	"unicode"

	"golang.org/x/text/unicode/runenames"
)

func serveCodepoint(writer http.ResponseWriter, request *http.Request, route string, timer time.Time) {
	writer = setHeaders(writer)

	var codepoint rune

	if strings.ContainsRune(route, '+') {
		// convert from U+ prefix (i.e. U+0061) to rune
		tempRoute := strings.SplitAfter(route, "+")
		codepointInt64, _ := strconv.ParseInt(tempRoute[1], 16, 0)
		codepoint = rune(codepointInt64)
	} else {
		// convert first character to rune
		runeArray := []rune(route)
		codepoint = rune(runeArray[0])
	}

	// check if codepoint exists, do something else if not
	if codepoint > unicode.MaxRune || codepoint < 0 || codepoint > 2147483647 {
		// TODO: create a dedicated 404 page with a JS-based automatic redirect
		http.Redirect(writer, request, "https://unicode.click/", http.StatusMovedPermanently)
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

	data := struct {
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

	serveFilesFromTemplate(writer, request, templateFiles, data, timer)
}
