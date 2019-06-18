package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/lhopki01/lexer-experiment/lexer"
	"github.com/lhopki01/lexer-experiment/parser"
)

func main() {
	if len(os.Args) < 2 {
		panic("no valid file name or path provided for file!")
	}

	path := os.Args[1]
	absPath, _ := filepath.Abs(path)
	data, err := ioutil.ReadFile(absPath)
	if err != nil {
		panic(err.Error())
	}

	l := lexer.NewLexer(data)

	p := parser.NewParser(l)

	ast := p.Parse()
	for s := range ast {
		fmt.Println(s)
	}
	js, _ := json.MarshalIndent(ast, "", "    ")
	fmt.Println(string(js))
}
