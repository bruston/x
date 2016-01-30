package main

import "fmt"

func main() {
	hashed := hash([]byte("testing"))
	fmt.Println(hashed)
}

// Reference:
// http://www.burtleburtle.net/bob/hash/doobs.html
// https://en.wikipedia.org/wiki/Jenkins_hash_function
func hash(key []byte) uint32 {
	var hash uint32
	var v uint32
	for _, b := range key {
		v = uint32(b)
		hash += v
		hash += (hash << 10)
		hash ^= (hash >> 6)
	}
	hash += (hash << 3)
	hash ^= (hash >> 11)
	hash += (hash << 15)
	return hash
}
