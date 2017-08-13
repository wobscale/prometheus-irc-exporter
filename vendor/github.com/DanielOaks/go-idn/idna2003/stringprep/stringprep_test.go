// Copyright 2012 Hannes Baldursson. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file is part of go-idn

package stringprep

import "testing"

type mappingtestcase struct {
	Mapping string
	Input   []rune
	Output  []rune
}

// from http://tools.ietf.org/html/draft-josefsson-idn-test-vectors-00#section-4
var mappingTests = []mappingtestcase{
	{"nameprep", []rune{0x0066, 0x006f, 0x006f, 0x00ad, 0x034f, 0x1806, 0x180b, 0x0062, 0x0061, 0x0072, 0x200b, 0x2060, 0x0062, 0x0061, 0x007a, 0xfe00, 0xfe08, 0xfe0f, 0xfeff}, []rune{0x0066, 0x006f, 0x006f, 0x0062, 0x0061, 0x0072, 0x0062, 0x0061, 0x007a}},
	{"nameprep", []rune{0x0043, 0x0041, 0x0046, 0x0045}, []rune{0x0063, 0x0061, 0x0066, 0x0065}},
	{"nameprep", []rune{0x00df}, []rune{0x0073, 0x0073}},
	{"nameprep", []rune{0x0130}, []rune{0x0069, 0x0307}},
	{"nameprep", []rune{0x0143, 0x037a}, []rune{0x0144, 0x0020, 0x03b9}},
	{"nameprep", []rune{0x2121, 0x33c6, 0x1d7bb}, []rune{0x0074, 0x0065, 0x006c, 0x0063, 0x2215, 0x006b, 0x0067, 0x03c3}},
	{"nameprep", []rune{0x006a, 0x030c, 0x00a0, 0x00aa}, []rune{0x01f0, 0x0020, 0x0061}},
	{"nameprep", []rune{0x1fb7}, []rune{0x1fb6, 0x03b9}},
	{"nameprep", []rune{0x01f0}, []rune{0x01f0}},
	{"nameprep", []rune{0x0390}, []rune{0x0390}},
	{"nameprep", []rune{0x03b0}, []rune{0x03b0}},
	{"nameprep", []rune{0x1e96}, []rune{0x1e96}},
	{"nameprep", []rune{0x1f56}, []rune{0x1f56}},
	{"nameprep", []rune{0x0020}, []rune{0x0020}},
	{"nameprep", []rune{0x00a0}, []rune{0x0020}},
	{"nameprep", []rune{0x2000}, []rune{0x0020}},
	{"nameprep", []rune{0x200b}, []rune{}},
	{"nameprep", []rune{0x3000}, []rune{0x0020}},
	{"nameprep", []rune{0x0010, 0x007f}, []rune{0x0010, 0x007f}},
	{"nameprep", []rune{0xfeff}, []rune{}},
	{"nameprep", []rune{0x0341}, []rune{0x0301}},
	{"nameprep", []rune{0x0066, 0x006f, 0x006f, 0xfe76, 0x0062, 0x0061, 0x0072}, []rune{0x0066, 0x006f, 0x006f, 0x0020, 0x064e, 0x0062, 0x0061, 0x0072}},
	{"nameprep", []rune{0x0627, 0x0031, 0x0628}, []rune{0x0627, 0x0031, 0x0628}},
	{"nameprep", []rune{0x0058, 0x00ad, 0x00df, 0x0130, 0x2121, 0x006a, 0x030c, 0x00a0, 0x00aa, 0x03b0, 0x2000}, []rune{0x0078, 0x0073, 0x0073, 0x0069, 0x0307, 0x0074, 0x0065, 0x006c, 0x01f0, 0x0020, 0x0061, 0x03b0, 0x0020}},
	{"nameprep", []rune{0x0058, 0x00df, 0x3316, 0x0130, 0x2121, 0x249f, 0x3300}, []rune{0x0078, 0x0073, 0x0073, 0x30ad, 0x30ed, 0x30e1, 0x30fc, 0x30c8, 0x30eb, 0x0069, 0x0307, 0x0074, 0x0065, 0x006c, 0x0028, 0x0064, 0x0029, 0x30a2, 0x30d1, 0x30fc, 0x30c8}},
}

type badmappingtestcase struct {
	Mapping string
	Input   []rune
}

// from http://tools.ietf.org/html/draft-josefsson-idn-test-vectors-00#section-4
var badMappingTests = []badmappingtestcase{
	{"nameprep", []rune{0x1680}},
	{"nameprep", []rune{0x0085}},
	{"nameprep", []rune{0x180e}},
	{"nameprep", []rune{0x1d175}},
	{"nameprep", []rune{0xf123}},
	{"nameprep", []rune{0xf1234}},
	{"nameprep", []rune{0x10f234}},
	{"nameprep", []rune{0x8fffe}},
	{"nameprep", []rune{0x10ffff}},
	{"nameprep", []rune{0xdf42}},
	{"nameprep", []rune{0xfffd}},
	{"nameprep", []rune{0x2ff5}},
	{"nameprep", []rune{0x200e}},
	{"nameprep", []rune{0x202a}},
	{"nameprep", []rune{0xe0001}},
	{"nameprep", []rune{0xe0042}},
	{"nameprep", []rune{0x0066, 0x006f, 0x006f, 0x05be, 0x0062, 0x0061, 0x0072}},
	{"nameprep", []rune{0x0066, 0x006f, 0x006f, 0xfd50, 0x0062, 0x0061, 0x0072}},
	{"nameprep", []rune{0x0627, 0x0031}},
	{"nameprep", []rune{0xe0002}},
}

func TestInTable(t *testing.T) {
	c := rune(0x000221)
	if !in_table(c, _A1) {
		t.Errorf("in_table(0x000221, _A1) = false; want true")
	}

	d := rune(0x000220)
	if in_table(d, _A1) {
		t.Errorf("in_table(0x000220, _A1) = true; want false")
	}
}

func TestMapping(t *testing.T) {
	for i, test := range mappingTests {
		output, err := PrepareRunes(Profiles[test.Mapping], test.Input)
		if err != nil {
			t.Error(
				"For test", i, test.Input,
				"got Error", err.Error(),
			)
		}
		if string(output) != string(test.Output) {
			t.Error(
				"For test", i, test.Input,
				"expected", test.Output,
				"got", output,
			)
		}
	}
}

func TestBadMapping(t *testing.T) {
	for i, test := range badMappingTests {
		output, err := PrepareRunes(Profiles[test.Mapping], test.Input)
		if err == nil || output != nil {
			t.Error(
				"For test", i, test.Input,
				"did not get Error",
			)
		}
	}
}
