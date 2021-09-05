package generatormock

type ShortGenerator struct {
	GenShortURLFunc func() string
}

func NewGenerator() *ShortGenerator {
	return &ShortGenerator{}
}

func (g *ShortGenerator) GenShortURL() string {
	return g.GenShortURLFunc()
}
