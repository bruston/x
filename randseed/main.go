// Yoinked from https://groups.google.com/forum/#!searchin/golang-nuts/rand$20seed/golang-nuts/NsxEM9T5c4A/hd0XalAYFvoJ
// Most examples show seeding the PRNG with the time in nanoseconds, which has poor entropy.

package main

import (
	crand "crypto/rand"
	"encoding/binary"
	"math/rand"
)

func seed() error {
	var s int64
	err := binary.Read(crand.Reader, binary.BigEndian, &s)
	if err != nil {
		return err
	}
	rand.Seed(s)
	return nil
}

func init() {
	seed()
}
