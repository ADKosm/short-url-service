package handlers

import "math/rand"

type ShortGenerator interface {
	GenShortURL() string
}

type generator struct {
}

func NewGenerator() *generator {
	return &generator{}
}

func (g *generator) GenShortURL() string {
	alphabet := []byte("qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM1234567890")
	rand.Shuffle(len(alphabet), func(i, j int) {
		alphabet[i], alphabet[j] = alphabet[j], alphabet[i]
	})
	id := string(alphabet[:8])
	return id
}
