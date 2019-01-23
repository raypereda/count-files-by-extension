// Cammand extdir list the counts of files by extension.
package main

import (
	"flag"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

var fileCount int
var extCount = make(map[string]int)

func walk(root string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}
	ext := filepath.Ext(root)
	ext = strings.ToLower(ext)
	extCount[ext]++
	fileCount++
	return nil
}

var done = make(chan bool)
var program string
var version = "0.1"

var flagV = flag.Bool("version", false, "Print version and exit")

func main() {
	program = path.Base(os.Args[0])
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr,
			"%s is a command that recurses a directory and reports file counts by extension.\n",
			program)
		fmt.Fprintf(os.Stderr, "Version %s\n", version)
		fmt.Fprintf(os.Stderr, "Usage: %s PATH\n", program)
		flag.PrintDefaults()
	}
	flag.Parse()

	if *flagV {
		fmt.Printf("%s version %s\n", program, version)
		return
	}
	args := flag.Args()

	var root string
	if len(args) == 1 {
		root = args[0]
	} else {
		flag.Usage()
		return
	}

	go markProgress()
	filepath.Walk(root, walk)
	done <- true

	printExtCount(extCount)
}

func markProgress() {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-done:
			fmt.Fprintf(os.Stderr, "Done!\n")
			fmt.Fprintf(os.Stderr, "File count: %d\n", fileCount)
			return
		case <-ticker.C:
			fmt.Fprintf(os.Stderr, "Files count: %d\n", fileCount)
		}
	}
}

type pair struct {
	Key   string
	Value int
}

type pairList []pair

func (p pairList) Len() int           { return len(p) }
func (p pairList) Less(i, j int) bool { return p[i].Value < p[j].Value }
func (p pairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

func printExtCount(counts map[string]int) {
	ranked := rankByExtCount(counts)
	fmt.Printf("%10s %s\n", "#", "extension")

	for _, pair := range ranked {
		fmt.Printf("%4d %s\n", pair.Value, pair.Key)
	}
}

func rankByExtCount(extFrequencies map[string]int) pairList {
	pl := make(pairList, len(extFrequencies))
	i := 0
	for k, v := range extFrequencies {
		pl[i] = pair{k, v}
		i++
	}
	sort.Sort(sort.Reverse(pl))
	return pl
}
