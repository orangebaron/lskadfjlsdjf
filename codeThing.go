package main

import "fmt"
import "strings"

type ast struct { //abstract syntax tree
	desc   string
	leaves []*ast
}

func removeComments(code string) string {
	//do things
	return code
}
func generateSyntaxTrees(code string) []*ast {
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
				thisAst.leaves[0].desc = "EVAL: true"
			} else {
				thisAst.leaves[0].desc = "EVAL: " + afterwards
			}
		} else if afterwards == "" {
			thisAst.desc = "EVAL: " + firstWord
		} else {
			thisAst.desc = firstWord
			thisAst.leaves = []*ast{new(ast)}
			thisAst.leaves[0].desc = "EVAL: " + afterwards
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
func main() {
	fmt.Println(generateSyntaxTrees("loop\n\tloop a\n\t\ta\na"))
}
