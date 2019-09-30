package main

import (
	"bufio"
	flag "github.com/spf13/pflag"
	"io"
	"os"
	"os/exec"
)

func main() {
	var start, end, linesPerPage int
	var pageDelimiter bool
	var destination string

	flag.IntVarP(&start, "s", "s", 0, "start page No.")
	flag.IntVarP(&end, "e", "e", 0, "end page No.")
	flag.IntVarP(&linesPerPage, "l", "l", 72,"lines per page")
	flag.BoolVarP(&pageDelimiter, "f", "f", false, "use \\f as page delimiter")
	flag.StringVarP(&destination, "d", "d", "", "destination destination name")

	flag.Parse()

	if flag.NFlag() < 2 {
		flag.Usage()
		return
	}

	if start <= 0 || end <= 0 {
		panic("page No. should be positive integer")
	}

	if start > end {
		panic("end page No. should not be greater than start page No.")
	}

	if linesPerPage <= 0 {
		panic("there should be positive lines per page")
	}

	scanner := bufio.NewScanner(os.Stdin)

	if flag.NArg() > 0 {
		file, err := os.Open(flag.Arg(0))
		if err != nil {
			panic(err)
		}

		scanner = bufio.NewScanner(file)
	}

	writer := bufio.NewWriter(os.Stdout)

	if destination != "" {
		cmd := exec.Command("lp", "-d" + destination)
		pipeReader, pipeWriter := io.Pipe()
		cmd.Stdin = pipeReader
		writer = bufio.NewWriter(pipeWriter)
	}

	page, line := 1, 1

	if pageDelimiter {
		pageSplit := func(data []byte, atEOF bool) (advance int, token []byte, err error) {
			for i := 0; i < len(data); i++ {
				if data[i] == '\f' {
					return i + 1, data[:i], nil
				}
			}
			if !atEOF {
				return 0, nil, nil
			}
			return 0, data, bufio.ErrFinalToken
		}

		scanner.Split(pageSplit)
	}

	for scanner.Scan() {
		if page >= start && page <= end {
			if pageDelimiter {
				_, err := writer.WriteString(scanner.Text() + "\f")
				if err != nil {
					panic(err)
				}
			} else {
				_, err := writer.WriteString(scanner.Text() + "\n")
				if err != nil {
					panic(err)
				}
			}
			_ = writer.Flush()
		}

		if pageDelimiter {
			page++
		} else {
			line++

			if line > linesPerPage {
				line = 1
				page++
				if page > end {
					break
				}
			}
		}
	}

	return
}