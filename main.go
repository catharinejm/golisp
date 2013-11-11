package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

type Form interface{}
type Number interface{}

type Pair struct {
	head Form
	tail Form
}

func analyzeToken(token string) (Form) {
	r, _ := utf8.DecodeRuneInString(token)
	switch {
	case unicode.IsNumber(r) || r == '-':
		return readNumber(token)
	default:
		panic(fmt.Sprint("What is this? I can't deal with this. Stop giving me crap:", token))
	}
}

func readNumber(token string) (Number) {
	var n Number
	var err error
	if strings.Contains(token, ".") {
		n, err = strconv.ParseFloat(token, 64)
	} else {
		n, err = strconv.ParseInt(token, 0, 64)
	}

	if err != nil {
		panic(err)
	}
	return n
}

func readList(in *Input) (*Pair) {
	cur := in.NextToken()
	if cur == ")" {
		return nil
	}

	in.ReplaceToken(cur)
	head := readForm(in)

	var tail Form
	cur = in.NextToken()

	if cur == "." {
		tail = readForm(in)

		cur = in.NextToken()
		if cur != ")" {
			panic("Invalid list structure.")
		}
	} else {
		in.ReplaceToken(cur)
		tail = readList(in)
	}

	return &Pair{head, tail}
}

func readForm(in *Input) (Form) {
	token := in.NextToken()
	if token == "(" {
		return readList(in)
	} else {
		return analyzeToken(token)
	}

	panic("Something weird happened.")
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

func repl(rdr *bufio.Reader) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("Error:", err)
			repl(rdr)
		}
	}()
	in := NewInput(rdr)
	for {
		var f Form

		fmt.Print("> ")
		tok := in.NextToken()
		if tok == "" {
			os.Exit(0)
		}
		in.ReplaceToken(tok)

		f = readForm(in)

		fmt.Print("VALUE: ")
		printForm(f)
		fmt.Println()
	}
}

func main() {
	var rdr *bufio.Reader
	rdr = bufio.NewReader(os.Stdin)
	repl(rdr)
}
