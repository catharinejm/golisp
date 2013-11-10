package main

import (
	"bufio"
	"fmt"
	"os"
	"unicode"
)

type Form interface{}
type Number int64

type Pair struct {
	head Form
	tail Form
}

func readRune(input *bufio.Reader) rune {
	cur, _, err := input.ReadRune()
	if err != nil {
		os.Exit(0)
	}
	return cur
}

func killWhitespace(input *bufio.Reader) {
	for {
		cur := readRune(input)
		if !unicode.IsSpace(cur) {
			input.UnreadRune()
			break
		}
	}
}

func nextBlackspaceRune(input *bufio.Reader) rune {
  killWhitespace(input)
  return readRune(input)
}

func flushInput(input *bufio.Reader) {
  for _, err := input.ReadByte(); err != nil; _, err = input.ReadByte() {}
}

func readNumber(input *bufio.Reader) Number {
	var n Number
	fmt.Fscanf(input, "%d", &n)
	return n
}

func readList(input *bufio.Reader) *Pair {
	cur := nextBlackspaceRune(input)

	if cur == ')' {
		return nil
	}

  input.UnreadRune()

	head := readForm(input)
	cur = nextBlackspaceRune(input)
	var tail Form
	if cur == '.' {
		tail = readForm(input)
    cur = nextBlackspaceRune(input)
    if cur != ')' {
      fmt.Println("ERROR: Invalid list structure.")
      flushInput(input)
      return nil
    }
	} else {
		input.UnreadRune()
		tail = readList(input)
	}

	return &Pair{head, tail}
}

func readForm(input *bufio.Reader) Form {
	cur := nextBlackspaceRune(input)

	switch {
	case cur == '(':
		return readList(input)
	case unicode.IsNumber(cur):
		input.UnreadRune()
		return readNumber(input)
	}
	return nil
}

func printList(list *Pair) {
  if list == nil { return }

  printForm(list.head)

  switch list.tail.(type) {
  case *Pair:
    tl := list.tail.(*Pair)
    if tl != nil {
      fmt.Print(" ")
      printList(tl)
    }
  default:
    fmt.Print(" . ")
    printForm(list.tail)
  }
}

func printForm(form Form) {
  switch form.(type) {
  case Number:
    fmt.Print(form)
  case *Pair:
    fmt.Print("(")
    printList(form.(*Pair))
    fmt.Print(")")
  }
}

func main() {
	rdr := bufio.NewReader(os.Stdin)
	for {
		var f Form

		fmt.Print("> ")
		f = readForm(rdr)
		printForm(f)
    fmt.Println()
	}
}
