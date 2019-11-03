package reader_test

import (
	"testing"
	"strings"

	"github.com/stretchr/testify/assert"

	r "github.com/abdulbahajaj/bel/pkg/reader"
)

func TestPraseTokens(t *testing.T) {
	testCases := []struct{
		input string
		expected []r.Token
	} {
		{
			"(+ 1 2)",
			[]r.Token{
				{r.OPEN_PARENTHESE, "("},
				{r.PLUS, "+"},
				{r.WS, " "},
				{r.NUMBER, "1"},
				{r.WS, " "},
				{r.NUMBER, "2"},
				{r.CLOSE_PARENTHESE, ")"},
				{r.EOF, ""},
			},
		},
	}

	
	for _, testCase := range testCases {
		actual := r.PraseTokens(strings.NewReader(testCase.input))
		assert.Equal(t, len(testCase.expected), len(actual), "Miss matching number of tokens")
	}
}