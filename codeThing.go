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
		} else if firstWord == "const" {
			thisAst.desc = line
		} else {
			thisAst.desc = "E " + line
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
type variable struct {
	vType varType
	vValue interface{}
}

var constVars map[string]*variable
var tempAstVars map[*ast]*variable
var astsFinished map[*ast]bool

func evalSyntaxTree(a *ast, parents []*ast) {
	fmt.Println(a.desc)
	if !astsFinished[a] {
		evalBranches := true
		done := true
		if len(a.desc) > 2 && a.desc[:2] == "E " { //e for eval
			word := a.desc[2:]
			var evalVar variable
			if word == "true" {
				evalVar = variable{varType{"bool"},true}
			} else if word == "false" {
				evalVar = variable{varType{"bool"},false}
			} else if word == "nil" {
				evalVar = variable{varType{"nil"},nil}
			} else {
				//err("UNDEFINED: "+word)
			}
			tempAstVars[a] = &evalVar
		} else if len(a.desc) > 6 && a.desc[:6] == "const " {
			word := a.desc[6:]
			equLoc := strings.Index(word,"=") //TODO: handle errors related to the = sign
			if len(a.leaves) == 0 {
				a.leaves = []*ast{&ast{"E "+word[equLoc+2:],[]*ast{}}}
			}
			evalSyntaxTree(a.leaves[0],append(parents, a))
			if !astsFinished[a.leaves[0]] {
				done = false
			} else {
				constVars[word[:equLoc-1]] = tempAstVars[a.leaves[0]]
			}
			evalBranches = false
		}
		if evalBranches {
			for _, a2 := range a.leaves {
				evalSyntaxTree(a2, append(parents, a))
				done = done && astsFinished[a2]
			}
		}
		astsFinished[a] = done
	}
}
func main() {
	code := "const a = true"
	code = removeComments(code)
	parent := &ast{"parent", generateSyntaxTrees(code)}
	tempAstVars = make(map[*ast]*variable)
	constVars = make(map[string]*variable)
	astsFinished = make(map[*ast]bool)
	evalSyntaxTree(parent, []*ast{})
	fmt.Println(constVars["a"])
	fmt.Println(parent.leaves)
}
