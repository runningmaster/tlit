package tlit

import (
	"bufio"
	"io"
	"regexp"
	"strings"
	"unicode/utf8"
)

// An Encoder writes transliteration to an output stream.
type Encoder struct {
	w   *bufio.Writer
	sys System
	tbl map[rune]string
}

// NewEncoder returns a new encoder that writes to w.
func NewEncoder(w io.Writer, sys System) *Encoder {
	return &Encoder{
		w:   bufio.NewWriter(w),
		sys: sys,
		tbl: tableTransliteration[sys],
	}
}

// Encode writes the transliteration encoding of data to the stream and flushes.
// Close is provided for io.Closer compatibility but is not required after Encode.
func (enc *Encoder) Encode(data []byte) error {
	var prev rune

	for len(data) > 0 {
		v, size := utf8.DecodeRune(data)
		data = data[size:]

		var next rune
		if len(data) > 0 {
			next, _ = utf8.DecodeRune(data)
		}

		var err error

		if s, ok := enc.tbl[v]; ok {
			if sFix, ok := fixRuleRune(prev, v, next, enc.sys); ok {
				s = sFix
			}

			_, err = enc.w.WriteString(s)
		} else {
			_, err = enc.w.WriteRune(v)
		}

		if err != nil {
			return err
		}

		prev = v
	}

	return enc.w.Flush()
}

// EncodeString is a convenience wrapper for Encode.
func (enc *Encoder) EncodeString(s string) error {
	return enc.Encode([]byte(s))
}

// Close flushes any buffered data to the underlying writer.
// Encode already flushes, so Close is only needed when using the Encoder
// as an io.Closer or when no Encode call has been made.
func (enc *Encoder) Close() error {
	return enc.w.Flush()
}

// translitString is the single implementation of the transliteration loop.
// Marshal and MarshalString are thin wrappers around it.
func translitString(s string, tbl map[rune]string, sys System) string {
	var b strings.Builder
	b.Grow(len(s) * 2)

	var prev rune

	for len(s) > 0 {
		v, size := utf8.DecodeRuneInString(s)
		s = s[size:]

		var next rune

		if len(s) > 0 {
			next, _ = utf8.DecodeRuneInString(s)
		}

		if out, ok := tbl[v]; ok {
			if fix, ok := fixRuleRune(prev, v, next, sys); ok {
				out = fix
			}
			b.WriteString(out)
		} else {
			b.WriteRune(v)
		}

		prev = v
	}

	return b.String()
}

// Marshal returns the translit encoding of data.
// It delegates to MarshalString, which incurs two extra copies vs a direct
// bytes.Buffer implementation (string(data) and []byte(result)).
// If Marshal is on a hot path, consider MarshalString or Encoder instead.
func Marshal(data []byte, sys System) []byte {
	return []byte(MarshalString(string(data), sys))
}

// MarshalString returns the translit encoding of s.
func MarshalString(s string, sys System) string {
	return translitString(s, tableTransliteration[sys], sys)
}

var reURL = regexp.MustCompile("[^A-Za-z0-9 ]+")

// MarshalStringURL transforms input string into a URL path segment.
func MarshalStringURL(s string, sys System) string {
	s = MarshalString(strings.ReplaceAll(s, "-", " "), sys)

	return strings.ToLower(strings.Join(strings.Fields(reURL.ReplaceAllString(s, "")), "-"))
}

// MarshalStringURLru is syntactic sugar for Russian (Default system).
func MarshalStringURLru(s string) string {
	return MarshalStringURL(s, Default)
}

// MarshalStringURLua is syntactic sugar for Ukrainian (UkrainianWeb system).
func MarshalStringURLua(s string) string {
	return MarshalStringURL(s, UkrainianWeb)
}
