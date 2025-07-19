package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"unicode"
)

var (
	flagA = flag.Bool("A", false, "equivalent to -vET")
	flagb = flag.Bool("b", false, "number nonempty output lines")
	flage = flag.Bool("e", false, "equivalent to -vE")
	flagE = flag.Bool("E", false, "display $ at end of each line")
	flagn = flag.Bool("n", false, "number all output lines")
	flags = flag.Bool("s", false, "suppress repeated empty output lines")
	flagT = flag.Bool("T", false, "display TAB characters as ^I")
	flagv = flag.Bool("v", false, "use ^ and M- notation, except for LFD and TAB")
)

func main() {
	flag.Parse()
	files := flag.Args()

	// Combine composite flags
	if *flagA {
		*flagv, *flagE, *flagT = true, true, true
	}
	if *flage {
		*flagv, *flagE = true, true
	}

	if len(files) == 0 {
		err := catFile(os.Stdin)
		if err != nil {
			fmt.Fprintln(os.Stderr, "cat:", err)
			os.Exit(1)
		}
		return
	}

	for _, fname := range files {
		var f *os.File
		var err error
		if fname == "-" {
			f = os.Stdin
		} else {
			f, err = os.Open(fname)
			if err != nil {
				fmt.Fprintf(os.Stderr, "cat: %s: %v\n", fname, err)
				continue
			}
			defer f.Close()
		}
		if err := catFile(f); err != nil {
			fmt.Fprintf(os.Stderr, "cat: %s: %v\n", fname, err)
		}
	}
}

func catFile(r io.Reader) error {
	scanner := bufio.NewScanner(r)
	var lineNum int
	var prevBlank bool

	for scanner.Scan() {
		line := scanner.Text()
		isBlank := len(line) == 0

		if *flags && isBlank && prevBlank {
			continue // squeeze blank lines
		}

		shouldNumber := false
		if *flagn {
			shouldNumber = true
		} else if *flagb && !isBlank {
			shouldNumber = true
		}

		if shouldNumber {
			lineNum++
			fmt.Printf("%6d\t", lineNum)
		}

		fmt.Println(transformLine(line))

		prevBlank = isBlank
	}
	return scanner.Err()
}

func transformLine(line string) string {
	var result string
	for _, r := range line {
		switch {
		case r == '\t' && *flagT:
			result += "^I"
		case *flagv && !unicode.IsPrint(r) && r != '\t':
			result += toCaretNotation(r)
		default:
			result += string(r)
		}
	}
	if *flagE {
		result += "$"
	}
	return result
}

func toCaretNotation(r rune) string {
	if r < 32 {
		return "^" + string(r+'@')
	}
	if r == 127 {
		return "^?"
	}
	if r > 127 {
		return "M-" + toCaretNotation(r-128)
	}
	return string(r)
}
