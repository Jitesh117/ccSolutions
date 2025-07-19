package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"unicode/utf8"
)

type counts struct {
	lines   int
	words   int
	bytes   int
	chars   int
	maxLine int
}

func countReader(r io.Reader) counts {
	var c counts
	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		line := scanner.Text()
		c.lines++
		c.bytes += len(line) + 1
		c.words += len(strings.Fields(line))
		c.chars += utf8.RuneCountInString(line) + 1

		if len(line) > c.maxLine {
			c.maxLine = len(line)
		}
	}

	return c
}

func printCounts(
	c counts,
	showLines, showWords, showBytes, showChars, showMaxLine bool,
	label string,
) {
	if !showLines && !showWords && !showBytes && !showChars && !showMaxLine {
		// Default: all except chars and maxLine
		showLines, showWords, showBytes = true, true, true
	}

	if showLines {
		fmt.Printf("%7d", c.lines)
	}
	if showWords {
		fmt.Printf("%7d", c.words)
	}
	if showBytes {
		fmt.Printf("%7d", c.bytes)
	}
	if showChars {
		fmt.Printf("%7d", c.chars)
	}
	if showMaxLine {
		fmt.Printf("%7d", c.maxLine)
	}
	if label != "" {
		fmt.Printf(" %s", label)
	}
	fmt.Println()
}

func main() {
	// Flags
	showLines := flag.Bool("l", false, "print the newline counts")
	showWords := flag.Bool("w", false, "print the word counts")
	showBytes := flag.Bool("c", false, "print the byte counts")
	showChars := flag.Bool("m", false, "print the character counts")
	showMaxLine := flag.Bool("L", false, "print the length of the longest line")
	flag.Parse()

	files := flag.Args()
	var total counts
	var fileCount int

	if len(files) == 0 {
		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeCharDevice) != 0 {
			fmt.Fprintln(os.Stderr, "Reading from stdin. Press Ctrl+D to end.")
		}
		c := countReader(os.Stdin)
		printCounts(c, *showLines, *showWords, *showBytes, *showChars, *showMaxLine, "")
		return
	}

	for _, file := range files {
		f, err := os.Open(file)
		if err != nil {
			fmt.Fprintf(os.Stderr, "wc: %s: %v\n", file, err)
			continue
		}
		defer f.Close()

		c := countReader(f)
		printCounts(c, *showLines, *showWords, *showBytes, *showChars, *showMaxLine, file)

		total.lines += c.lines
		total.words += c.words
		total.bytes += c.bytes
		total.chars += c.chars
		if c.maxLine > total.maxLine {
			total.maxLine = c.maxLine
		}
		fileCount++
	}

	if fileCount > 1 {
		printCounts(total, *showLines, *showWords, *showBytes, *showChars, *showMaxLine, "total")
	}
}
