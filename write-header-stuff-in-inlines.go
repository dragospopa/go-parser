package main

import (
	"fmt"
	"os"
	"io/ioutil"
	"path/filepath"
)

const GITHUBPATH = "https://github.com/cloud66/help/edit/feature/inlines/_includes/"

func fuckthis(mapperino map[string][]string) {
	for inline, pages := range mapperino {
		inline = filepath.Join("/Users/dragos/work/help/_includes", inline)
		pagez := ""
		for _, page := range pages {
			dirs, _ := filepath.Split(page)
			dirs = dirs[24:]
			page = filepath.Join(dirs)
			if pagez != "" {
				pagez = pagez + ", " + page
			} else {
				pagez = "[ " + page
			}
		}
		pagez += "]"
		fmt.Println(pagez)
		_, er := os.Open(inline)
		if er != nil {
			fmt.Errorf("You obviously fucked up ...because...%s\n", er)
		}
		text, _ := ioutil.ReadFile(inline)
		if len(text) > 3 {
			text = []byte("<!-- usedin: " + pagez + " " + string(text[3:]))
		}
		err := ioutil.WriteFile(inline, text, 0777)
		if err != nil {
			fmt.Errorf("That's deffo not gonna print.\n")
			break
		}
	}
}

func populateGitLinks(mapperino map[string][]string) {
	for page, includes := range mapperino {
		gitlinks := "[ "
		for _, include := range includes {
			include = GITHUBPATH + include[:len(include)-3] + ".html"
			if gitlinks == "[ " {
				gitlinks += include
			} else {
				gitlinks += ", " + include
			}
		}
		gitlinks += " ]"
		_, err := os.Open(page)
		if err != nil {
			fmt.Errorf("broken.\n")
			break
		}
		fmt.Println(gitlinks)
		text, _ := ioutil.ReadFile(page)
		text = []byte("---\n" + "gitlinks: " + gitlinks + "\n" + string(text[4:]))
		err = ioutil.WriteFile(page, text, 0777)
		if err!=nil{
			fmt.Errorf("It didnt write!\n")
			break
		}
	}
}
