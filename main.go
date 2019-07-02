package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"sort"
	"strings"

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

	if jenkinsFile.Function == "CICD" {
		convertContainerImages(&jenkinsFile)

		convertMakeTargets(&jenkinsFile)
		convertNpmRunTargets(&jenkinsFile)
		convertRakeTargets(&jenkinsFile)

		convertMoveToAll(&jenkinsFile)

		convertStepConfig(&jenkinsFile, "pr", "Constants.PR_VALIDATION_STEPS")
		convertStepConfig(&jenkinsFile, "promoteToProduction", "Constants.PROMOTION_JOB_STEPS")
		convertStepConfig(&jenkinsFile, "master", "Constants.MASTER_BRANCH_STEPS")

		delete(jenkinsFile.Values["all"].(map[string]interface{}), "stepConfig")
	} else {
		convertContainerImages(&jenkinsFile)
	}

	fmt.Println(jenkinsFile.Library)
	for _, i := range jenkinsFile.Imports {
		fmt.Printf("import %s\n", i)
	}
	fmt.Println("")
	fmt.Printf("%s ", jenkinsFile.Function)
	reg := regexp.MustCompile(`\[\s*\]`)
	fmt.Print(reg.ReplaceAllString(
		printBody("  ", "{", "}", "=", jenkinsFile.Values),
		"[]",
	))
}

func convertStepConfig(js *ast.JenkinsFile, key string, constant string) {
	all := js.Values["all"].(map[string]interface{})
	allStepConfig, ok := all["stepConfig"]

	stepConfig := ast.ConcatenatedItem{
		Primary: constant,
	}
	v, exists := js.Values[key]
	if exists {
		if innerStepConfig, innerExists := v.(map[string]interface{})["stepConfig"]; innerExists {
			stepConfig.Append = innerStepConfig
		} else if ok {
			stepConfig.Append = allStepConfig
		}
		js.Values[key].(map[string]interface{})["stepConfig"] = stepConfig
	} else {
		if ok {
			stepConfig.Append = allStepConfig
		}
		js.Values[key] = map[string]interface{}{

			"stepConfig": stepConfig,
		}
	}
}

func convertContainerImages(jf *ast.JenkinsFile) {
	images := map[string]interface{}{}
	re := regexp.MustCompile(`containerImages\[(.*)\]`)
	for k, v := range jf.Values {
		matches := re.FindStringSubmatch(k)
		if matches != nil {
			//images[matches[1]] = v
			vm := v.(map[string]interface{})
			if val, ok := vm["uri"]; ok {
				_, nameOk := vm["name"]
				_, tagOk := vm["tag"]
				if nameOk || tagOk {
					panic("Can't have name or tag with uri")
				} else {
					images[matches[1]] = val
				}
			} else {
				images[matches[1]] = fmt.Sprintf(
					"eu.gcr.io/karhoo-common/%s:%s",
					stripQuotes(vm["name"].(string)),
					stripQuotes(vm["tag"].(string)),
				)
			}

			delete(jf.Values, k)
		}
	}
	if len(images) > 0 {
		jf.Values["containerImages"] = images
	}
}

func convertMoveToAll(jf *ast.JenkinsFile) {
	allValues := map[string]interface{}{}
	for k, v := range jf.Values {
		if k != "pr" && k != "master" && k != "promoteToProd" {
			allValues[k] = v
			delete(jf.Values, k)
		}
	}
	if len(allValues) > 0 {
		jf.Values["all"] = allValues
	}
}

func convertMakeTargets(jf *ast.JenkinsFile) {
	scriptTargets := []interface{}{}

	if val, ok := jf.Values["scriptTargets"]; ok {
		scriptTargets = append(scriptTargets, val.([]interface{})...)
	}

	if val, ok := jf.Values["makeTargets"]; ok {
		for _, target := range val.([]interface{}) {
			scriptTargets = append(scriptTargets, fmt.Sprintf(
				`'make %s'`,
				stripQuotes(target.(string)),
			))
		}
	}

	jf.Values["scriptTargets"] = scriptTargets

	delete(jf.Values, "makeTargets")
}

func convertRakeTargets(jf *ast.JenkinsFile) {
	scriptTargets := []interface{}{}

	if val, ok := jf.Values["scriptTargets"]; ok {
		scriptTargets = append(scriptTargets, val.([]interface{})...)
	}

	if val, ok := jf.Values["rakeTargets"]; ok {
		for _, target := range val.([]interface{}) {
			scriptTargets = append(scriptTargets, fmt.Sprintf(
				`'rake %s'`,
				stripQuotes(target.(string)),
			))
		}
	}

	jf.Values["scriptTargets"] = scriptTargets

	delete(jf.Values, "rakeTargets")
}

func convertNpmRunTargets(jf *ast.JenkinsFile) {
	scriptTargets := []interface{}{}

	if val, ok := jf.Values["scriptTargets"]; ok {
		scriptTargets = append(scriptTargets, val.([]interface{})...)
	}

	if val, ok := jf.Values["npmRunTargets"]; ok {
		for _, target := range val.([]interface{}) {
			scriptTargets = append(scriptTargets, fmt.Sprintf(
				`'npm run %s'`,
				stripQuotes(target.(string)),
			))
		}
	}

	jf.Values["scriptTargets"] = scriptTargets

	delete(jf.Values, "npmRunTargets")
}

func stripQuotes(s string) string {
	start := string([]rune(s)[0])
	if start == `"` {
		s = strings.TrimPrefix(s, `"`)
		s = strings.TrimSuffix(s, `"`)
		s = strings.Replace(s, `\"`, `"`, -1)
		s = strings.Replace(s, `'`, `\'`, -1)
	} else if start == `'` {
		s = strings.TrimPrefix(s, `'`)
		s = strings.TrimSuffix(s, `'`)
	}
	return s

}

func printBody(indent string, lbracket string, rbracket string, assignment string, object interface{}) string {
	switch vt := object.(type) {
	case ast.ConcatenatedItem:
		//fmt.Println("ConcatenatedItem")
		if vt.Append == nil {
			return fmt.Sprint(vt.Primary)
		} else {
			return fmt.Sprintf("%v << %s", vt.Primary, printBody(""+indent, "[", "]", ":", vt.Append))
		}
	case []interface{}:
		//fmt.Println("Slice")
		s := ""
		s = s + fmt.Sprintf("%s\n", lbracket)
		for _, str := range vt {
			s = s + fmt.Sprintf(
				"%s%s,\n",
				indent,
				printBody("  "+indent, "[", "]", ":", str),
			)
		}
		s = s + fmt.Sprintf("%s%s", strings.TrimPrefix(indent, "  "), rbracket)
		return s
	case map[string]interface{}:
		//fmt.Println("Map")
		s := ""
		delimiter := "\n\n"
		if lbracket == "[" {
			delimiter = ",\n"
		}
		if assignment == "=" {
			assignment = " ="
		}
		s = s + fmt.Sprintf("%s\n", lbracket)

		keys := make([]string, 0, len(vt))
		for key := range vt {
			keys = append(keys, key)
		}
		sort.Strings(keys)

		for _, v := range keys {
			s = s + fmt.Sprintf(
				"%s%s%s %s%s",
				indent,
				stripQuotes(v),
				assignment,
				printBody("  "+indent, "[", "]", ":", vt[v]),
				delimiter,
			)
		}
		s = s + fmt.Sprintf("%s%s", strings.TrimPrefix(indent, "  "), rbracket)
		return s
	case string:
		//fmt.Println("String")
		if strings.Contains(vt, "Constants") {
			return fmt.Sprintf("%s", vt)
		} else {
			return fmt.Sprintf("'%s'", stripQuotes(vt))
		}
	case bool:
		//fmt.Println("Bool")
		return fmt.Sprint(vt)
	default:
		return fmt.Sprintf("%s", reflect.TypeOf(vt))
	}
	//fmt.Sprintf("  %s = %v\n", key, value)
}
