package main

import (
	"bufio"
	"os"
	"unicode"
)

func IsWhitespace(r rune) bool {
  return unicode.IsSpace(r) || r == ','
}

type Input struct {
	*bufio.Reader
}

func (in Input) NextRune() rune {
	cur, _, err := in.ReadRune()
	if err != nil {
		os.Exit(0)
	}
	return cur
}

func (in *Input) GetRune() rune {
	in.StripWhitespace()
	return in.NextRune()
}

func (in *Input) Backtrack() {
	err := in.UnreadRune()
	if err != nil {
		panic(err)
	}
}

func (in *Input) BacktrackBytes(n int) {
  for i := 0; i < n; i++ {
    err := in.UnreadByte()
    if err != nil {
      panic(err)
    }
  }
}

func (in *Input) StripWhitespace() {
	for {
		cur := in.NextRune()
		if !IsWhitespace(cur) {
			in.Backtrack()
			break
		}
	}
}
