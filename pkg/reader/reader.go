package reader

import (
	"fmt"
	"regexp"
	"github.com/abdulbahajaj/brutus/pkg/types"
	"errors"
	"strings"
	"strconv"
)


/*
* Tokenizer
*/

type token struct{
	name string
	val string
	line int
	start int
	end int
}

func (token token) String() string {
	return fmt.Sprintf("type: %v, val: %v, line: %v, start: %v, end: %v",
		token.name, token.val, token.line, token.start, token.end)
}

type tokenPattern struct{
	name string
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

func tokenize(in string) []token{
	allPatterns := []tokenPattern{
		tokenPattern{ name: "COMMENT",  pattern: `;`  },
		tokenPattern{ name: "WHITE_SPACE",  pattern: ` `  },
		tokenPattern{ name: "NUMBER",  pattern: `[+-]?([0-9]+(\.[0-9]*)?)`},
		tokenPattern{ name: "OPEN_CIRCLE_BRACKET",   pattern: `\(` },
		tokenPattern{ name: "CLOSE_CIRCLE_BRACKET",  pattern: `\)` },
		tokenPattern{ name: "ESCAPED",  pattern: `(\\bel|\\.)` },
		tokenPattern{ name: "SYMBOL",  pattern: `[^ ]*` },
		tokenPattern{ name: "OTHER",  pattern: `.` },
	}

	allPatterns = compilePatterns(allPatterns)

	allTokens := make([]token,0,0)

	for in != ""{
		token, newIn := matchToken(in, allPatterns)
		in = newIn
		allTokens = append(allTokens, token)
	}

	return allTokens
}


/*
* Reader
*/

func consume(in []token) (token, []token, error){
	if len(in) == 0 {
		return token{}, []token{}, errors.New("Empty token list")
	} else if len(in) == 1 {
		return in[0], []token{}, nil
	}

	return in[0], in[1:], nil
}

func unConsume(oneToken token ,in []token) []token{
	return append([]token{oneToken}, in...)
}

func PrintExp(exp types.BrutList, intendation int){
	indString := strings.Repeat("\t", intendation)
	fmt.Println(indString + " ===============")
	fmt.Println(indString + "Exp begin")
	for _, el := range exp.Elements {
		if el.GetType() == types.LIST {
			PrintExp(el.(types.BrutList), intendation + 1)
		} else {
			fmt.Println(indString + el.String())
		}
	}
}

func PrintModule(module types.BrutModule){
	fmt.Println("===============")
	fmt.Println("Module Begin")
	for _,exp  := range module.Expressions{
		PrintExp(exp, 1)
	}
}

func readSymbol(allTokens []token) (types.BrutSymbol, []token, error){
	current_token, remaining_tokens, _ := consume(allTokens)
	allTokens = remaining_tokens
	return types.NewBrutSymbol(current_token.val), allTokens, nil
}

func readExp(allTokens []token)(types.BrutList, []token, error){
	exp := types.NewBrutList()
	_, remaining_tokens, _ := consume(allTokens)
	allTokens = remaining_tokens
	for {
		current, remaining_tokens, err := consume(allTokens)
		allTokens = remaining_tokens
		if err != nil {
			return exp, []token{}, errors.New("Missing (")
		} else if current.name == "CLOSE_CIRCLE_BRACKET" {
			return exp, remaining_tokens, nil
		} else if current.name == "WHITE_SPACE" {
			continue
		} else {
			allTokens = unConsume(current, allTokens)
			result, remaining_tokens, err := Read(allTokens)
			if err != nil {
				return exp, []token{}, err
			}
			allTokens = remaining_tokens
			exp = exp.Append(result)
		}
	}
	return exp, allTokens, nil
}

// TODO return the real number
func readNum(allTokens []token) (types.BrutNumber, []token, error){
	current_token, remaining_tokens, _ := consume(allTokens)
	allTokens = remaining_tokens
	i, _ := strconv.ParseFloat(current_token.val, 64)
	return types.NewBrutNumber(i), allTokens, nil
}

func Read(allTokens []token)(types.BrutType, []token, error){
	current_token, remaining_token, err := consume(allTokens)

	if err != nil {
		return types.BrutList{}, []token{}, err
	}
	allTokens = unConsume(current_token, remaining_token)

	if current_token.name == "OPEN_CIRCLE_BRACKET" {
		return readExp(allTokens)
	} else if current_token.name == "CLOSE_CIRCLE_BRACKET" {
		return types.BrutList{}, []token{}, errors.New("Unmatched )")
	} else if current_token.name == "NUMBER" {
		return readNum(allTokens)
	} else if current_token.name == "SYMBOL" {
		return readSymbol(allTokens)
	}

	return types.BrutList{}, []token{}, errors.New(
		"Unidentified expression: " + current_token.val + " " + current_token.name)
}

func ReadModule(in string) (types.BrutModule, error){
	module := types.NewBrutModule()
	allTokens := tokenize(in)
	for {
		exp, remaining_tokens, err := Read(allTokens)

		if err != nil {
			return module, err
		}

		module = module.AppendExp(exp.(types.BrutList))

		if len(remaining_tokens) == 0 {
			break
		}

		allTokens = remaining_tokens
	}
	return module, nil
}
