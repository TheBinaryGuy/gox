package tokenizer

import (
	"fmt"
	"strconv"
	"strings"
)

func PrintTokens(tokens []Token) {
	for _, token := range tokens {
		fmt.Printf("%s\n", token)
	}
}

func formatFloat(num float64) string {
	formatted := strconv.FormatFloat(num, 'f', -1, 64)
	if !strings.Contains(formatted, ".") {
		formatted += ".0"
	}
	return formatted
}
