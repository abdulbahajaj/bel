package main

import (
  "fmt"
  "strings"

  "github.com/abdulbahajaj/bel/pkg/reader"
)

var tt = map[reader.TokenType]string {
	reader.ILLEGAL: "ILLEGAL",
  reader.EOF: "EOF",
  reader.WS: "WS",

  reader.OPEN_PARENTHESE: "OPEN_PARENTHESE",
  reader.CLOSE_PARENTHESE: "CLOSE_PARENTHESE",

  reader.SYMBOL: "SYMBOL",
  reader.PLUS: "PLUS",
  reader.MINUS: "MINUS",
  reader.NUMBER: "NUMBER",
}

func main() {
  TestString := "(+ -1 2)"
	

	tokens := reader.PraseTokens(strings.NewReader(TestString))
	
	fmt.Println(len(tokens))
	for _, t := range tokens {
		fmt.Printf("%s: '%s'\n", tt[t.Type] ,t.Lit)	
	}
  
}
