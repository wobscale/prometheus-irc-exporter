// Copyright 2012 Hannes Baldursson. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file is part of go-idn

// Package idna2003 implements IDNA as described in RFC 3490.
//
// This package is in beta and has not been extensively tested.
package idna2003

import (
	"errors"
	"strings"

	"github.com/DanielOaks/go-idn/idna2003/punycode"
	"github.com/DanielOaks/go-idn/idna2003/stringprep"
)

// IDNA section 5
const (
	AcePrefix = "xn--"
)

// Converts a Unicode string to ASCII using the procedure in RFC 3490
// section 4.1. Unassigned characters are not allowed and STD3 ASCII rules are
// enforced. The input string may be a domain name containing dots.
func ToASCII(label string) (string, error) {

	label = strings.ToLower(label)
	o := ""
	h := ""

	for _, cp := range label {

		if cp == 0x2E /* dot */ || cp == 0x3002 || cp == 0xff0e || cp == 0xff61 {
			uh, err := toASCIIRaw(h)
			if err != nil {
				return label, err
			}
			o += uh
			o += string(cp)
			h = ""
		} else {
			h += string(cp)
		}
	}
	uh, err := toASCIIRaw(h)
	if err != nil {
		return label, err
	}
	o += uh
	return o, nil
}

func toASCIIRaw(label string) (string, error) {
	original := label

	// Step 1: If the sequence contains any code points outside the ASCII range
	// (0..7F) then proceed to step 2, otherwise skip to step 3.
	for i := 0; i < len(label); i++ {
		if label[i] > 127 {
			// Step 2: Perform the smake teps specified in [NAMEPREP] and fail if there is an error.
			// The AllowUnassigned flag is used in [NAMEPREP].
			// BUG(dan) we should check error value here
			label, _ = stringprep.Nameprep(label)
			break
		}
	}

	// Step 3: - Verify the absence of non-LDH ASCII code points
	for _, c := range label {
		if (c <= 0x2c) || (c >= 0x2e && c <= 0x2f) || (c >= 0x3a && c <= 0x40) || (c >= 0x5b && c <= 0x60) || (c >= 0x7b && c <= 0x7f) {
			return original, errors.New("Contains non-LDH ASCII codepoints")
		}

	}
	if strings.HasPrefix(label, "-") || strings.HasSuffix(label, "-") {
		return original, errors.New("Contains hyphen at either end of the string")
	}

	// Step 4: If the sequence contains any code points outside the ASCII range
	// (0..7F) then proceed to step 5, otherwise skip to step 8.

	isASCII := true
	for i := 0; i < len(label); i++ {
		if label[i] > 127 {
			isASCII = false
			break
		}

	}

	if !isASCII {

		// Step 5 Verify that the sequence does NOT begin with the ACE prefix.
		if strings.HasPrefix(label, AcePrefix) {
			return label, errors.New("Label starts with ACE prefix")
		}

		var err error

		// Step 6: Encode with punycode
		label, err = punycode.EncodeString(label)
		if err != nil {
			return "", err // delegate err
		}
		// Step 7: Prepend ACE prefix
		label = AcePrefix + label
	}

	// 8. Verify that the number of code points is in the range 1 to 63 inclusive.
	if 0 < len(label) && len(label) < 64 {
		return label, nil
	}

	return original, errors.New("label empty or too long")
}

//
// Converts a Punycode string to Unicode using the procedure in RFC 3490
// section 4.2. Unassigned characters are not allowed and STD3 ASCII
// rules are enforced. The input string may be a domain name
// containing dots.
//
// ToUnicode never fails.  If any step fails, then the original input
// sequence is returned immediately in that step.
func ToUnicode(label string) (string, error) {

	label = strings.ToLower(label)
	o := ""
	h := ""

	for _, cp := range label {

		if cp == 0x2E /* dot */ || cp == 0x3002 || cp == 0xff0e || cp == 0xff61 {
			uh, err := toUnicodeRaw(h)
			if err != nil {
				return label, err
			}
			o += uh
			o += string(cp)
			h = ""
		} else {
			h += string(cp)
		}
	}
	uh, err := toUnicodeRaw(h)
	if err != nil {
		return label, err
	}
	o += uh
	return o, nil
}

func toUnicodeRaw(label string) (string, error) {

	original := label

	// Step 1: If all code points in the sequence are in the ASCII range (0..7F) then skip to step 3.
	for i := 0; i < len(label); i++ {
		if label[i] > 127 {
			// Step 2: Perform the steps specified in [NAMEPREP] and fail if there is an error.
			// BUG(dan) We should check error value here
			label, _ = stringprep.Nameprep(label)
			break
		}
	}

	// Step 3: Verify that the sequence begins with the ACE prefix, and save a copy of the sequence.
	if !strings.HasPrefix(label, AcePrefix) {
		return label, errors.New("Label doesn't begin with the ACE prefix")
	} // else

	// 4. Remove the ACE prefix.
	label = strings.SplitN(label, AcePrefix, -1)[1]

	// 5. Decode the sequence using the decoding algorithm in [PUNYCODE] and fail if there is an error.
	//fmt.Printf(label+"\n")
	results, err := punycode.DecodeString(label)

	if err != nil {
		return original, errors.New("Failed punycode decoding: " + err.Error())
	}

	// 6. Apply ToASCII.
	verification, err := ToASCII(label)

	if err != nil {
		return original, errors.New("Failed ToASCII on the decoded sequence: " + err.Error())
	}

	// 7. Verify that the result of step 6 matches the saved copy from step 3,
	// 	  using a case-insensitive ASCII comparison.
	if strings.ToLower(verification) == strings.ToLower(original) {
		return results, nil
	}

	return original, errors.New("Failed verification step")
}

// Returns true if c is a label separator as defined by section 3.1 in RFC 3490
func isSeparator(c rune) bool {
	if c == 0x02E || c == 0x3002 || c == 0xFF0E || c == 0xFF61 {
		return true
	}
	return false
}
