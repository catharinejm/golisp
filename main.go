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

func readNumber(input *bufio.Reader) (Number, error) {
	var n Number
  _, err := fmt.Fscanf(input, "%d", &n)
  if err != nil {
    return 0, err
  }
  next := readRune(input)
  if !unicode.IsSpace(next) && next != ')' {
    return 0, fmt.Errorf("Invalid number")
  }
  input.UnreadRune()
	return n, nil
}

func readList(input *bufio.Reader) (*Pair, error) {
	cur := nextBlackspaceRune(input)

	if cur == ')' {
		return nil, nil
	}

  input.UnreadRune()

	head, err := readForm(input)
  if err != nil {
    return nil, err
  }

	cur = nextBlackspaceRune(input)
	var tail Form

	if cur == '.' {
		tail, err = readForm(input)
    if err != nil {
      return nil, err
    }

    cur = nextBlackspaceRune(input)
    if cur != ')' {
      return nil, fmt.Errorf("Invalid list structure.")
    }
	} else {
		input.UnreadRune()
		tail, err = readList(input)
    if err != nil {
      return nil, err
    }
	}

	return &Pair{head, tail}, nil
}

func readForm(input *bufio.Reader) (Form, error) {
	cur := nextBlackspaceRune(input)

	switch {
	case cur == '(':
		return readList(input)
	case unicode.IsNumber(cur):
		input.UnreadRune()
		return readNumber(input)
	}

	return nil, fmt.Errorf("Something weird happened.")
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
    var err error

		fmt.Print("> ")
		f, err = readForm(rdr)
    if (err != nil) {
      rdr.ReadLine()
      fmt.Println("Error:", err);
    } else {
      printForm(f)
      fmt.Println()
    }
	}
}
