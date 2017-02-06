package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/unixpickle/bnf"
	"github.com/unixpickle/essentials"
)

func main() {
	if len(os.Args) != 2 && len(os.Args) != 3 {
		essentials.Die("Usage: bnf-sample <grammar.txt> [num]")
	}
	gFile := os.Args[1]
	num := 1
	if len(os.Args) == 3 {
		var err error
		num, err = strconv.Atoi(os.Args[2])
		if err != nil {
			essentials.Die("invalid sample count:", os.Args[2])
		}
	}

	f, err := os.Open(gFile)
	if err != nil {
		essentials.Die(err)
	}
	defer f.Close()

	g, err := bnf.ReadGrammar(f)
	if err != nil {
		essentials.Die(err)
	}
	if len(g) == 0 {
		essentials.Die("grammar has no rules")
	}

	rootName := g[0].Name
	for i := 0; i < num; i++ {
		out, err := bnf.Sample(g, rootName)
		if err != nil {
			essentials.Die(err)
		}
		fmt.Println(out)
	}
}
