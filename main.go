package main

import (
	"bufio"
	"fmt"
	"os"
	"unicode"
	"unicode/utf8"
	"strings"
	"strconv"
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
    return nil, fmt.Errorf("What is this? I can't deal with this. Stop giving me crap. (%s)", token)
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

func readList(in *Input, cur string) (*Pair, error) {
	head, err := readForm(in, cur)
	if err != nil {
		return nil, err
	}

	cur, err = in.NextToken()
  if err != nil { return nil, err }

	var tail Form

  if cur == ")" {
    tail = nil
  } else if cur == "." {
    cur, err = in.NextToken()
		if err != nil { return nil, err }

		tail, err = readForm(in, cur)
		if err != nil { return nil, err }

		cur, err = in.NextToken()
		if err != nil { return nil, err }
		if cur != ")" {
			return nil, fmt.Errorf("Invalid list structure.")
		}
	} else {
		tail, err = readList(in, cur)
		if err != nil {
			return nil, err
		}
	}

	return &Pair{head, tail}, nil
}

func readForm(in *Input, token string) (Form, error) {
	switch {
	case token == "(":
    token, err := in.NextToken()
    if err != nil { return nil, err }
    if token == ")" { return nil, nil }
		return readList(in, token)
  default:
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
	scanner := NewInput(bufio.NewReader(os.Stdin))
	for {
		//var f Form
		var err error

		fmt.Print("> ")
    tok, err := scanner.NextToken()
    if err != nil {
      fmt.Println("Error:", err)
      continue
    }
    fmt.Printf("read: %s\n", tok)

//		f, err = readForm(scanner, tok)
//		if err != nil {
//      // rdr.ReadLine()
//			fmt.Println("Error:", err)
//		} else {
//			printForm(f)
//			fmt.Println()
//		}
	}
}
