// TODO Note: I am note sure if this is correct please fix if it is not
package reader

import (
  "io"
  "bytes"
  "bufio"
)

type TokenType int
const (
  ILLEGAL TokenType = iota
  EOF
  WS

  OPEN_PARENTHESE
  CLOSE_PARENTHESE

  SYMBOL
  PLUS
  MINUS
  NUMBER
)

var eof = rune(0)

type Token struct {
  Type TokenType
  Lit string
}

type Scanner struct {
  src *bufio.Reader
}

func NewScanner(r io.Reader) *Scanner {
  return &Scanner{src: bufio.NewReader(r)}
}

func (s *Scanner) read() rune {
  ch, _, err := s.src.ReadRune()
  if err != nil {
    return eof
  }
  return ch
}
// Used for look ahead
func (s *Scanner) unread() error { return s.src.UnreadRune() }
func isWhitespace(ch rune) bool {
  return ch == ' ' || ch == '\t' || ch == '\n'
}
func isLetter(ch rune) bool {
  return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')
}
func isDigit(ch rune) bool {
  return (ch >= '0' && ch <= '9')
}

func (s *Scanner) scanWhitespace() (tok TokenType, lit string) {
  var buf bytes.Buffer
  buf.WriteRune(s.read())

  for {
    if ch := s.read(); ch == eof {
      break
    } else if !isWhitespace(ch) {
      s.unread()
      break
    } else {
      buf.WriteRune(ch)
    }
  }

  return WS, buf.String()
}

func (s *Scanner) scanNumber() (tok TokenType, lit string) {
  var buf bytes.Buffer
  buf.WriteRune(s.read())

  for {
    if ch := s.read(); ch == eof {
      break
    } else if !isDigit(ch) && ch != '.' {
      s.unread()
      break
    } else {
      buf.WriteRune(ch)
    }
  }

  // Otherwise return as a regular identifier.
  return NUMBER, buf.String()
}

func (s *Scanner) scanSymbol() (tok TokenType, lit string) {
  var buf bytes.Buffer
  buf.WriteRune(s.read())

  for {
    if ch := s.read(); ch == eof {
      break
    } else if !isLetter(ch) && !isDigit(ch) {
      s.unread()
      break
    } else {
      buf.WriteRune(ch)
    }
  }

  return SYMBOL, buf.String()
}

func (s *Scanner) Scan() (tok TokenType, lit string) {
  ch := s.read()

  // Multiple line token
  if isWhitespace(ch) {
    s.unread()
    return s.scanWhitespace()
  } else if isDigit(ch) || ch == '.' { // .0 would be considered 0.0
    s.unread()
    return s.scanNumber()
  } else if isLetter(ch) {
    s.unread()
    return s.scanSymbol()
  }

  // individual character token
  switch ch {
  case eof:
    return EOF, ""
  case '+':
    return PLUS, string(ch)
  case '-':
    return MINUS, string(ch)
  case '(':
    return OPEN_PARENTHESE, string(ch)
  case ')':
    return CLOSE_PARENTHESE, string(ch)
  }

  return ILLEGAL, string(ch)
}

func PraseTokens(r io.Reader) []Token {
  s := NewScanner(r)

  tokens := []Token{}
  for {
    tok, lit := s.Scan()
    tokens = append(tokens, Token{tok, lit})

    if tok == EOF {
      break
    }
  }
  
  return tokens
}

