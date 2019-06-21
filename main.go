package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/lhopki01/lexer-experiment/ast"
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

	jenkinsFile := p.ParseJenkinsFile()
	fmt.Println("==============")
	fmt.Println(jenkinsFile.Library)
	fmt.Println(jenkinsFile.Function)
	fmt.Println(jenkinsFile.Imports)
	spew.Dump(jenkinsFile)

	js, _ := json.MarshalIndent(jenkinsFile, "", "    ")
	fmt.Println(string(js))
	fmt.Println(reflect.TypeOf(jenkinsFile.Library))
	fmt.Println(reflect.TypeOf(jenkinsFile.Values))
	makeTargets := jenkinsFile.Values["makeTargets"]
	fmt.Println(reflect.TypeOf(makeTargets))
	for _, target := range makeTargets.([]interface{}) {
		t := strings.TrimSuffix(strings.TrimPrefix(target.(string), "\""), "\"")
		fmt.Printf("%s ==> make %s\n", t, t)

	}

	fmt.Println(jenkinsFile.Library)
	for _, i := range jenkinsFile.Imports {
		fmt.Printf("import %s\n", i)
	}
	fmt.Println("")
	fmt.Printf("%s ", jenkinsFile.Function)
	printBody("  ", "{", "}", "=", jenkinsFile.Values)
	//for target := range jenkinsFile.Values["makeTargets"] {
	//	spew.Dump(target)
	//}
	//makeTargets := jenkinsFile.Values["makeTargets"].([]string)
	//for target := range makeTargets {
	//	fmt.Printf("%s ==> make %s", target, target)
	//}
}

func printBody(indent string, lbracket string, rbracket string, assignment string, object interface{}) {
	switch vt := object.(type) {
	case ast.ConcatenatedItem:
		//fmt.Println("========")
		//fmt.Println(reflect.TypeOf(vt))
		//fmt.Println("========")
		fmt.Printf("%v << ", vt.Primary)
		printBody("  "+indent, "[", "]", ":", vt.Append)
	case []interface{}:
		fmt.Println(lbracket)
		for _, s := range vt {
			fmt.Printf("%s%v,\n", indent, s)
		}
		fmt.Printf("%s%s\n", strings.TrimPrefix(indent, "  "), rbracket)
	case map[string]interface{}:
		fmt.Println(lbracket)
		for k, v := range vt {
			fmt.Printf("%s%s %s ", indent, k, assignment)
			printBody("  "+indent, "[", "]", ":", v)
		}
		fmt.Printf("%s%s\n", strings.TrimPrefix(indent, "  "), rbracket)
	case string:
		fmt.Printf("%s\n", vt)
	default:
		fmt.Println(reflect.TypeOf(vt))
	}
	//fmt.Printf("  %s = %v\n", key, value)
}
