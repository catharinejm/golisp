package main

import (
	"bufio"
	"unicode"
)

type Input struct {
	*bufio.Reader
}

func (in Input) NextRune() rune {
	cur, _, err := in.ReadRune()
	if err != nil {
		panic(err)
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

func (in *Input) StripWhitespace() {
	for {
		cur := in.NextRune()
		if !unicode.IsSpace(cur) {
			in.Backtrack()
			break
		}
	}
}
