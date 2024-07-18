package main

import (
	"crypto/sha256"
	"encoding/hex"
	"math/rand"
	"time"
)

// GenerateUniqueHash generates a unique hash based on the current time and random bytes
// * I think it should be shorter but it's good so there is no conflicts
func GenerateUniqueHash() string {
	timeBytes := []byte(time.Now().Format("2006-01-02 15:04:05.000"))
	hash := sha256.New()
	hash.Write(timeBytes)
	b := make([]byte, 10)
	for i := 0; i < len(b); i++ {
		b[i] = byte(rand.Intn(255))
	}
	hash.Write(b)
	return hex.EncodeToString(hash.Sum(nil))
}
