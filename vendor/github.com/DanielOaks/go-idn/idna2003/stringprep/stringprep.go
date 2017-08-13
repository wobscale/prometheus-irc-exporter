// Copyright 2012 Hannes Baldursson. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file is part of go-idn

// Package stringprep implements Stringprep as described in RFC 3454,
// including the Nameprep profile.
//
// This package is in beta and is still being written and tested.
package stringprep

import (
	"errors"

	"golang.org/x/text/unicode/norm"
)

// Steps in a stringprep profile.
const (
	NFKC                = 1
	BIDI                = 2
	MAP_TABLE           = 3
	UNASSIGNED_TABLE    = 4
	PROHIBIT_TABLE      = 5
	BIDI_PROHIBIT_TABLE = 6
	BIDI_RAL_TABLE      = 7
	BIDI_L_TABLE        = 8
)

// MaxMapChars is the largest number of runes/bytes a mapping will take up.
const MaxMapChars = 4

type d [MaxMapChars]rune

/*

// Append returns p(append(out, b...)). The buffer out must be nil, empty or
// equal to p(out).
func (p *Profile) Append(out []byte, src ...byte) []byte { return nil }

// AppendString returns p(append(out, []byte(s))). The buffer out must be nil,
//empty, or equal to p(out).
func (p *Profile) AppendString(out []byte, src string) []byte { return nil }

// Bytes returns p(b). May return b if p(b) = b.
func (p *Profile) Bytes(b []byte) []byte { return nil }

// Reader returns a new reader that implements Read by reading data from r and
// returning p(data).
func (p *Profile) Reader(r io.Reader) io.Reader { return nil }

// String returns p(s).
func (p *Profile) String(s string) string { return "" }

// Writer returns a new writer that implements Write(b) by writing p(b) to w.
// The returned writer may use an an internal buffer to maintain state across
// Write calls. Calling its Close method writes any buffered data to w.
func (p *Profile) Writer(w io.Writer) io.WriteCloser { return nil }
*/

// PrepareRunes prepares the input rune array according to the stringprep
// profile, and returns the results as a rune array.
func PrepareRunes(profile Profile, input []rune) ([]rune, error) {
	output := make([]rune, len(input))
	copy(output[0:], input[0:])

	for i := 0; i < len(profile); i++ {
		switch profile[i].Step {
		case NFKC:
			// ew, so many conversions here
			output = []rune(string(norm.NFKC.Bytes([]byte(string(output)))))
			break
		case BIDI:
			doneProhibited := 0
			doneRAL := 0
			doneL := 0
			containsRAL := -1
			containsL := -1
			startswithRAL := 0
			endswithRAL := 0

			for j := 0; j < len(profile); j++ {
				switch profile[j].Step {
				case BIDI_PROHIBIT_TABLE:
					doneProhibited = 1
					for k := 0; k < len(output); k++ {
						if in_table(output[k], profile[j].Table) {
							return nil, errors.New("stringprep: BIDI prohibited table")
						}
					}

				case BIDI_RAL_TABLE:
					doneRAL = 1
					for k := 0; k < len(output); k++ {
						if in_table(output[k], profile[j].Table) {
							containsRAL = j
							if k == 0 {
								startswithRAL = 1
							} else if k == len(output)-1 {
								endswithRAL = 1
							}
						}
					}

				case BIDI_L_TABLE:
					doneL = 1
					for k := 0; k < len(output); k++ {
						if in_table(output[k], profile[j].Table) {
							containsL = j
						}
					}
				}
			}

			if doneProhibited != 1 || doneRAL != 1 || doneL != 1 {
				return nil, errors.New("stringprep: Profile error")
			}

			if containsRAL != -1 && containsL != -1 {
				return nil, errors.New("stringprep: BIDI both L and RAL")
			}

			if containsRAL != -1 && (startswithRAL+endswithRAL != 2) {
				return nil, errors.New("stringprep: Contains RAL but does not start and end with RAL characters")
			}

			break
		case MAP_TABLE:
			output = map_table(output, profile[i].Table)
			break
		case UNASSIGNED_TABLE:
			for k := 0; k < len(output); k++ {
				if in_table(output[k], profile[i].Table) {
					return nil, errors.New("stringprep: Unassigned character in input runes")
				}
			}
		case PROHIBIT_TABLE:
			for k := 0; k < len(output); k++ {
				if in_table(output[k], profile[i].Table) {
					return nil, errors.New("stringprep: Prohibited character, cannot casefold this")
				}
			}
			break
		case BIDI_PROHIBIT_TABLE:
			break
		case BIDI_RAL_TABLE:
			break
		case BIDI_L_TABLE:
			break
		default:
			return nil, errors.New("stringprep: Profile error")
		}
	}

	return output, nil
}
