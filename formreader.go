package main

import (
	"bufio"
	"io"
  "unicode/utf8"
	"unicode"
)

func IsWhitespace(r rune) bool {
  return unicode.IsSpace(r) || r == ','
}

type Input struct {
  *bufio.Scanner
}

func NewInput(in io.Reader) *Input {
  input := &Input{
    bufio.NewScanner(in),
  }

  split := func(data []byte, atEOF bool) (advance int, token []byte, err error) {
  	// Skip leading spaces.
     start := 0
     for width := 0; start < len(data); start += width {
     	var r rune
     	r, width = utf8.DecodeRune(data[start:])
     	if !IsWhitespace(r) {
     		break
     	}
     }
     if atEOF && len(data) == 0 {
     	return 0, nil, nil
     }

     var r rune
     var width int
     r, width = utf8.DecodeRune(data[start:])
     if r == '(' || r == ')' {
       return start + width, data[start:start+width], nil
     }

     // Scan until space, marking end of word.
     for width, i := 0, start; i < len(data); i += width {
     	r, width = utf8.DecodeRune(data[start:])
     	if IsWhitespace(r) || r == '(' || r == ')' {
     		return i + width, data[start:i], nil
     	}
     }
     // If we're at EOF, we have a final, non-empty, non-terminated word. Return it.
     if atEOF && len(data) > start {
     	return len(data), data[start:], nil
     }
     // Request more data.
     return 0, nil, nil
  }
  input.Split(split)

  return input
}

func (in *Input) NextToken() (string, error) {
  if ok := in.Scan(); ok {
    return in.Text(), nil
  } else {
    return "", in.Err()
  }
}