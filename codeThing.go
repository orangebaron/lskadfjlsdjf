package main

import "fmt"
import "strings"
import "os"

func err(e string) {
	fmt.Println(e)
	os.Exit(1)
}

type ast struct { //abstract syntax tree
	desc   string
	leaves []*ast
}

func removeComments(code string) string {
	//do things
	return code
}
func generateSyntaxTrees(code string) []*ast { //TODO: error handling?
	s := strings.Split(code, "\n")
	retVal := []*ast{}
	var currentAst *ast
	for _, line := range s {
		appendTo := currentAst //ast to append to
		amtOfTabs := 0
		for len(line) > 0 && line[0] == '\t' {
			if amtOfTabs != 0 {
				if appendTo == nil || len(appendTo.leaves) == 0 {
					continue
				}
				appendTo = appendTo.leaves[len(appendTo.leaves)-1]
			}

			line = line[1:]
			amtOfTabs++
		}
		if amtOfTabs == 0 && appendTo != nil {
			retVal = append(retVal, appendTo)
			appendTo = nil
		}

		firstSpace := strings.Index(line, " ")
		firstWord := ""
		afterwards := ""
		if firstSpace == 0 {
			continue
		} else if firstSpace == -1 {
			firstWord = line
		} else {
			firstWord = line[:firstSpace]
			afterwards = line[firstSpace+1:]
		}
		thisAst := new(ast)
		if firstWord == "loop" || firstWord == "if" || firstWord == "elif" || firstWord == "else" {
			thisAst.desc = firstWord
			thisAst.leaves = []*ast{new(ast)}
			if afterwards == "" {
				thisAst.leaves[0].desc = "E true"
			} else {
				thisAst.leaves[0].desc = "E " + afterwards
			}
		} else if afterwards == "" {
			thisAst.desc = "E " + firstWord
		} else {
			thisAst.desc = firstWord
			thisAst.leaves = []*ast{new(ast)}
			thisAst.leaves[0].desc = "EV " + afterwards
		}
		fmt.Println(thisAst)
		if appendTo == nil {
			currentAst = thisAst
		} else {
			appendTo.leaves = append(appendTo.leaves, thisAst)
		}
	}
	if currentAst != nil {
		retVal = append(retVal, currentAst)
	}
	return retVal
}

type varType struct {
	desc string //make this more elaborate later
}

var tempAstVals map[*ast]interface{}
var tempAstTypes map[*ast]varType

func evalSyntaxTree(a *ast, parents []*ast) {
	fmt.Println(a.desc)
	if len(a.desc) > 2 && a.desc[:2] == "E " { //e for eval
		word := a.desc[2:]
		var evalVal interface{}
		var evalType varType
		if word == "true" {
			evalVal = true
			evalType = varType{"bool"}
		} else if word == "false" {
			evalVal = false
			evalType = varType{"bool"}
		} else if word == "nil" {
			evalVal = nil
			evalType = varType{"nil"}
		} else {
			//err("UNDEFINED: "+word)
		}
		tempAstVals[a] = evalVal
		tempAstTypes[a] = evalType
		a.desc="D" //d for done
	}
	for _, a2 := range a.leaves {
		evalSyntaxTree(a2, append(parents,a))
	}
}
func main() {
	code := "loop\n\tloop a\n\t\tb\nc"
	code = removeComments(code)
	parent := &ast{"parent", generateSyntaxTrees(code)}
	tempAstVals = make(map[*ast]interface{})
	tempAstTypes = make(map[*ast]varType)
	evalSyntaxTree(parent, []*ast{})
	fmt.Println("---------------------------")
	evalSyntaxTree(parent, []*ast{})
	fmt.Println(parent.leaves)
}
