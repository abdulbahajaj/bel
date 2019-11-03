// TODO Note: I am note sure if this is correct please fix if it is not
package reader

import (
    "fmt"
    "strings"
)

type Token struct{
    val string
    name string
    line int
    start int
    end int
    valid bool
    ignore bool
}

type tokenPattern struct{
    name string
    pattern string
    ignore bool
    valid bool
}

func getTokenNames()map[string]string{
    namesList := []string{
        "newLine",
        "comment",
        "openRoundBracket",
        "closeRoundBracket",
        "symbol"
    }
    var namesMap map[string]string

    for name := range namesList{
        namesMap[name] = name
    }

    return namesMap
}

func genTokenPatterns()[]Token{
    names := genTokenNames()
    return []Token{
        Token{name: names["newLine"], pattern: ""},
    }
}

func matchToken(in string) Token, bool{

}

func tokenize(in string){
    // TODO Set a more efficient default size/capacity
    allTokens := make([]Token, 1)

    lineNum = 0
    more := true

    names := getTokenNames()

    for more{
        var token Token
        token, more = matchToken(in)
        if token.name == names['newLine']{
            lineNum++
        }
    }
}
