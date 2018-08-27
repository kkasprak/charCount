package main

import (
	"strings"
	"fmt"
	"io/ioutil"
	"os"
	"flag"
)

const proverbs = `Don't communicate by sharing memory, share memory by communicating.
Concurrency is not parallelism.
Channels orchestrate; mutexes serialize.
The bigger the interface, the weaker the abstraction.
Make the zero value useful.
interface{} says nothing.
Gofmt's style is no one's favorite, yet gofmt is everyone's favorite.
A little copying is better than a little dependency.
Syscall must always be guarded with build tags.
Cgo must always be guarded with build tags.
Cgo is not Go.
With the unsafe package there are no guarantees.
Clear is better than clever.
Reflection is never clear.
Errors are values.
Don't just check errors, handle them gracefully.
Design the architecture, name the components, document the details.
Documentation is for users.
Don't panic.`

type proverb struct {
	line string
	chars map[rune]int
}


func main() {
	/*Get file using flags.
	To use:
	go run main.go -f proverbs.txt
	or
	FILE = proverbs.txt go run main.go*/

	path := pathFromFlag()
	if path == "" {
		path = pathFromEnv()
	}
	if path == "" {
		fmt.Println("You must specify one of the file with with -f or as FILE environment var.")
		os.Exit(1)
	}

	proverbs, err := loadProverbs(path)
	if err != nil {
		fmt.Errorf("error loading proverbs")
	}


	//channels synchronize access to shared resources
	ch := make(chan *proverb)
	go printProverbs(ch)
	//range over all the proverbs
	for _, p := range proverbs {
		//send the value to channel
		ch <- p
	}
	//nothing else to be sent, so close the channel
	close(ch)
}

func printProverbs(pc chan *proverb) {
	//range over the channel
	for p := range pc {
		fmt.Printf("%s\n", p.line)
		for k, v := range p.newCount() {
			fmt.Printf("'%c'=%d, ", k, v)
		}
		fmt.Print("\n\n")
	}
}

func loadProverbs(path string) ([]*proverb, error) {
	var proverbs []* proverb
	bs, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(bs), "\n")
	for _, line := range lines {
		// for each one, create a proverb and set proverb.line
		p := &proverb{line: line}
		proverbs = append(proverbs, p)
	}
	return proverbs, nil
}

// Count takes a line (string) and returns a map of characters and their counts
// Rune can be a multi byte character that encompasses multiple character sets (eg different langs)
func count(line string) map[rune]int {
	//same as:
	//var m map[rune]int
	//but specifying a size, the map will grow automatically
	m := make(map[rune]int, 0)
	for _, c := range line {
		m[c] = m[c] + 1
	}
	return m
}

//this count function has a receiver
func (p proverb) newCount() map[rune]int {
	if p.chars != nil {
		return p.chars
	}
	m := make(map[rune]int, 0)
	for _, c := range p.line {
		m[c] = m[c] + 1
	}
	p.chars = m
	return p.chars
}

func pathFromFlag() string {
	//value: is the default value in case -f isnt specified
	path := flag.String("f", "", "file flag")
	//Each flag will be pointed to back in memory with the Parse function
	flag.Parse()
	return *path
}

func pathFromEnv() string {
	return os.Getenv("FILE")
}




