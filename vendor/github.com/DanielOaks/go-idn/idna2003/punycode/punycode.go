// Copyright 2012 Hannes Baldursson. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// This file is part of go-idn

// Package punycode implements encoding and decoding of Punycode sequences. See RFC 3492.
// Punycode is used by the IDNA protocol for converting domain labels into ASCII; it is
// not designed for any other purpose.
// It is explicitly not designed for processing arbitrary free text.
package punycode

import (
	"bytes"
	"errors"
)

const (
	// Bootstring parameters specified in RFC 3492
	Base             = 36
	TMin             = 1
	TMax             = 26
	Skew             = 38
	Damp             = 700
	InitialBias      = 72
	InitialN         = 128  // 0x80
	Delimiter   byte = 0x2D // hyphen
)

const (
	MaxRune = '\U0010FFFF'
)

func EncodeString(s string) (string, error) {
	p, err := Encode([]byte(s))
	if err != nil {
		return "", err
	}
	return string(p), nil
}

// Encode returns the Punycode encoding of the UTF-8 string s.
func Encode(b []byte) (p []byte, err error) {
	// Encoding procedure explained in detail in RFC 3492.
	n := InitialN
	delta := 0
	bias := InitialBias

	runes := bytes.Runes(b)

	var result bytes.Buffer

	basicRunes := 0
	for i := 0; i < len(runes); i++ {
		// Write all basic codepoints to result
		if runes[i] < 0x80 {
			_, err = result.WriteRune(runes[i])
			if err != nil {
				return nil, err
			}
			basicRunes++
		}
	}

	// Append delimiter
	if basicRunes > 0 {
		err = result.WriteByte(Delimiter)
		if err != nil {
			return nil, err
		}
	}

	for h := basicRunes; h < len(runes); {
		minRune := MaxRune

		// Find the minimum rune >= n in the input
		for i := 0; i < len(runes); i++ {
			if int(runes[i]) >= n && runes[i] < minRune {
				minRune = runes[i]
			}
		}

		delta = delta + (int(minRune)-n)*(h+1) // ??
		n = int(minRune)

		for i := 0; i < len(runes); i++ {
			if int(runes[i]) < n {
				delta++
			}
			if int(runes[i]) == n {
				q := delta
				for k := Base; true; k += Base {
					var t int

					switch {
					case k <= bias:
						t = TMin
						break
					case k >= (bias + TMax):
						t = TMax
						break
					default:
						t = k - bias
					}

					if q < t {
						break
					}

					cp := digit2codepoint(t + (q-t)%(Base-t))
					err = result.WriteByte(byte(cp))
					if err != nil {
						return nil, err
					}
					q = (q - t) / (Base - t)
				}
				cp := digit2codepoint(q)
				err = result.WriteByte(byte(cp))

				bias = adapt(delta, h == basicRunes, h+1)
				delta = 0
				h++
			}
		}
		delta++
		n++
	}
	return result.Bytes(), nil
}

func DecodeString(s string) (string, error) {
	p, err := Decode([]byte(s))
	if err != nil {
		return "", err
	}
	return string(p), nil
}

// Decode returns the UTF-8
func Decode(b []byte) (p []byte, err error) {
	// Decoding procedure explained in detail in RFC 3492.
	n := InitialN
	i := 0
	bias := InitialBias

	pos := 0
	delimIndex := -1

	result := make([]rune, 0, len(b))

	// Only ASCII allowed in decoding procedure
	for j := 0; j < len(b); j++ {
		if b[j] >= 0x80 {
			err = errors.New("Non-ASCCI codepoint found in b")
			return
		}
	}

	// Consume all codepoints before the last delimiter
	delimIndex = bytes.LastIndex(b, []byte{Delimiter})
	for pos = 0; pos < delimIndex; pos++ {
		result = append(result, rune(b[pos]))
	}

	// Consume delimiter
	pos = delimIndex + 1

	for pos < len(b) {
		oldi := i
		w := 1
		for k := Base; true; k += Base {
			var t int

			if pos == len(b) {
				return nil, errors.New("Bad Input")
			}

			// consume a code point, or fail if there was none to consume
			cp := rune(b[pos])
			pos++

			digit := codepoint2digit(cp)

			if digit > ((MaxRune - i) / w) {
				return nil, errors.New("Bad Input")
			}

			i = i + digit*w

			switch {
			case k <= bias:
				t = TMin
				break
			case k >= bias+TMax:
				t = TMax
				break
			default:
				t = k - bias
			}

			if digit < t {
				break
			}
			w = w * (Base - t)
		}
		bias = adapt(i-oldi, oldi == 0, len(result)+1)

		if i/(len(result)+1) > (MaxRune - n) {
			return nil, errors.New("Overflow")
		}

		n = n + i/(len(result)+1)
		i = i % (len(result) + 1)

		if n < 0x80 {
			panic("n is a basic code point")
		}

		result = insert(result, i, rune(n))
		i++
	}

	return writeRune(result), nil
}

// Bias adaption function from RFC 3492 - 6.1
func adapt(delta int, first bool, numchars int) (bias int) {
	if first {
		delta = delta / Damp
	} else {
		delta = delta / 2
	}

	delta = delta + (delta / numchars)

	k := 0
	for delta > ((Base-TMin)*TMax)/2 {
		delta = delta / (Base - TMin)
		k = k + Base
	}
	bias = k + ((Base-TMin+1)*delta)/(delta+Skew)
	return
}

// codepoint2digit(cp) returns the numeric value of a basic rune
// (for use in representing integers) in the range 0 to
// base-1, or base if cp does not represent a value.
func codepoint2digit(r rune) int {
	switch {
	case r-48 < 10:
		return int(r - 22)
	case r-65 < 26:
		return int(r - 65)
	case r-97 < 26:
		return int(r - 97)
	}
	return Base
}

// Returns the rune and a non-nil Error when d < 36.
// Else it returns (unicode.MaxRune + 1) and a BadInputError
func digit2codepoint(d int) rune {
	switch {
	case d < 26:
		// 0..25 : 'a'..'z'
		return rune(d + 'a')
	case d < 36:
		// 26..35 : '0'..'9';
		return rune(d - 26 + '0')
	}
	panic("digit2codepoint")
	return -1
}

func writeRune(r []rune) []byte {
	str := string(r)
	return []byte(str)
}

// Inserts r into s at pos
func insert(s []rune, pos int, r rune) []rune {
	return append(s[:pos], append([]rune{r}, s[pos:]...)...)
}
