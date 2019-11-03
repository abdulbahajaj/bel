// TODO Note: I am note sure if this is correct please fix if it is not
package reader

import (
    "fmt"
    "regexp"
)

type tokenName int
const (
    NEWLINE tokenName=iota
    COMMENT
    OPEN_CIRCLE_BRACKET
    CLOSE_CIRCLE_BRACKET
    INVALID

)

type token struct{
    name tokenName
    val string
    line int
    start int
    end int
}

type tokenPattern struct{
    name tokenName
    pattern string
    compiledPattern *regexp.Regexp
}

func compilePatterns(allPatterns []tokenPattern) []tokenPattern{
    compiledPatterns := make([]tokenPattern, 0, len(allPatterns))
    for _, pattern := range allPatterns{
        pattern.pattern = "^" + pattern.pattern
        pattern.compiledPattern, _= regexp.Compile(pattern.pattern)
        compiledPatterns = append(compiledPatterns, pattern)
    }
    return compiledPatterns
}

func matchToken(in string, allPatterns []tokenPattern) (token, string){

    for _, pattern := range allPatterns {
        match := pattern.compiledPattern.FindString(in)
        if match != "" {
            newIn := in[len(match):]
            return token{name: pattern.name, val: match}, newIn
        }
    }
    var t token
    var s string
    return t, s
}

func printToken(token token){
    if token.name == NEWLINE{
        fmt.Print("NEWLINE: ")
    }
    if token.name == COMMENT{
        fmt.Print("COMMENT: ")
    }
    if token.name == OPEN_CIRCLE_BRACKET{
        fmt.Print("OPEN_CIRCLE_BRACKET: ")
    }
    if token.name == CLOSE_CIRCLE_BRACKET{
        fmt.Print("CLOSE_CIRCLE_BRACKET: ")
    }
    if token.name == INVALID{
        fmt.Print("INVALID: ")
    }
    fmt.Println(token.val)
}

func printAllTokens(allTokens []token){
    for _, token := range allTokens {
        printToken(token)
    }
}

func tokenize(in string) []token{
    allPatterns := []tokenPattern{
        tokenPattern{ name: OPEN_CIRCLE_BRACKET,   pattern: `\(` },
        tokenPattern{ name: CLOSE_CIRCLE_BRACKET,  pattern: `\)` },
        tokenPattern{ name: COMMENT,  pattern: `;.*`  },
        tokenPattern{ name: INVALID,  pattern: `.` },
    }
    allPatterns = compilePatterns(allPatterns)

    allTokens := make([]token,0,1)

    for in != ""{
        token, newIn := matchToken(in, allPatterns)
        in = newIn
        allTokens = append(allTokens, token)
    }

    printAllTokens(allTokens)
    return allTokens
}
