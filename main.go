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

func readNumber(in Input) (Number, error) {
	var n Number
	_, err := fmt.Fscanf(in, "%d", &n)
	if err != nil {
		return 0, err
	}
	next := in.NextRune()
	if !unicode.IsSpace(next) && next != ')' {
		return 0, fmt.Errorf("Invalid number")
	}
	in.Backtrack()
	return n, nil
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
	case Number:
		fmt.Print(form)
	case *Pair:
		fmt.Print("(")
		printList(form.(*Pair))
		fmt.Print(")")
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
