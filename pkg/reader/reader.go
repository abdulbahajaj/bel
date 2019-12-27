package reader

import (
	"fmt"
	"regexp"
	"errors"
	"strings"
	"strconv"

	"github.com/abdulbahajaj/brutus/pkg/types"
	"github.com/abdulbahajaj/brutus/pkg/common"
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

func removeComments(allTokens []token) []token{
	newAllTokens := make([]token,0)
	skip := false
	for _, token := range allTokens{
		if skip {
			if token.name == "NEW_LINE"{
				skip = false
			}
			continue
		}

		if token.name == "COMMENT" {
			skip = true
			continue
		}

		newAllTokens = append(newAllTokens, token)
	}
	return newAllTokens
}

func tokenize(in string) []token{

	allPatterns := []tokenPattern{
		tokenPattern{ name: "NEW_LINE",  pattern: `\n` },
		tokenPattern{ name: "BACKTICK",  pattern: `` + "`"},
		tokenPattern{ name: "QUOTE",  pattern: `'` },
		tokenPattern{ name: "COMMENT",  pattern: `;`  },
		tokenPattern{ name: "WHITE_SPACE",  pattern: ` `  },
		tokenPattern{ name: "NUMBER",  pattern: `[+-]?([0-9]+(\.[0-9]*)?)`},
		tokenPattern{ name: "OPEN_CIRCLE_BRACKET",   pattern: `\(` },
		tokenPattern{ name: "CLOSE_CIRCLE_BRACKET",  pattern: `\)` },
		tokenPattern{ name: "ESCAPED",  pattern: `(\\bel|\\.)` },
		tokenPattern{ name: "SYMBOL",  pattern: `[-+a-zA-Z,0-9@]*` },
		tokenPattern{ name: "OTHER",  pattern: `.` },
	}

	allPatterns = compilePatterns(allPatterns)

	allTokens := make([]token,0,0)

	for in != ""{
		token, newIn := matchToken(in, allPatterns)
		in = newIn
		allTokens = append(allTokens, token)
	}
	allTokens = removeComments(allTokens)
	return allTokens
}


/*
* Reader helpers
*/

func consume(in []token) (token, []token, error){
	if len(in) == 0 {
		return token{}, []token{}, errors.New("EmptyTokenList")
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
	for _, el := range exp {
		if el.GetType() == types.LIST {
			PrintExp(el.(types.BrutList), intendation + 1)
		} else {
			fmt.Println(indString + el.String())
		}
	}
}

func PrintModule(expStack types.BrutList){
	fmt.Println("===============")
	fmt.Println("Module Begin")
	for _,exp  := range expStack{
		PrintExp(exp.(types.BrutList), 1)
	}
}

func putQuote(child types.BrutType) types.BrutList{
	parent := types.NewBrutList()
	parent = parent.Append(types.BrutSymbol("quote"))
	parent = parent.Append(child)
	return parent
}

func readQuote(allTokens []token) (types.BrutList, []token, error){
	child, remainingTokens, err := readRec(allTokens)
	allTokens = remainingTokens

	this := putQuote(child)
	return this, allTokens, err
}

/*
* Back tick
*/

func putBTWrapper(el types.BrutType) types.BrutType {

	if el.GetType() == types.SYMBOL{
		if sym := el.(types.BrutSymbol); sym[0] == '@'{
			return sym[1:]
		}
	}

	wrapper := types.NewBrutList()
	wrapper = wrapper.Append(types.BrutSymbol("list"))
	wrapper = wrapper.Append(el)

	return wrapper
}

func putBackTick(bType types.BrutType)(types.BrutType){
	if common.IsAtom(bType){
		// TODO Why quote every non symbol atom?
		if bType.GetType() != types.SYMBOL {
			return putQuote(bType)
		} else if sym := bType.(types.BrutSymbol); sym[0] == ','{
			return sym[1:]
		}
		return putQuote(bType)
	} else {
		list := bType.(types.BrutList)
		exp := types.NewBrutList()
		exp = exp.Append(types.BrutSymbol("append"))
		listLength := len(list)
		for cursor := 0; cursor < listLength; cursor++{
			el := list[cursor]
			btEl := putBackTick(el)
			wrappedBtEl := putBTWrapper(btEl)
			exp = exp.Append(wrappedBtEl)

		}
		return exp
	}
}

func readBackTick(allTokens []token)(types.BrutType, []token, error){
	readStructure, remaining_tokens, err := readRec(allTokens)
	allTokens = remaining_tokens
	return putBackTick(readStructure), allTokens, err
}

/*
* Reader
*/

func readSymbol(allTokens []token) (types.BrutSymbol, []token, error){
	current_token, remaining_tokens, _ := consume(allTokens)
	allTokens = remaining_tokens
	return types.BrutSymbol(current_token.val), allTokens, nil
}

func readExp(allTokens []token)(types.BrutList, []token, error){
	exp := types.NewBrutList()
	// _, remaining_tokens, _ := consume(allTokens)
	// allTokens = remaining_tokens
	for {
		current, remaining_tokens, err := consume(allTokens)
		allTokens = remaining_tokens
		if err != nil {
			return exp, []token{}, errors.New("Missing (")
		} else if current.name == "CLOSE_CIRCLE_BRACKET" {
			return exp, remaining_tokens, nil
		} else if current.name == "WHITE_SPACE" {
			continue
		} else if current.name == "NEW_LINE"{
			continue
		} else {
			allTokens = unConsume(current, allTokens)
			result, remaining_tokens, err := readRec(allTokens)
			if err != nil {
				return exp, []token{}, err
			}
			allTokens = remaining_tokens
			exp = exp.Append(result)
		}
	}
	return exp, allTokens, nil
}

func readNum(allTokens []token) (types.BrutNumber, []token, error){
	current_token, remaining_tokens, _ := consume(allTokens)
	allTokens = remaining_tokens
	i, _ := strconv.ParseFloat(current_token.val, 64)
	return types.BrutNumber(i), allTokens, nil
}

func readRec(allTokens []token)(types.BrutType, []token, error){
	current_token, remaining_tokens, err := consume(allTokens)
	allTokens = remaining_tokens
	if err != nil {
		return types.BrutList{}, []token{}, err
	}
	// allTokens = unConsume(current_token, remaining_tokens)

	if current_token.name == "OPEN_CIRCLE_BRACKET" {
		return readExp(allTokens)
	} else if current_token.name == "CLOSE_CIRCLE_BRACKET" {
		return types.BrutList{}, []token{}, errors.New("Unmatched )")
	} else if current_token.name == "NUMBER" {
		allTokens = unConsume(current_token, remaining_tokens)
		return readNum(allTokens)
	} else if current_token.name == "SYMBOL" {
		allTokens = unConsume(current_token, remaining_tokens)
		sym, allTokens, err := readSymbol(allTokens)
		if sym == "nil"{
			return types.NewBrutList(), allTokens, err
		}
		return sym, allTokens, err
	} else if current_token.name == "NEW_LINE" || current_token.name == "WHITE_SPACE"{
		return readRec(remaining_tokens)
	} else if current_token.name == "QUOTE"{
		return readQuote(allTokens)
	} else if current_token.name == "BACKTICK"{
		return readBackTick(allTokens)
	}
	return types.BrutList{}, []token{}, errors.New(
		"Unidentified expression: " + current_token.val + " " + current_token.name)
}

func Read(in string) (types.BrutList, error){
	expression_stack := types.NewBrutList()  //make(types.BrutStack, 0)
	allTokens := tokenize(in)
	for {
		exp, remaining_tokens, err := readRec(allTokens)

		if err != nil {
			return expression_stack, err
		}

		expression_stack = append(expression_stack, (exp.(types.BrutList)))

		if len(remaining_tokens) == 0 {
			break
		}

		allTokens = remaining_tokens
	}
	return expression_stack, nil
}
