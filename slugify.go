package slugify

import (
	"bytes"
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"golang.org/x/text/unicode/norm"
)

var SKIP = []*unicode.RangeTable{
	unicode.Mark,
	unicode.Sk,
	unicode.Lm,
}

var SAFE = []*unicode.RangeTable{
	unicode.Letter,
	unicode.Number,
}

var DASH = []*unicode.RangeTable{
	unicode.Pd,
}

var SPACE = []*unicode.RangeTable{
	unicode.Space,
}

var OK = "-_~."
var TO_DASH = "/\\—–"
var ID_OK = "-_"
var ID_TO_DASH = "/\\—–.~"
var extra_dashes = regexp.MustCompile("[-]{2,}")

// A very limited list of transliterations to catch common european names translated to urls.
// This set could be expanded with at least caps and many more characters.
var transliterations = map[rune]string{
	'À': "A",
	'Á': "A",
	'Â': "A",
	'Ã': "A",
	'Ä': "A",
	'Å': "AA",
	'Æ': "AE",
	'Ç': "C",
	'È': "E",
	'É': "E",
	'Ê': "E",
	'Ë': "E",
	'Ì': "I",
	'Í': "I",
	'Î': "I",
	'Ï': "I",
	'Ð': "D",
	'Ł': "L",
	'Ñ': "N",
	'Ò': "O",
	'Ó': "O",
	'Ô': "O",
	'Õ': "O",
	'Ö': "O",
	'Ø': "OE",
	'Ù': "U",
	'Ú': "U",
	'Ü': "U",
	'Û': "U",
	'Ý': "Y",
	'Þ': "Th",
	'ß': "ss",
	'à': "a",
	'á': "a",
	'â': "a",
	'ã': "a",
	'ä': "a",
	'å': "aa",
	'æ': "ae",
	'ç': "c",
	'è': "e",
	'é': "e",
	'ê': "e",
	'ë': "e",
	'ì': "i",
	'í': "i",
	'î': "i",
	'ï': "i",
	'ð': "d",
	'ł': "l",
	'ñ': "n",
	'ń': "n",
	'ò': "o",
	'ó': "o",
	'ô': "o",
	'õ': "o",
	'ō': "o",
	'ö': "o",
	'ø': "oe",
	'ś': "s",
	'ù': "u",
	'ú': "u",
	'û': "u",
	'ū': "u",
	'ü': "u",
	'ý': "y",
	'þ': "th",
	'ÿ': "y",
	'ż': "z",
	'Œ': "OE",
	'œ': "oe",
}

// Slugify a string. The result will only contain lowercase letters,
// digits and dashes. It will not begin or end with a dash, and it
// will not contain runs of multiple dashes.
//
// It is NOT forced into being ASCII, but may contain any Unicode
// characters, with the above restrictions.
func Slugify(text string) string {
	buf := make([]rune, 0, len(text))
	text = SanatizeText(text)
	for _, r := range norm.NFKD.String(text) {
		s := strconv.QuoteRune(r)
		switch {
		case unicode.IsOneOf(SAFE, r):
			buf = append(buf, unicode.ToLower(r))
		case strings.ContainsAny(s, OK):
			buf = append(buf, r)
		case unicode.IsOneOf(SPACE, r):
			buf = append(buf, '-')
		case unicode.IsOneOf(DASH, r):
			buf = append(buf, '-')
		case strings.ContainsAny(s, TO_DASH):
			buf = append(buf, '-')
		}
	}
	return cleanup(string(buf))
}

func IDify(text string) string {
	buf := make([]rune, 0, len(text))
	text = SanatizeText(text)
	for _, r := range norm.NFKD.String(text) {
		s := strconv.QuoteRune(r)
		switch {
		case unicode.IsOneOf(SAFE, r):
			buf = append(buf, unicode.ToLower(r))
		case strings.ContainsAny(s, ID_OK):
			buf = append(buf, r)
		case unicode.IsOneOf(SPACE, r):
			buf = append(buf, '-')
		case unicode.IsOneOf(DASH, r):
			buf = append(buf, '-')
		case strings.ContainsAny(s, ID_TO_DASH):
			buf = append(buf, '-')
		}
	}
	return cleanup(string(buf))
}

func SanatizeText(text string) string {
	text = strings.ToLower(text)
	b := bytes.NewBufferString("")
	for _, c := range text {
		// Check transliterations first
		if val, ok := transliterations[c]; ok {
			b.WriteString(val)
		} else {
			b.WriteRune(c)
		}
	}
	return b.String()
}

func cleanup(text string) string {
	text = strings.Trim(text, "-")
	text = extra_dashes.ReplaceAllString(text, "-")
	return text
}
