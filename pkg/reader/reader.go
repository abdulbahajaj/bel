package reader


func readList(){

}

func readAtom(){

}

func Read(in string){
	allTokens := tokenize(in)
	// allTokens,  = cleanTokens(allTokens)
	printAllTokens(allTokens)
}
