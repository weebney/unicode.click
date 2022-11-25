package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"text/template"
	"time"
	"unicode"

	"github.com/julienschmidt/httprouter"
	"golang.org/x/text/unicode/runenames"
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

func codepointFromRoute(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {

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

	var mcatSlice string
	var catSlice string
	var categories []string
	var majorCategories []string
	for categoryName, categoryRangeTable := range unicode.Categories {
		if unicode.Is(categoryRangeTable, codepoint) {
			if len(categoryName) == 1 {
				switch categoryName {
				case "C":
					categoryName = "Other (C)"
					mcatSlice = (mcatSlice + "c")
				case "L":
					categoryName = "Letter (L)"
					mcatSlice = (mcatSlice + "l")
				case "M":
					categoryName = "Mark (M)"
					mcatSlice = (mcatSlice + "m")
				case "N":
					categoryName = "Number (N)"
					mcatSlice = (mcatSlice + "n")
				case "P":
					categoryName = "Punctuation (P)"
					mcatSlice = (mcatSlice + "p")
				case "S":
					categoryName = "Symbol (S)"
					mcatSlice = (mcatSlice + "s")
				case "Z":
					categoryName = "Separator (Z)"
					mcatSlice = (mcatSlice + "z")
				}

				majorCategories = append(majorCategories, categoryName)
			} else {
				switch categoryName {

				// THE Cs
				case "Cc":
					categoryName = "Control (Cc)"
					catSlice = (catSlice + "cc")
				case "Cf":
					categoryName = "Format (Cf)"
					catSlice = (catSlice + "cf")
				case "Co":
					categoryName = "Private use (Co)"
					catSlice = (catSlice + "co")
				case "Cs":
					categoryName = "Surrogate (Cs)"
					catSlice = (catSlice + "cs")
				// THE Ns
				case "Nd":
					categoryName = "Decimal (Nd)"
					catSlice = (catSlice + "nd")
				case "Nl":
					categoryName = "Letter (Nl)"
					catSlice = (catSlice + "nl")
				case "No":
					categoryName = "Other (No)"
					catSlice = (catSlice + "no")

				// THE Ls
				case "Ll":
					categoryName = "Lowercase (Ll)"
					catSlice = (catSlice + "ll")
				case "Lm":
					categoryName = "Modifier (Lm)"
					catSlice = (catSlice + "lm")
				case "Lo":
					categoryName = "Other (Lo)"
					catSlice = (catSlice + "lo")
				case "Lt":
					categoryName = "Titlecase (Lt)"
					catSlice = (catSlice + "lt")
				case "Lu":
					categoryName = "Uppercase (Lu)"
					catSlice = (catSlice + "lu")

				// THE Ms
				case "Mc":
					categoryName = "Spacing (Mc)"
					catSlice = (catSlice + "mc")
				case "Me":
					categoryName = "Enclosing (Me)"
					catSlice = (catSlice + "me")
				case "Mn":
					categoryName = "Nonspacing (Mn)"
					catSlice = (catSlice + "mn")

				// THE Ps
				case "Pc":
					categoryName = "Connector (Pc)"
					catSlice = (catSlice + "pc")
				case "Pd":
					categoryName = "Dash (Pd)"
					catSlice = (catSlice + "pd")
				case "Pe":
					categoryName = "Close (Pe)"
					catSlice = (catSlice + "pe")
				case "Pf":
					categoryName = "Final quote (Pf)"
					catSlice = (catSlice + "pf")
				case "Pi":
					categoryName = "Initial quote (Pi)"
					catSlice = (catSlice + "pi")
				case "Po":
					categoryName = "Other (Po)"
					catSlice = (catSlice + "po")
				case "Ps":
					categoryName = "Open (Ps)"
					catSlice = (catSlice + "ps")

				// THE Ss
				case "Sc":
					categoryName = "Currency (Sc)"
					catSlice = (catSlice + "sc")
				case "Sk":
					categoryName = "Modifier (Sk)"
					catSlice = (catSlice + "sk")
				case "Sm":
					categoryName = "Math (Sm)"
					catSlice = (catSlice + "sm")
				case "So":
					categoryName = "Other (So)"
					catSlice = (catSlice + "so")

				// THE Zs
				case "Zl":
					categoryName = "Line (Zl)"
					catSlice = (catSlice + "zl")
				case "Zp":
					categoryName = "Paragraph (Zp)"
					catSlice = (catSlice + "zp")
				case "Zs":
					categoryName = "Space (Zs)"
					catSlice = (catSlice + "zs")
				}
				categories = append(categories, categoryName)
			}
		}
	}

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

		MajCatLiteral:   mcatSlice,
		CatLiteral:      catSlice,
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

	serveFilesFromTemplate(writer, request, params, templateFiles, data)

}

func rangeFromRoute(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {

	route := strings.ToLower(params.ByName("rangeRoute"))

	writer.Header().Set("Content-Type", "text/html; charset=utf-8")
	writer.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	writer.Header().Set("Pragma", "no-cache")
	writer.Header().Set("Expires", "0")

	// ------------------------------------------------------
	// EVERYTHING BELOW THIS IS FOR THE RANGETABLE PAGE

	var rtLiteral *unicode.RangeTable

	switch route {
	case "adlam":
		rtLiteral = unicode.Adlam
	case "ahom":
		rtLiteral = unicode.Ahom
	case "anatolian_hieroglyphs":
		rtLiteral = unicode.Anatolian_Hieroglyphs
	case "arabic":
		rtLiteral = unicode.Arabic
	case "armenian":
		rtLiteral = unicode.Armenian
	case "avestan":
		rtLiteral = unicode.Avestan
	case "balinese":
		rtLiteral = unicode.Balinese
	case "bamum":
		rtLiteral = unicode.Bamum
	case "bassa_vah":
		rtLiteral = unicode.Bassa_Vah
	case "batak":
		rtLiteral = unicode.Batak
	case "bengali":
		rtLiteral = unicode.Bengali
	case "bhaiksuki":
		rtLiteral = unicode.Bhaiksuki
	case "bopomofo":
		rtLiteral = unicode.Bopomofo
	case "brahmi":
		rtLiteral = unicode.Brahmi
	case "braille":
		rtLiteral = unicode.Braille
	case "buginese":
		rtLiteral = unicode.Buginese
	case "buhid":
		rtLiteral = unicode.Buhid
	case "canadian_aboriginal":
		rtLiteral = unicode.Canadian_Aboriginal
	case "carian":
		rtLiteral = unicode.Carian
	case "caucasian_albanian":
		rtLiteral = unicode.Caucasian_Albanian
	case "chakma":
		rtLiteral = unicode.Chakma
	case "cham":
		rtLiteral = unicode.Cham
	case "cherokee":
		rtLiteral = unicode.Cherokee
	case "chorasmian":
		rtLiteral = unicode.Chorasmian
	case "common":
		rtLiteral = unicode.Common
	case "coptic":
		rtLiteral = unicode.Coptic
	case "cuneiform":
		rtLiteral = unicode.Cuneiform
	case "cypriot":
		rtLiteral = unicode.Cypriot
	case "cyrillic":
		rtLiteral = unicode.Cyrillic
	case "deseret":
		rtLiteral = unicode.Deseret
	case "devanagari":
		rtLiteral = unicode.Devanagari
	case "dives_akuru":
		rtLiteral = unicode.Dives_Akuru
	case "dogra":
		rtLiteral = unicode.Dogra
	case "duployan":
		rtLiteral = unicode.Duployan
	case "egyptian_hieroglyphs":
		rtLiteral = unicode.Egyptian_Hieroglyphs
	case "elbasan":
		rtLiteral = unicode.Elbasan
	case "elymaic":
		rtLiteral = unicode.Elymaic
	case "ethiopic":
		rtLiteral = unicode.Ethiopic
	case "georgian":
		rtLiteral = unicode.Georgian
	case "glagolitic":
		rtLiteral = unicode.Glagolitic
	case "gothic":
		rtLiteral = unicode.Gothic
	case "grantha":
		rtLiteral = unicode.Grantha
	case "greek":
		rtLiteral = unicode.Greek
	case "gujarati":
		rtLiteral = unicode.Gujarati
	case "gunjala_gondi":
		rtLiteral = unicode.Gunjala_Gondi
	case "gurmukhi":
		rtLiteral = unicode.Gurmukhi
	case "han":
		rtLiteral = unicode.Han
	case "hangul":
		rtLiteral = unicode.Hangul
	case "hanifi_rohingya":
		rtLiteral = unicode.Hanifi_Rohingya
	case "hanunoo":
		rtLiteral = unicode.Hanunoo
	case "hatran":
		rtLiteral = unicode.Hatran
	case "hebrew":
		rtLiteral = unicode.Hebrew
	case "hiragana":
		rtLiteral = unicode.Hiragana
	case "imperial_aramaic":
		rtLiteral = unicode.Imperial_Aramaic
	case "inherited":
		rtLiteral = unicode.Inherited
	case "inscriptional_pahlavi":
		rtLiteral = unicode.Inscriptional_Pahlavi
	case "inscriptional_parthian":
		rtLiteral = unicode.Inscriptional_Parthian
	case "javanese":
		rtLiteral = unicode.Javanese
	case "kaithi":
		rtLiteral = unicode.Kaithi
	case "kannada":
		rtLiteral = unicode.Kannada
	case "katakana":
		rtLiteral = unicode.Katakana
	case "kayah_li":
		rtLiteral = unicode.Kayah_Li
	case "kharoshthi":
		rtLiteral = unicode.Kharoshthi
	case "khitan_small_script":
		rtLiteral = unicode.Khitan_Small_Script
	case "khmer":
		rtLiteral = unicode.Khmer
	case "khojki":
		rtLiteral = unicode.Khojki
	case "khudawadi":
		rtLiteral = unicode.Khudawadi
	case "lao":
		rtLiteral = unicode.Lao
	case "latin":
		rtLiteral = unicode.Latin
	case "lepcha":
		rtLiteral = unicode.Lepcha
	case "limbu":
		rtLiteral = unicode.Limbu
	case "linear_a":
		rtLiteral = unicode.Linear_A
	case "linear_b":
		rtLiteral = unicode.Linear_B
	case "lisu":
		rtLiteral = unicode.Lisu
	case "lycian":
		rtLiteral = unicode.Lycian
	case "lydian":
		rtLiteral = unicode.Lydian
	case "mahajani":
		rtLiteral = unicode.Mahajani
	case "makasar":
		rtLiteral = unicode.Makasar
	case "malayalam":
		rtLiteral = unicode.Malayalam
	case "mandaic":
		rtLiteral = unicode.Mandaic
	case "manichaean":
		rtLiteral = unicode.Manichaean
	case "marchen":
		rtLiteral = unicode.Marchen
	case "masaram_gondi":
		rtLiteral = unicode.Masaram_Gondi
	case "medefaidrin":
		rtLiteral = unicode.Medefaidrin
	case "meetei_mayek":
		rtLiteral = unicode.Meetei_Mayek
	case "mende_kikakui":
		rtLiteral = unicode.Mende_Kikakui
	case "meroitic_cursive":
		rtLiteral = unicode.Meroitic_Cursive
	case "meroitic_hieroglyphs":
		rtLiteral = unicode.Meroitic_Hieroglyphs
	case "miao":
		rtLiteral = unicode.Miao
	case "modi":
		rtLiteral = unicode.Modi
	case "mongolian":
		rtLiteral = unicode.Mongolian
	case "mro":
		rtLiteral = unicode.Mro
	case "multani":
		rtLiteral = unicode.Multani
	case "myanmar":
		rtLiteral = unicode.Myanmar
	case "nabataean":
		rtLiteral = unicode.Nabataean
	case "nandinagar":
		rtLiteral = unicode.Nandinagari
	case "new_tai_lue":
		rtLiteral = unicode.New_Tai_Lue
	case "newa":
		rtLiteral = unicode.Newa
	case "nko":
		rtLiteral = unicode.Nko
	case "nushu":
		rtLiteral = unicode.Nushu
	case "nyiakeng_puachue_hmong":
		rtLiteral = unicode.Nyiakeng_Puachue_Hmong
	case "ogham":
		rtLiteral = unicode.Ogham
	case "ol_chiki":
		rtLiteral = unicode.Ol_Chiki
	case "old_hungarian":
		rtLiteral = unicode.Old_Hungarian
	case "old_italic":
		rtLiteral = unicode.Old_Italic
	case "old_north_arabian":
		rtLiteral = unicode.Old_North_Arabian
	case "old_permic":
		rtLiteral = unicode.Old_Permic
	case "old_persian":
		rtLiteral = unicode.Old_Persian
	case "old_sogdian":
		rtLiteral = unicode.Old_Sogdian
	case "old_south_arabian":
		rtLiteral = unicode.Old_South_Arabian
	case "old_turkic":
		rtLiteral = unicode.Old_Turkic
	case "oriya":
		rtLiteral = unicode.Oriya
	case "osage":
		rtLiteral = unicode.Osage
	case "osmanya":
		rtLiteral = unicode.Osmanya
	case "pahawh_hmong":
		rtLiteral = unicode.Pahawh_Hmong
	case "palmyrene":
		rtLiteral = unicode.Palmyrene
	case "pau_cin_hau":
		rtLiteral = unicode.Pau_Cin_Hau
	case "phags_pa":
		rtLiteral = unicode.Phags_Pa
	case "phoenician":
		rtLiteral = unicode.Phoenician
	case "psalter_pahlavi":
		rtLiteral = unicode.Psalter_Pahlavi
	case "rejang":
		rtLiteral = unicode.Rejang
	case "runic":
		rtLiteral = unicode.Runic
	case "samaritan":
		rtLiteral = unicode.Samaritan
	case "saurashtra":
		rtLiteral = unicode.Saurashtra
	case "sharada":
		rtLiteral = unicode.Sharada
	case "shavian":
		rtLiteral = unicode.Shavian
	case "siddham":
		rtLiteral = unicode.Siddham
	case "signwriting":
		rtLiteral = unicode.SignWriting
	case "sinhala":
		rtLiteral = unicode.Sinhala
	case "sogdian":
		rtLiteral = unicode.Sogdian
	case "sora_sompeng":
		rtLiteral = unicode.Sora_Sompeng
	case "soyombo":
		rtLiteral = unicode.Soyombo
	case "sundanese":
		rtLiteral = unicode.Sundanese
	case "syloti_nagri":
		rtLiteral = unicode.Syloti_Nagri
	case "syriac":
		rtLiteral = unicode.Syriac
	case "tagalog":
		rtLiteral = unicode.Tagalog
	case "tagbanwa":
		rtLiteral = unicode.Tagbanwa
	case "tai_le":
		rtLiteral = unicode.Tai_Le
	case "tai_tham":
		rtLiteral = unicode.Tai_Tham
	case "tai_viet":
		rtLiteral = unicode.Tai_Viet
	case "takri":
		rtLiteral = unicode.Takri
	case "tamil":
		rtLiteral = unicode.Tamil
	case "tangut":
		rtLiteral = unicode.Tangut
	case "telugu":
		rtLiteral = unicode.Telugu
	case "thaana":
		rtLiteral = unicode.Thaana
	case "thai":
		rtLiteral = unicode.Thai
	case "tibetan":
		rtLiteral = unicode.Tibetan
	case "tifinagh":
		rtLiteral = unicode.Tifinagh
	case "tirhuta":
		rtLiteral = unicode.Tirhuta
	case "ugaritic":
		rtLiteral = unicode.Ugaritic
	case "vai":
		rtLiteral = unicode.Vai
	case "wancho":
		rtLiteral = unicode.Wancho
	case "warang_citi":
		rtLiteral = unicode.Warang_Citi
	case "yezidi":
		rtLiteral = unicode.Yezidi
	case "yi":
		rtLiteral = unicode.Yi
	case "zanabazar_square":
		rtLiteral = unicode.Zanabazar_Square

		// properties

	case "ascii_hex_digit":
		rtLiteral = unicode.ASCII_Hex_Digit
	case "bidi_control":
		rtLiteral = unicode.Bidi_Control
	case "dash":
		rtLiteral = unicode.Dash
	case "deprecated":
		rtLiteral = unicode.Deprecated
	case "diacritic":
		rtLiteral = unicode.Diacritic
	case "extender":
		rtLiteral = unicode.Extender
	case "hex_digit":
		rtLiteral = unicode.Hex_Digit
	case "hyphen":
		rtLiteral = unicode.Hyphen
	case "ids_binary_operator":
		rtLiteral = unicode.IDS_Binary_Operator
	case "ids_trinary_operator":
		rtLiteral = unicode.IDS_Trinary_Operator
	case "ideographic":
		rtLiteral = unicode.Ideographic
	case "join_control":
		rtLiteral = unicode.Join_Control
	case "logical_order_exception":
		rtLiteral = unicode.Logical_Order_Exception
		/* 	this causes the server to fail epicly
		case "noncharacter_code_point":
			rtLiteral = unicode.Noncharacter_Code_Point */
	case "other_alphabetic":
		rtLiteral = unicode.Other_Alphabetic
	case "other_default_ignorable_code_point":
		rtLiteral = unicode.Other_Default_Ignorable_Code_Point
	case "other_grapheme_extend":
		rtLiteral = unicode.Other_Grapheme_Extend
	case "other_id_continue":
		rtLiteral = unicode.Other_ID_Continue
	case "other_id_start":
		rtLiteral = unicode.Other_ID_Start
	case "other_lowercase":
		rtLiteral = unicode.Other_Lowercase
	case "other_math":
		rtLiteral = unicode.Other_Math
	case "other_uppercase":
		rtLiteral = unicode.Other_Uppercase
	case "pattern_syntax":
		rtLiteral = unicode.Pattern_Syntax
	case "pattern_white_space":
		rtLiteral = unicode.Pattern_White_Space
	case "prepended_concatenation_mark":
		rtLiteral = unicode.Prepended_Concatenation_Mark
	case "quotation_mark":
		rtLiteral = unicode.Quotation_Mark
	case "radical":
		rtLiteral = unicode.Radical
	case "regional_indicator":
		rtLiteral = unicode.Regional_Indicator
	case "sterm":
		rtLiteral = unicode.Sentence_Terminal
	case "sentence_terminal":
		rtLiteral = unicode.Sentence_Terminal
	case "soft_dotted":
		rtLiteral = unicode.Soft_Dotted
	case "terminal_punctuation":
		rtLiteral = unicode.Terminal_Punctuation
	case "unified_ideograph":
		rtLiteral = unicode.Unified_Ideograph
	case "variation_selector":
		rtLiteral = unicode.Variation_Selector
	case "white_space":
		rtLiteral = unicode.White_Space

		// categories

	case "cc":
		rtLiteral = unicode.Cc
	case "cf":
		rtLiteral = unicode.Cf
	case "co":
		rtLiteral = unicode.Co
	case "cs":
		rtLiteral = unicode.Cs
	case "nd":
		rtLiteral = unicode.Nd
	case "letter":
		rtLiteral = unicode.L
	case "l":
		rtLiteral = unicode.L
	case "lm":
		rtLiteral = unicode.Lm
	case "lo":
		rtLiteral = unicode.Lo
	case "lower":
		rtLiteral = unicode.Ll
	case "ll":
		rtLiteral = unicode.Ll
	case "mark":
		rtLiteral = unicode.M
	case "m":
		rtLiteral = unicode.M
	case "mc":
		rtLiteral = unicode.Mc
	case "me":
		rtLiteral = unicode.Me
	case "mn":
		rtLiteral = unicode.Mn
	case "nl":
		rtLiteral = unicode.Nl
	case "no":
		rtLiteral = unicode.No
	case "number":
		rtLiteral = unicode.N
	case "n":
		rtLiteral = unicode.N
	case "other":
		rtLiteral = unicode.C
	case "c":
		rtLiteral = unicode.C
	case "pc":
		rtLiteral = unicode.Pc
	case "pd":
		rtLiteral = unicode.Pd
	case "pe":
		rtLiteral = unicode.Pe
	case "pf":
		rtLiteral = unicode.Pf
	case "pi":
		rtLiteral = unicode.Pi
	case "po":
		rtLiteral = unicode.Po
	case "ps":
		rtLiteral = unicode.Ps
	case "punct":
		rtLiteral = unicode.P
	case "p":
		rtLiteral = unicode.P
	case "sc":
		rtLiteral = unicode.Sc
	case "sk":
		rtLiteral = unicode.Sk
	case "sm":
		rtLiteral = unicode.Sm
	case "so":
		rtLiteral = unicode.So
	case "space":
		rtLiteral = unicode.Z
	case "z":
		rtLiteral = unicode.Z
	case "symbol":
		rtLiteral = unicode.S
	case "s":
		rtLiteral = unicode.S
	case "title":
		rtLiteral = unicode.Lt
	case "lt":
		rtLiteral = unicode.Lt
	case "upper":
		rtLiteral = unicode.Lu
	case "lu":
		rtLiteral = unicode.Lu
	case "zl":
		rtLiteral = unicode.Zl
	case "zp":
		rtLiteral = unicode.Zp
	case "zs":
		rtLiteral = unicode.Zs

		// custom

	case "control":
		rtLiteral = unicode.Cc
	case "digit":
		rtLiteral = unicode.Nd

	default:
		rtLiteral = unicode.Latin
		route = "latin"
	}

	// Table Generation Code ================================
	// ++++++++++++++++++++++++++++++++++++++++++++++++++++++
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

	tables := []table{}
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

	// End Table Code =======================================
	// ++++++++++++++++++++++++++++++++++++++++++++++++++++++

	var tableLengths []int
	for i := 0; i < len(tables); i++ {
		tableLengths = append(tableLengths, len(tables[i].rows))
	}

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

	serveFilesFromTemplate(writer, request, params, templateFiles, data)

}

func generateTableHtml(tables []table, tableLengths []int, literalRT *unicode.RangeTable) string {
	var literal []string

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
			literal = append(literal, "<td>U+", tables[i].rows[row].name, "</td>")
			for j := 0; j < len(tables[i].rows[row].row); j++ {
				if unicode.In(tables[i].rows[row].row[j], literalRT) {
					literal = append(literal, `<td><a href="/cp/`, fmt.Sprintf("%U", tables[i].rows[row].row[j]), `">`, string(tables[i].rows[row].row[j]), "</a></td>")
				} else {
					literal = append(literal, `<td class="invalid"><a href="/cp/`, fmt.Sprintf("%U", tables[i].rows[row].row[j]), `">`, string(tables[i].rows[row].row[j]), "</a></td>")
				}
			}
			literal = append(literal, "</tr>")
		}
		literal = append(literal, "</table><br>")
	}

	literal = append(literal, "</div>")

	output := strings.Join(literal, "")

	return output
}

func serveFilesFromTemplate(writer http.ResponseWriter, request *http.Request, params httprouter.Params, templates []string, data interface{}) {

	// Define the template functions
	funcMap := template.FuncMap{
		"until":   templateUntil,
		"iterate": templateIterate,
	}

	// Parse the templates and add the function map
	tmpl, err := template.New("").Funcs(funcMap).ParseFiles(templates...)
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

	log.Printf("%s %s %s %s %s \n", request.UserAgent(), request.RemoteAddr, request.Method, request.URL, request.Proto)
}

func templateUntil(n int) []int {
	result := make([]int, n)
	for i := 0; i < n; i++ {
		result[i] = i
	}
	return result
}

func templateIterate(count int) []int {
	var i int
	var Items []int
	for i = 0; i < (count); i++ {
		Items = append(Items, i)
	}
	return Items
}

func serveFavicon(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	http.ServeFile(writer, request, "./favicon.ico")
}

func serveIndex(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {

	templateFiles := []string{
		"./template/base.template.html",
		"./template/index.template.html",
	}

	serveFilesFromTemplate(writer, request, params, templateFiles, struct{ UnicodeVersion string }{UnicodeVersion: unicode.Version})
}

func redirectToSSL(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	route := params.ByName("route")
	if route == "/" {
		route = ""
	}
	log.Printf("REDIRECTING TO SSL: %s %s %s \n", request.RemoteAddr, request.URL, request.Proto)
	http.Redirect(writer, request, "https://localhost/"+route, http.StatusMovedPermanently)
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

	sslRedirect := httprouter.New()
	sslRedirect.GET("/*route", redirectToSSL)

	fmt.Println("unicode.click listening on 443 and 8080")

	// log.Fatal(http.ListenAndServe(":8080", router))
	log.Fatal(http.ListenAndServeTLS(":443", "./ssl/domain.cert.pem", "./ssl/private.key.pem", router))
}
