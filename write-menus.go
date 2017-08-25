package main

import (
	"fmt"
	"os"
	"io/ioutil"
	"regexp"
	"path/filepath"
)

func writeMenus(mapperino map[string][]string) {
	for page, includes := range mapperino {
		links := []string{}
		menuheaders := "[ "
		menulinks:=[]string{}
		for _, include := range includes {
			menulinks= append(menulinks, include)
			r, _ := regexp.Compile(".*include (_inline.*md)")
			res := r.FindAllStringSubmatch(include, -1)
			links = append(links, res[0][1])
			inlinePath := filepath.Join("/Users/dragos/work/help/_includes", res[0][1])
			_, err := os.Open(inlinePath)
			if err != nil {
				fmt.Errorf("----> did not open!!!!\n")
			}
			text, _ := ioutil.ReadFile(inlinePath)
			heads, _ := regexp.Compile("# (.*)\n")
			resHeads := heads.FindAllStringSubmatch(string(text), -1)
			if len(resHeads) > 0 {
				if menuheaders == "[ " {
					menuheaders += "\"" + resHeads[0][1] + "\""
				} else {
					menuheaders += ", " + "\"" + resHeads[0][1] + "\""
				}
			}
		}
		menuheaders += " ]"
		_, err := os.Open(page)
		if err != nil {
			fmt.Errorf("broken.\n")
			break
		}
		isCode := false
		for _, link := range menulinks{
			code, _ := regexp.Match("code_", []byte(link))
			if code {
				isCode = true
				break
			}
		}
		if !isCode && menuheaders != "[  ]" {
			fmt.Println(menuheaders)

			pageText, _ := ioutil.ReadFile(page)
			j := 0
			for ; j < len(pageText); j++ {
				if pageText[j] == '\n' {
					break
				}
			}
			pageText = []byte(string(pageText[:j+1])+"menuheaders: " + menuheaders + "\n" + string(pageText[j+1:]))

			for j = 0; j < len(pageText); j++ {
				if pageText[j] == '{' {
					break
				}
			}
			pageText = pageText[:j]
			for _, href := range menulinks {
				pageText = []byte(string(pageText)+"\n" + href)
			}
			err = ioutil.WriteFile(page, pageText, 0777)
			if err != nil {
				fmt.Errorf("It didnt write!\n")
				break
			}
		}
	}
}
