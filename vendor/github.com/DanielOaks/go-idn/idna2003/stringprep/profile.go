// Copyright 2012 Hannes Baldursson. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file is part of go-idn

package stringprep

type Profile []ProfileElement
type ProfileElement struct {
	Step  int // see Step const's
	Table Table
}

const (
	typeMask  = 0xC000 // 11000000 00000000
	valueMask = 0x3FFF // 00111111 11111111
)

type valueType int

const (
	Unassigned valueType = iota
	Map
	Prohibited
	Delete
)
