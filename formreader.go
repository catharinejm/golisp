package main

import (
	//"fmt"
	"bufio"
	"io"
	"unicode"
	"unicode/utf8"
)

func IsWhitespace(r rune) bool {
	return unicode.IsSpace(r) || r == ','
}

type Input struct {
	lastTok string
	*bufio.Scanner
}

func NewInput(in io.Reader) *Input {
	input := &Input{
		"",
		bufio.NewScanner(in),
	}

	split := func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		//fmt.Println("Input:", string(data))
		// Skip leading spaces.
		start := 0
		for width := 0; start < len(data); start += width {
			var r rune
			r, width = utf8.DecodeRune(data[start:])
			if !IsWhitespace(r) {
				break
			}
		}
		if atEOF && len(data[start:]) == 0 {
			//fmt.Println("need more data 1")
			return 0, nil, nil
		}

		//fmt.Println("After WS Skip:", string(data[start:]))

		var r rune
		var width int
		r, width = utf8.DecodeRune(data[start:])
		if r == '(' || r == ')' {
			//fmt.Println("returning token:", string(data[start:start+width]))
			return start + width, data[start : start+width], nil
		}
		//fmt.Println("After paren check:", string(data[start:]))

		// Scan until space, marking end of word.
		for width, i := 0, start; i < len(data); i += width {
			r, width = utf8.DecodeRune(data[i:])
			//fmt.Printf("rune %d: %s\n", i, string(r))
			if IsWhitespace(r) || r == '(' || r == ')' {
				//fmt.Println("returning token:", string(data[start:i]))
				return i, data[start:i], nil
			}
		}
		// If we're at EOF, we have a final, non-empty, non-terminated word. Return it.
		if atEOF && len(data) > start {
			//fmt.Println("returning token:", string(data[start:]))
			return len(data) - start, data[start:], nil
		}
		// Request more data.
		//fmt.Println("need more data 2")
		return 0, nil, nil
	}
	input.Split(split)

	return input
}

func (in *Input) NextToken() (string) {
	if !in.ReadNextToken() {
		panic(in.Err())
	}
	if in.lastTok == "" {
		return in.Text()
	} else {
		return in.lastTok
	}
}

func (in *Input) ReplaceToken(tok string) {
	in.lastTok = tok
}

func (in *Input) ReadNextToken() bool {
	if in.lastTok != "" {
		in.lastTok = ""
		return true
	}
	return in.Scan()
}
