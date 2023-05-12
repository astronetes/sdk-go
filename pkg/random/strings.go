package random

import (
	"math/rand"
	"time"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandSeq(n int) string {
	return RandSeqWithCharacters(n, letters)
}

func RandSeqWithCharacters(n int, characters []rune) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[seededRand.Intn(len(letters))]
	}
	return string(b)
}
