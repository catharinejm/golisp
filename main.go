package main

/*
#cgo CFLAGS: -I/usr/local/Cellar/readline/6.2.4/include
#cgo CFLAGS: -Qunused-arguments
#cgo LDFLAGS: -L/usr/local/Cellar/readline/6.2.4/lib
#cgo LDFLAGS: -lreadline
#include <readline/readline.h>
*/

import (
	"C"
	"unsafe"
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

func analyzeToken(token string) (Form, error) {
	r, _ := utf8.DecodeRuneInString(token)
	switch {
	case unicode.IsNumber(r) || r == '-':
		return readNumber(token)
	default:
		return nil, fmt.Errorf("What is this? I can't deal with this. Stop giving me crap: %s", token)
	}
}

func readNumber(token string) (Number, error) {
	if strings.Contains(token, ".") {
		return strconv.ParseFloat(token, 64)
	} else {
		return strconv.ParseInt(token, 0, 64)
	}

	return 0, fmt.Errorf("What kind of number is this: %s", token)
}

func readList(in *Input) (*Pair, error) {
	cur, err := in.NextToken()
	if err != nil {
		return nil, err
	}
	if cur == ")" {
		return nil, nil
	}

	in.ReplaceToken(cur)
	head, err := readForm(in)
	if err != nil {
		return nil, err
	}

	var tail Form
	cur, err = in.NextToken()
	if err != nil {
		return nil, err
	}

	if cur == "." {
		tail, err = readForm(in)
		if err != nil {
			return nil, err
		}

		cur, err = in.NextToken()
		if err != nil {
			return nil, err
		}
		if cur != ")" {
			return nil, fmt.Errorf("Invalid list structure.")
		}
	} else {
		in.ReplaceToken(cur)
		tail, err = readList(in)
		if err != nil {
			return nil, err
		}
	}

	return &Pair{head, tail}, nil
}

func readForm(in *Input) (Form, error) {
	token, err := in.NextToken()
	if err != nil {
		return nil, err
	}
	if token == "(" {
		return readList(in)
	} else {
		return analyzeToken(token)
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
	prompt := C.CString("> ")
	for {
		var f Form
		var err error

		line := C.readline(prompt)
		if line == nil {
			os.Exit(0)
		}

		scanner := NewInput(strings.Reader(C.GoString(line)))
		C.free(unsafe.Pointer(line))

		f, err = readForm(scanner)
		if err != nil {
			fmt.Println("Error:", err)
		} else {
			fmt.Print("VALUE: ")
			printForm(f)
			fmt.Println()
		}
	}
}
