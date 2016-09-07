package tlit

import (
	"bufio"
	"bytes"
	"io"
	"regexp"
	"strings"
)

// An Encoder writes transliteration to an output stream.
type Encoder struct {
	*bufio.Writer
	sys System
	tbl map[rune]string
}

// NewEncoder returns a new encoder that writes to w.
func NewEncoder(w io.Writer, sys System) *Encoder {
	return &Encoder{
		Writer: bufio.NewWriter(w),
		sys:    sys,
		tbl:    tableTransliteration[sys],
	}
}

// Encode writes the transliteration encoding of data to the stream.
func (enc *Encoder) Encode(data []byte) error {
	r := bytes.Runes(data)
	l := len(r)
	var rPrev, rNext rune
	for i, rCurr := range r {
		if i+1 <= l {
			rNext = r[i]
		} else {
			rNext = 0
		}
		var err error
		if s, ok := enc.tbl[rCurr]; ok {
			if sFix, ok := fixRuleRune(rPrev, rCurr, rNext, enc.sys); ok {
				s = sFix
			}
			_, err = enc.WriteString(s)
			if err != nil {
				return err
			}
		} else {
			_, err = enc.WriteRune(rCurr)
			if err != nil {
				return err
			}
		}
		rPrev = rCurr
	}
	return enc.Flush()
}

// EncodeString is a convenience wrapper for Encode()
func (enc *Encoder) EncodeString(s string) error {
	return enc.Encode([]byte(s))
}

// Marshal returns the translit encoding of data.
func Marshal(data []byte, sys System) ([]byte, error) {
	var b bytes.Buffer
	if err := NewEncoder(&b, sys).Encode(data); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

// MarshalString is like Marshal but applies string in the input and output.
func MarshalString(s string, sys System) (string, error) {
	b, err := Marshal([]byte(s), sys)
	return string(b), err
}

// MarshalStringURL transforms input string into part of URL
func MarshalStringURL(s string, sys System) string {
	reg := regexp.MustCompile("[^A-Za-z0-9 ]+")
	s, _ = MarshalString(strings.Replace(s, "-", " ", -1), sys)
	return strings.ToLower(strings.Join(strings.Fields(reg.ReplaceAllString(s, "")), "-"))
}

// MarshalStringURLru is syntactic sugar
func MarshalStringURLru(s string) string {
	return MarshalStringURL(s, Default)
}

// MarshalStringURLua is syntactic sugar
func MarshalStringURLua(s string) string {
	return MarshalStringURL(s, UkrainianWeb)
}
