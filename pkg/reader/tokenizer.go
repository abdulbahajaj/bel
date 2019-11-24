// TODO Note: I am note sure if this is correct please fix if it is not
package reader

import (
    "fmt"
    "regexp"
    "errors"
)

type tokenName int
const (
    //Compiler specific token-types
    NEWLINE tokenName=iota
    COMMENT
    INVALID
    WHITE_SPACE

    //Bel specific token-types
    OPEN_CIRCLE_BRACKET
    CLOSE_CIRCLE_BRACKET
    NUMBER
    CHARACHTER
    STRING
    WORD // Anything other than the above, could be symbol,
         // special form, etc
)

type token struct{
    name tokenName
    val string
    line int
    start int
    end int
}

func (token token) String() string {

    var tokenName string

    if token.name == NEWLINE{
        tokenName = "NEWLINE"
    } else if token.name == COMMENT{
        tokenName = "COMMENT"
    } else if token.name == OPEN_CIRCLE_BRACKET{
        tokenName = "OPEN_CIRCLE_BRACKET"
    } else if token.name == CLOSE_CIRCLE_BRACKET{
        tokenName = "CLOSE_CIRCLE_BRACKET"
    } else if token.name == INVALID{
        tokenName = "INVALID"
    } else if token.name == WHITE_SPACE{
        tokenName = "WHITE_SPACE"
    } else if token.name == NUMBER{
        tokenName = "NUMBER"
    } else if token.name == CHARACHTER{
        tokenName = "CHARACHTER"
    } else if token.name == STRING{
        tokenName = "STRING"
    } else if token.name == WORD{
        tokenName = "WORD"
    }

    return fmt.Sprintf("type: %v, val: %v, line: %v, start: %v, end: %v",
        tokenName, token.val, token.line, token.start, token.end)
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


func printAllTokens(allTokens []token){
    for _, token := range allTokens {
        fmt.Println(token)
    }
}

func cleanTokens(allTokens []token) ([]token, error){
    cleanedTokens := make([]token, 0,len(allTokens))

    for _, token := range allTokens {
        if token.name == NEWLINE{
            continue
        } else if token.name == COMMENT{
            continue
        } else if token.name == INVALID{
            return cleanedTokens, errors.New(fmt.Sprintf("Invalid token <%s>", token.String()))
        }
        cleanedTokens = append(cleanedTokens, token)
    }
    return cleanedTokens, nil
}

func tokenize(in string) []token{
    allPatterns := []tokenPattern{
        tokenPattern{ name: COMMENT,  pattern: `;`  },
        tokenPattern{ name: STRING,  pattern: `\".*?\"`},
        tokenPattern{ name: WHITE_SPACE,  pattern: ` `  },
        tokenPattern{ name: CHARACHTER ,  pattern: `\\(bel|[a-z])`},
        tokenPattern{ name: NUMBER,  pattern: `[+-]?([0-9]+(\.[0-9]*)?)`},
        tokenPattern{ name: OPEN_CIRCLE_BRACKET,   pattern: `\(` },
        tokenPattern{ name: CLOSE_CIRCLE_BRACKET,  pattern: `\)` },
        tokenPattern{ name: WORD,  pattern: `[^ ]+` },
        tokenPattern{ name: INVALID,  pattern: `.` },
    }

    allPatterns = compilePatterns(allPatterns)

    allTokens := make([]token,0,1)

    for in != ""{
        token, newIn := matchToken(in, allPatterns)
        in = newIn
        allTokens = append(allTokens, token)
    }

    cleanedTokens, err := cleanTokens(allTokens)

    if err != nil {
        panic(err)
    }

    return cleanedTokens
}
