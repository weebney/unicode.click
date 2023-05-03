package main

import "unicode"

func getCategoryData(codepoint rune) (majorCategoryLiteral string, categoryLiteral string, categories []string, majorCategories []string) {
	for categoryName, categoryRangeTable := range unicode.Categories {
		if unicode.Is(categoryRangeTable, codepoint) {
			if len(categoryName) == 1 {
				switch categoryName {
				case "C":
					categoryName = "Other (C)"
					majorCategoryLiteral = (majorCategoryLiteral + "c")
				case "L":
					categoryName = "Letter (L)"
					majorCategoryLiteral = (majorCategoryLiteral + "l")
				case "M":
					categoryName = "Mark (M)"
					majorCategoryLiteral = (majorCategoryLiteral + "m")
				case "N":
					categoryName = "Number (N)"
					majorCategoryLiteral = (majorCategoryLiteral + "n")
				case "P":
					categoryName = "Punctuation (P)"
					majorCategoryLiteral = (majorCategoryLiteral + "p")
				case "S":
					categoryName = "Symbol (S)"
					majorCategoryLiteral = (majorCategoryLiteral + "s")
				case "Z":
					categoryName = "Separator (Z)"
					majorCategoryLiteral = (majorCategoryLiteral + "z")
				}

				majorCategories = append(majorCategories, categoryName)
			} else {
				switch categoryName {

				// THE Cs
				case "Cc":
					categoryName = "Control (Cc)"
					categoryLiteral = (categoryLiteral + "cc")
				case "Cf":
					categoryName = "Format (Cf)"
					categoryLiteral = (categoryLiteral + "cf")
				case "Co":
					categoryName = "Private use (Co)"
					categoryLiteral = (categoryLiteral + "co")
				case "Cs":
					categoryName = "Surrogate (Cs)"
					categoryLiteral = (categoryLiteral + "cs")
				// THE Ns
				case "Nd":
					categoryName = "Decimal (Nd)"
					categoryLiteral = (categoryLiteral + "nd")
				case "Nl":
					categoryName = "Letter (Nl)"
					categoryLiteral = (categoryLiteral + "nl")
				case "No":
					categoryName = "Other (No)"
					categoryLiteral = (categoryLiteral + "no")

				// THE Ls
				case "Ll":
					categoryName = "Lowercase (Ll)"
					categoryLiteral = (categoryLiteral + "ll")
				case "Lm":
					categoryName = "Modifier (Lm)"
					categoryLiteral = (categoryLiteral + "lm")
				case "Lo":
					categoryName = "Other (Lo)"
					categoryLiteral = (categoryLiteral + "lo")
				case "Lt":
					categoryName = "Titlecase (Lt)"
					categoryLiteral = (categoryLiteral + "lt")
				case "Lu":
					categoryName = "Uppercase (Lu)"
					categoryLiteral = (categoryLiteral + "lu")

				// THE Ms
				case "Mc":
					categoryName = "Spacing (Mc)"
					categoryLiteral = (categoryLiteral + "mc")
				case "Me":
					categoryName = "Enclosing (Me)"
					categoryLiteral = (categoryLiteral + "me")
				case "Mn":
					categoryName = "Nonspacing (Mn)"
					categoryLiteral = (categoryLiteral + "mn")

				// THE Ps
				case "Pc":
					categoryName = "Connector (Pc)"
					categoryLiteral = (categoryLiteral + "pc")
				case "Pd":
					categoryName = "Dash (Pd)"
					categoryLiteral = (categoryLiteral + "pd")
				case "Pe":
					categoryName = "Close (Pe)"
					categoryLiteral = (categoryLiteral + "pe")
				case "Pf":
					categoryName = "Final quote (Pf)"
					categoryLiteral = (categoryLiteral + "pf")
				case "Pi":
					categoryName = "Initial quote (Pi)"
					categoryLiteral = (categoryLiteral + "pi")
				case "Po":
					categoryName = "Other (Po)"
					categoryLiteral = (categoryLiteral + "po")
				case "Ps":
					categoryName = "Open (Ps)"
					categoryLiteral = (categoryLiteral + "ps")

				// THE Ss
				case "Sc":
					categoryName = "Currency (Sc)"
					categoryLiteral = (categoryLiteral + "sc")
				case "Sk":
					categoryName = "Modifier (Sk)"
					categoryLiteral = (categoryLiteral + "sk")
				case "Sm":
					categoryName = "Math (Sm)"
					categoryLiteral = (categoryLiteral + "sm")
				case "So":
					categoryName = "Other (So)"
					categoryLiteral = (categoryLiteral + "so")

				// THE Zs
				case "Zl":
					categoryName = "Line (Zl)"
					categoryLiteral = (categoryLiteral + "zl")
				case "Zp":
					categoryName = "Paragraph (Zp)"
					categoryLiteral = (categoryLiteral + "zp")
				case "Zs":
					categoryName = "Space (Zs)"
					categoryLiteral = (categoryLiteral + "zs")
				}
				categories = append(categories, categoryName)
			}
		}
	}
	return
}

func getRangeTableLiteral(route string) (rtLiteral *unicode.RangeTable) {
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

	return
}
