package blockchain

type block struct {
	data     string
	hash     string
	prevHash string
}

type blockchain struct {
	blocks []block
}

var b *blockchain

func GenerateBlockchain() *blockchain {
	if b == nil {
		b = &blockchain{}
	}
	return b
}
