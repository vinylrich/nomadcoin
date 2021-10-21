package main

func main() {
	chain := blockchain{}
	chain.addBlock("Genesis Block")
	chain.addBlock("Second Block")
	chain.addBlock("Third Block")
	chain.listofBlocks()
}
