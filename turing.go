package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var help string = `usage: turing [-h | -t]

Turing accepts a definition of a Turing machine from standard input and 
simulates its execution.

The -t flag accepts an extended definition that includes a tape and 
initial position for the head.

The format of the definition follows:

	head position    (if -t is specified)
	tape             (if -t is specified)
	blank character
	initial state
	final state
	transition rules

The format of a transition rule is:

	current_state input next_state output move`

func main() {
	if len(os.Args) > 1 && os.Args[1] == "-h" {
		fmt.Println(help)
		os.Exit(0)
	}

	bs, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	
	s := string(bs)
	lines := getLines(s)

	var (
		program Program
		head    int
		tape    []string
	)

	if len(os.Args) > 1 && os.Args[1] == "-t" {
		program, head, tape = parse(true, lines)
	} else {
		program, _, _ = parse(false, lines)
		head = 0
		tape = []string{program.blank}
	}

	exec(program, head, &tape)
	fmt.Println(strings.Join(tape, " "))
}

type Program struct {
	blank      string
	initial    string
	final      string
	transition map[[2]string][3]string // (state, input) -> (state, output, move)
}

var newlineRegexp = regexp.MustCompile(`\r?\n`)

func getLines(input string) []string {
	trimmed := strings.TrimSpace(input)
	lines := newlineRegexp.Split(trimmed, -1)
	return lines
}

func parse(hasTape bool, lines []string) (Program, int, []string) {
	var head int
	var tape []string
	if hasTape {
		head, _ = strconv.Atoi(lines[0])
		tape = strings.Split(lines[1], " ")
		lines = lines[2:]
	}
	blank := lines[0]
	initial := lines[1]
	final := lines[2]
	transition := make(map[[2]string][3]string)
	for i := 3; i < len(lines); i++ {
		line := lines[i]
		symbols := strings.Split(line, " ")
		transition[[2]string{symbols[0], symbols[1]}] = [3]string{symbols[2], symbols[3], symbols[4]}
	}
	return Program{blank, initial, final, transition}, head, tape
}

func exec(program Program, head int, pTape *[]string) {
	currentState := program.initial
	for currentState != program.final {
		if head == -1 {
			newTape := []string{program.blank}
			*pTape = append(newTape, (*pTape)...)
			head = 0
		} else if head == len(*pTape) {
			*pTape = append(*pTape, program.blank)
		}

		input := (*pTape)[head]

		t := program.transition[[2]string{currentState, input}]
		currentState = t[0]
		(*pTape)[head] = t[1]
		if t[2] == "l" {
			head--
		} else {
			head++
		}
	}
}
