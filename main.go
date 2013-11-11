package main

import (
	"bufio"
	"fmt"
	"os"
	"unicode"
	"strings"
	"strconv"
)

type Form interface{}
type Number interface{}

type Pair struct {
	head Form
	tail Form
}

func readNumber(in Input) (Number, error) {
  var num string
	_, err := fmt.Fscan(in, &num)
	if err != nil {
		return 0, err
	}
  if strings.HasSuffix(num, ")") {
    in.BacktrackBytes(1)
    num = strings.TrimSuffix(num, ")")
  }

  if strings.Contains(num, ".") {
    return strconv.ParseFloat(num, 64)
  } else {
    return strconv.ParseInt(num, 0, 64)
  }

  return 0, fmt.Errorf("What kind of number is this: %s", num)
}

func readList(in Input) (*Pair, error) {
	cur := in.GetRune()

	if cur == ')' {
		return nil, nil
	}

	in.Backtrack()

	head, err := readForm(in)
	if err != nil {
		return nil, err
	}

	cur = in.GetRune()
	var tail Form

	if cur == '.' {
		tail, err = readForm(in)
		if err != nil {
			return nil, err
		}

		cur = in.GetRune()
		if cur != ')' {
			return nil, fmt.Errorf("Invalid list structure.")
		}
	} else {
		in.Backtrack()
		tail, err = readList(in)
		if err != nil {
			return nil, err
		}
	}

	return &Pair{head, tail}, nil
}

func readForm(in Input) (Form, error) {
	cur := in.GetRune()

	switch {
	case cur == '(':
		return readList(in)
	case unicode.IsNumber(cur):
		in.Backtrack()
		return readNumber(in)
	}

	return nil, fmt.Errorf("Something weird happened.")
}

func printList(list *Pair) {
	if list == nil {
		return
	}

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
	case *Pair:
		fmt.Print("(")
		printList(form.(*Pair))
		fmt.Print(")")
	case int64, float64:
		fmt.Print(form)
	}
}

func main() {
	rdr := Input{bufio.NewReader(os.Stdin)}
	for {
		var f Form
		var err error

		fmt.Print("> ")
		f, err = readForm(rdr)
		if err != nil {
			rdr.ReadLine()
			fmt.Println("Error:", err)
		} else {
			printForm(f)
			fmt.Println()
		}
	}
}
