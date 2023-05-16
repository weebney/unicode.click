package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/julienschmidt/httprouter"
)

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

	TableLiteral template.HTML
}

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

func generateTableFromRTLiteral(rtLiteral *unicode.RangeTable) (tables []table, tableLengths []int) {
	var rawRuneSlice []int32

	for i := 0; i < len(rtLiteral.R16); i++ {
		rtLo := rtLiteral.R16[i].Lo
		rtHi := rtLiteral.R16[i].Hi
		for runeLiteral := rtLo; runeLiteral <= rtHi; runeLiteral += rtLiteral.R16[i].Stride {
			rawRuneSlice = append(rawRuneSlice, int32(runeLiteral))
		}
	}
	for i := 0; i < len(rtLiteral.R32); i++ {
		rtLo := rtLiteral.R32[i].Lo
		rtHi := rtLiteral.R32[i].Hi
		for runeLiteral := rtLo; runeLiteral <= rtHi; runeLiteral += rtLiteral.R32[i].Stride {
			rawRuneSlice = append(rawRuneSlice, int32(runeLiteral))
		}
	}

	rows := []row{}

	for i := 0; i < len(rawRuneSlice); i++ {
		currentRune := fmt.Sprintf("%U", rawRuneSlice[i])
		currentRune = currentRune[2:]

		if !unicode.In(rawRuneSlice[i], rtLiteral) {
			break
		}

		rowPrefix := currentRune[:len(currentRune)-1]
		if !((len(rows) != 0) && rows[len(rows)-1].name == rowPrefix) {
			rowInt, _ := strconv.ParseInt(rowPrefix+"0", 16, 32)
			var localRow []rune
			for i := 0; i < 16; i++ {
				localRow = append(localRow, rune(rowInt+int64(i)))
			}
			var assembledRow row = row{rowPrefix, localRow}
			rows = append(rows, assembledRow)
		}
	}

	// validate rows, pop them if invalid

	tables = []table{}
	for i := 0; i < len(rows); i++ {
		tablePrefix := rows[i].name[:len(rows[i].name)-1]
		rowSlice := []row{rows[i]}

		for j := i + 1; j < len(rows); j++ {
			if rows[j].name[:len(rows[j].name)-1] == tablePrefix {
				rowSlice = append(rowSlice, rows[j])
				i++
			}
		}

		var assembledTable table = table{tablePrefix, rowSlice}
		tables = append(tables, assembledTable)
	}

	for i := 0; i < len(tables); i++ {
		tableLengths = append(tableLengths, len(tables[i].rows))
	}

	return
}

func rangeFromRoute(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	timer := time.Now()
	route := strings.ToLower(params.ByName("rangeRoute"))

	writer.Header().Set("Content-Type", "text/html; charset=utf-8")
	writer.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	writer.Header().Set("Pragma", "no-cache")
	writer.Header().Set("Expires", "0")

	rtLiteral := getRangeTableLiteral(route)

	tables, tableLengths := generateTableFromRTLiteral(rtLiteral)

	var data = rangeData{
		RangeTableName: route,
		Tables:         tables,

		UnicodeVersion: unicode.Version,

		NumberOfTables: len(tables),
		TableLengths:   tableLengths,

		TableLiteral: "",
	}

	data.TableLiteral = generateTableHtml(data.Tables, data.TableLengths, rtLiteral)

	templateFiles := []string{
		"./template/base.template.html",
		"./template/range.template.html",
	}

	serveFilesFromTemplate(writer, request, params, templateFiles, data, timer)
}

func generateTableHtml(tables []table, tableLengths []int, literalRT *unicode.RangeTable) template.HTML {
	var literal []template.HTML

	literal = append(literal, "<div id=tables>")

	for i := 0; i < len(tables); i++ {
		literal = append(literal, `
		<table>
		<tr>
		<th></th>
        <th>0</th>
        <th>1</th>
        <th>2</th>
        <th>3</th>
        <th>4</th>
        <th>5</th>
        <th>6</th>
        <th>7</th>
        <th>8</th>
        <th>9</th>
        <th>A</th>
        <th>B</th>
        <th>C</th>
        <th>D</th>
        <th>E</th>
        <th>F</th></tr>
		`)
		for row := 0; row < tableLengths[i]; row++ {
			literal = append(literal, "<tr>")
			literal = append(literal, "<td>U+", template.HTML(tables[i].rows[row].name), "</td>")
			for j := 0; j < len(tables[i].rows[row].row); j++ {
				if unicode.In(tables[i].rows[row].row[j], literalRT) {
					literal = append(literal, `<td><a href="/cp/`, template.HTML(fmt.Sprintf("%U", tables[i].rows[row].row[j])), `">`, template.HTML(tables[i].rows[row].row[j]), "</a></td>")
				} else {
					literal = append(literal, `<td class="invalid"><a href="/cp/`, template.HTML(fmt.Sprintf("%U", tables[i].rows[row].row[j])), `">`, template.HTML(tables[i].rows[row].row[j]), "</a></td>")
				}
			}
			literal = append(literal, "</tr>")
		}
		literal = append(literal, "</table><br>")
	}

	literal = append(literal, "</div>")

	output := htmlJoin(literal, "")

	return output
}

func htmlJoin(elems []template.HTML, sep string) template.HTML {
	switch len(elems) {
	case 0:
		return ""
	case 1:
		return elems[0]
	}
	n := len(sep) * (len(elems) - 1)
	for i := 0; i < len(elems); i++ {
		n += len(elems[i])
	}

	var b strings.Builder
	b.Grow(n)
	b.WriteString(string(elems[0]))
	for _, s := range elems[1:] {
		b.WriteString(sep)
		b.WriteString(string(s))
	}
	return template.HTML(b.String())
}
