package tokenizer

import (
	"bytes"
	"fmt"
	"os"
	"strconv"
)

type TokenType int

const (
	EOF TokenType = iota
	LEFT_PAREN
	RIGHT_PAREN
	LEFT_BRACE
	RIGHT_BRACE
	COMMA
	DOT
	MINUS
	PLUS
	SEMICOLON
	STAR
	EQUAL
	EQUAL_EQUAL
	BANG
	BANG_EQUAL
	LESS
	LESS_EQUAL
	GREATER
	GREATER_EQUAL
	SLASH
	STRING
	NUMBER
	IDENTIFIER

	AND
	CLASS
	ELSE
	FALSE
	FOR
	FUN
	IF
	NIL
	OR
	PRINT
	RETURN
	SUPER
	THIS
	TRUE
	VAR
	WHILE
)

var tokens = map[TokenType]Token{
	EOF:           {EOF, "EOF", "", nil},
	LEFT_PAREN:    {LEFT_PAREN, "LEFT_PAREN", "(", nil},
	RIGHT_PAREN:   {RIGHT_PAREN, "RIGHT_PAREN", ")", nil},
	LEFT_BRACE:    {LEFT_BRACE, "LEFT_BRACE", "{", nil},
	RIGHT_BRACE:   {RIGHT_BRACE, "RIGHT_BRACE", "}", nil},
	COMMA:         {COMMA, "COMMA", ",", nil},
	DOT:           {DOT, "DOT", ".", nil},
	MINUS:         {MINUS, "MINUS", "-", nil},
	PLUS:          {PLUS, "PLUS", "+", nil},
	SEMICOLON:     {SEMICOLON, "SEMICOLON", ";", nil},
	STAR:          {STAR, "STAR", "*", nil},
	EQUAL:         {EQUAL, "EQUAL", "=", nil},
	EQUAL_EQUAL:   {EQUAL_EQUAL, "EQUAL_EQUAL", "==", nil},
	BANG:          {BANG, "BANG", "!", nil},
	BANG_EQUAL:    {BANG_EQUAL, "BANG_EQUAL", "!=", nil},
	LESS:          {LESS, "LESS", "<", nil},
	LESS_EQUAL:    {LESS_EQUAL, "LESS_EQUAL", "<=", nil},
	GREATER:       {GREATER, "GREATER", ">", nil},
	GREATER_EQUAL: {GREATER_EQUAL, "GREATER_EQUAL", ">=", nil},
	SLASH:         {SLASH, "SLASH", "/", nil},
	STRING:        {STRING, "STRING", "", ""},
	NUMBER:        {NUMBER, "NUMBER", "", ""},
	IDENTIFIER:    {IDENTIFIER, "IDENTIFIER", "", nil},

	AND:    {AND, "AND", "and", nil},
	CLASS:  {CLASS, "CLASS", "class", nil},
	ELSE:   {ELSE, "ELSE", "else", nil},
	FALSE:  {FALSE, "FALSE", "false", nil},
	FOR:    {FOR, "FOR", "for", nil},
	FUN:    {FUN, "FUN", "fun", nil},
	IF:     {IF, "IF", "if", nil},
	NIL:    {NIL, "NIL", "nil", nil},
	OR:     {OR, "OR", "or", nil},
	PRINT:  {PRINT, "PRINT", "print", nil},
	RETURN: {RETURN, "RETURN", "return", nil},
	SUPER:  {SUPER, "SUPER", "super", nil},
	THIS:   {THIS, "THIS", "this", nil},
	TRUE:   {TRUE, "TRUE", "true", nil},
	VAR:    {VAR, "VAR", "var", nil},
	WHILE:  {WHILE, "WHILE", "while", nil},
}

var reserved = map[string]Token{
	"and":    tokens[AND],
	"class":  tokens[CLASS],
	"else":   tokens[ELSE],
	"false":  tokens[FALSE],
	"for":    tokens[FOR],
	"fun":    tokens[FUN],
	"if":     tokens[IF],
	"nil":    tokens[NIL],
	"or":     tokens[OR],
	"print":  tokens[PRINT],
	"return": tokens[RETURN],
	"super":  tokens[SUPER],
	"this":   tokens[THIS],
	"true":   tokens[TRUE],
	"var":    tokens[VAR],
	"while":  tokens[WHILE],
}

type Token struct {
	tokenType     TokenType
	tokenTypeName string
	lexeme        string
	literal       interface{}
}

func (t Token) String() string {
	if t.literal == nil {
		return fmt.Sprintf("%s %s null", t.tokenTypeName, t.lexeme)
	}

	return fmt.Sprintf("%s %s %v", t.tokenTypeName, t.lexeme, t.literal)
}

type LexingError struct {
	line    int
	message string
	lexeme  string
}

func (e LexingError) String() string {
	if e.lexeme == "" {
		return fmt.Sprintf("[line %d] Error: %s.", e.line, e.message)
	}
	return fmt.Sprintf("[line %d] Error: %s: %s", e.line, e.message, e.lexeme)
}

func matchEqualVariants(tokens []Token, fileContents []byte, cursor int, singleToken Token, dualToken Token) (int, []Token) {
	if cursor < len(fileContents) && fileContents[cursor] == '=' {
		tokens = append(tokens, dualToken)
		cursor++
	} else {
		tokens = append(tokens, singleToken)
	}

	return cursor, tokens
}

func Tokenize(fileContents []byte) ([]Token, error) {
	parsedTokens := []Token{}
	errors := []LexingError{}

	currentLine := 0
	for cursor := 0; cursor < len(fileContents); {
		token := fileContents[cursor]
		cursor++
		if token == '\n' {
			currentLine++
			continue
		}

		if token == ' ' || token == '\r' || token == '\t' {
			continue
		}

		switch token {
		case '(':
			parsedTokens = append(parsedTokens, tokens[LEFT_PAREN])
		case ')':
			parsedTokens = append(parsedTokens, tokens[RIGHT_PAREN])
		case '{':
			parsedTokens = append(parsedTokens, tokens[LEFT_BRACE])
		case '}':
			parsedTokens = append(parsedTokens, tokens[RIGHT_BRACE])
		case ',':
			parsedTokens = append(parsedTokens, tokens[COMMA])
		case '.':
			parsedTokens = append(parsedTokens, tokens[DOT])
		case '-':
			parsedTokens = append(parsedTokens, tokens[MINUS])
		case '+':
			parsedTokens = append(parsedTokens, tokens[PLUS])
		case ';':
			parsedTokens = append(parsedTokens, tokens[SEMICOLON])
		case '*':
			parsedTokens = append(parsedTokens, tokens[STAR])
		case '=':
			cursor, parsedTokens = matchEqualVariants(parsedTokens, fileContents, cursor, tokens[EQUAL], tokens[EQUAL_EQUAL])
		case '!':
			cursor, parsedTokens = matchEqualVariants(parsedTokens, fileContents, cursor, tokens[BANG], tokens[BANG_EQUAL])
		case '<':
			cursor, parsedTokens = matchEqualVariants(parsedTokens, fileContents, cursor, tokens[LESS], tokens[LESS_EQUAL])
		case '>':
			cursor, parsedTokens = matchEqualVariants(parsedTokens, fileContents, cursor, tokens[GREATER], tokens[GREATER_EQUAL])
		case '/':
			if cursor < len(fileContents) && fileContents[cursor] == '/' {
				newLineCursor := bytes.Index(fileContents[cursor:], []byte("\n"))
				if newLineCursor == -1 {
					cursor = len(fileContents)
				} else {
					cursor += newLineCursor
				}
			} else {
				parsedTokens = append(parsedTokens, tokens[SLASH])
			}
		case '"':
			start := cursor
			end := bytes.Index(fileContents[start:], []byte{'"'})
			if end == -1 {
				error := LexingError{currentLine + 1, "Unterminated string", ""}
				errors = append(errors, error)
				fmt.Fprintln(os.Stderr, error)
				cursor = len(fileContents)
				continue
			}

			value := string(fileContents[start : start+end])
			valueWithQuotes := string(fileContents[start-1 : start+end+1])
			stringToken := tokens[STRING]
			stringToken.lexeme = valueWithQuotes
			stringToken.literal = value
			parsedTokens = append(parsedTokens, stringToken)

			cursor += end + 1
		default:
			if token >= '0' && token <= '9' {
				start := cursor - 1
				dotSeen := false
				for cursor < len(fileContents) &&
					(fileContents[cursor] >= '0' &&
						fileContents[cursor] <= '9' ||
						fileContents[cursor] == '.') {
					if fileContents[cursor] == '.' && dotSeen {
						break
					}

					if fileContents[cursor] == '.' {
						dotSeen = true
					}

					cursor++
				}

				if fileContents[cursor-1] == '.' {
					cursor--
					dotSeen = false
				}

				lexeme := string(fileContents[start:cursor])

				num, err := strconv.ParseFloat(lexeme, 64)
				if err != nil {
					error := LexingError{currentLine + 1, "Invalid number", lexeme}
					errors = append(errors, error)
					fmt.Fprintln(os.Stderr, error)
					break
				}

				numberToken := tokens[NUMBER]
				numberToken.lexeme = lexeme
				numberToken.literal = formatFloat(num)
				parsedTokens = append(parsedTokens, numberToken)
				break
			}

			if (token >= 'a' && token <= 'z') || (token >= 'A' && token <= 'Z') || token == '_' {
				start := cursor - 1
				for cursor < len(fileContents) &&
					((fileContents[cursor] >= 'a' && fileContents[cursor] <= 'z') ||
						(fileContents[cursor] >= 'A' && fileContents[cursor] <= 'Z') ||
						(fileContents[cursor] >= '0' && fileContents[cursor] <= '9') ||
						fileContents[cursor] == '_') {
					cursor++
				}

				token, ok := reserved[string(fileContents[start:cursor])]
				if ok {
					parsedTokens = append(parsedTokens, token)
					break
				}

				identifierToken := tokens[IDENTIFIER]
				identifierToken.lexeme = string(fileContents[start:cursor])
				parsedTokens = append(parsedTokens, identifierToken)
				break
			}

			error := LexingError{currentLine + 1, "Unexpected character", string(token)}
			errors = append(errors, error)
			fmt.Fprintln(os.Stderr, error)
		}
	}

	parsedTokens = append(parsedTokens, tokens[EOF])

	if len(errors) > 0 {
		return parsedTokens, fmt.Errorf("lexing errors found")
	}

	return parsedTokens, nil
}
