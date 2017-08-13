// Copyright 2012 Hannes Baldursson. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package stringprep

// An Iter iterates over a string or byte slice, while preparing it with a
// given Profile.
type Iter struct{}

// Done returns true if there is no more input to process.
func (i *Iter) Done() bool { return false }

// Next writes p(i.input[i.Pos():n]...) to buffer buf, where n is the largest
// boundary of i.input such that the result fits in buf. It returns the number 
// of bytes written to buf. len(buf) should be at least MaxSegmentSize. Done 
// must be false before calling Next.
func (i *Iter) Next(buf []byte) int { return -1 }

// Pos returns the byte position at which the next call to Next will commence 
// processing.
func (i *Iter) Pos() int { return -1 }

// SetInput initializes i to iterate over src after normalizing it to 
// Profile p.
func (i *Iter) SetInput(p *Profile, src []byte) {}

// SetInputString initializes i to iterate over src after normalizing it to 
// Profile p.
func (i *Iter) SetInputString(p *Profile, src string) {}
