package main

import (
	"flag"
	"fmt"
	"github.com/lunny/html2md"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"unicode"
	"regexp"
)

var (
	flagPath  string
	flagParse bool
	flagMove  bool
	flagRelate bool
)

func init() {
	flag.StringVar(&flagPath, "p", "", "project path")
	flag.BoolVar(&flagMove, "move", false, "call it in the exact folder where you want stuff to happen")
	flag.BoolVar(&flagParse, "parse", false, "runs starter to do nothing related to starter")
	flag.BoolVar(&flagRelate, "relate", false, "tries to work out where is what")
}

func main() {
	flag.Parse()
	var ok1, ok2, ok int
	var visited, largerVisited []string

	if flagMove {
		// OUTPUT: CREATES FILE NAMED BY THE NAME OF THE FOLDER IN THE RIGHT PLACE IN STRUCTURE OF THE ACTUAL POSTS(OR AS CLOSE AS POSSBILE)
		// THE FILE HAS THE HEADER COMPLETED AS MUCH AS HUMANLY POSSIBLE
		// {% assign product=""[NAME_OF_THE_PRODUCT]" %} - taken from the path
		// {% list of includes that matches the number of files that are note code %} - path should be completed

		var includes []string
		fPath, _ := os.Getwd()
		//EVERYTHING
		filepath.Walk(fPath, func(path string, info os.FileInfo, err error) error {
			ok++
			if info.IsDir() && ok > 1 {
				folderPath := filepath.Join(fPath, info.Name())
				//INSIDE CATEGORIES
				filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
					ok1++
					laster := strings.Split(folderPath, "/")[len(strings.Split(folderPath, "/"))-1]
					for k := 0; k < len(largerVisited); k++ {
						//fmt.Println("**** " + largerVisited[k] + " ****\n")
						if largerVisited[k] == laster {
							return nil
						}
					}
					if info.IsDir() && ok1 > 1 {
						childFolderPath := filepath.Join(folderPath, info.Name())
						//INSIDE PRODUCTS
						filepath.Walk(childFolderPath, func(path string, info os.FileInfo, err error) error {
							ok2 ++
							//fmt.Println("THIS IS " + childFolderPath)
							last := strings.Split(childFolderPath, "/")[len(strings.Split(childFolderPath, "/"))-1]
							for k := 0; k < len(visited); k++ {
								//fmt.Println("---->" + visited[k] + "<-----\n")
								if visited[k] == last {
									return nil
								}
							}
							largerVisited = append(largerVisited, info.Name())
							if info.IsDir() && ok2 > 1 {
								childchildFolderPath := filepath.Join(childFolderPath, info.Name())
								fmt.Println(childchildFolderPath)
								visited = append(visited, info.Name())
								//INSIDE TOPICS
								includes = []string{}
								filepath.Walk(childchildFolderPath, func(path string, info os.FileInfo, err error) error {
									if !info.IsDir() {
										_, _ = os.Open(info.Name())
										oldPath := filepath.Join(childchildFolderPath, info.Name())
										newPath := filepath.Join(childchildFolderPath, getFileName(info.Name()))
										os.Rename(oldPath, newPath)
										codeFile, _ := regexp.Match("code", []byte(info.Name()[:5]))
										if !codeFile {
											file, _ := ioutil.ReadFile(newPath)
											takeCareOfIncludes(string(file), newPath)
											includes = append(includes, getFileName(info.Name()))
										}
									}
									return nil
								})
								generatePost(childchildFolderPath, includes)
							}
							return nil
						})
					}
					ok2 = 0
					return nil
				})
			}
			ok1 = 0
			return nil
		})
	}

	if flagParse {
		//use -p to set dir

		if flagPath != "" {

			filepath.Walk(flagPath, func(path string, info os.FileInfo, err error) error {
				if !info.IsDir() {

					_, _ = os.Open(info.Name())

					text, _ := ioutil.ReadFile(flagPath + info.Name())

					html2md.AddRule("pre", &html2md.Rule{
						Patterns: []string{"pre"},
						Tp:       0,
						Replacement: func(innerHTML string, attrs []string) string {
							if len(attrs) > 1 {
								return "```" + attrs[1] + "```\n"
							}
							return ""
						},
					})

					md := html2md.Convert(string(text))

					generateFilesFromThis(md, info.Name()[:len(info.Name())-3])
				} else {

				}
				return nil
			})
		} else {
			fmt.Printf("lol - give me a path usign the -p flag boii")
		}

		return
	}
	if flagRelate {
		fmt.Printf("\n")
		if flagPath != "" {
			mapperino := make(map[string][]string)
			mappy := &mapperino
			OperinoNoperino := make(map[string]int)
			bop := &OperinoNoperino
			filepath.Walk(flagPath, func( path string, info os.FileInfo, err error,) error {
				if !info.IsDir() {
					OperinoNoperinolocal := *bop
					if _, ok := OperinoNoperinolocal[info.Name()]; ok {
						OperinoNoperinolocal[info.Name()] ++
					} else{
						_, _ = os.Open(path)
						text, e := ioutil.ReadFile(path)
						if e != nil {
							fmt.Printf(e.Error())
						}
						OperinoNoperino[info.Name()] = 1
						lookForIncludes(string(text), info.Name(), mappy, path)
					}
				}
				return nil
			})
			fuckthis(mapperino)
			/*for inline, pages := range mapperino {
				for _, page := range pages {
					fmt.Printf(inline + "\n is used in \n" + page + "\n\n")
				}
			}*/
		} else {
			fmt.Printf("lol - give me a path usign the -p flag boii")
		}
		return
	}
}
func lookForIncludes(text string, filename string, mapAddr *(map[string][]string), path string) {
	r, _ := regexp.Compile("include (_inline.*md)")
	res := r.FindAllStringSubmatch(text, -1)
	for _, element := range res {
		mapperino := *mapAddr
		mapperino[element[1]] = append(mapperino[element[1]], path);
	}

}

func generateFilesFromThis(text string, filename_base string) {
	var inlineContent, filename string
	if len(text) < 2 {
		return
	}
	os.Mkdir(filename_base, 0777)

	text += "\n\n"

	header := hasHeader(text)
	if header {
		for i := 0; i < len(text)-4; i++ {
			if text[i] == '-' && text[i+1] == '-' && text[i+2] == '-' {
				for i += 6; i < len(text); i++ {
					if text[i] == '-' && text[i-1] == '-' && text[i-2] == '-' {
						break
					}
				}
				text = "\n#" + text[i+3:]
				break
			}
		}
	}

	for i := 1; i < len(text); i++ {
		if text[i] == '#' && text[i-1] == '\n' {
			filename = filename_base + "_"
			inlineContent = ""
			for ; i < len(text); i++ {
				if text[i] == '#' || text[i] == ' ' || text[i] == '\n' {
					inlineContent = string(append([]byte(inlineContent), text[i]))
				} else {
					break
				}
			}
			for ; i < len(text); i++ {
				if text[i] == '\n' {
					break
				} else {
					inlineContent = string(append([]byte(inlineContent), text[i]))
					if text[i] != '#' && text[i] != '\n' {
						filename = string(append([]byte(filename), text[i]))
					}
				}
			}
			filename = getFileName(filename)

			for ; i < len(text)-1; i++ {
				if text[i+1] == '#' && text[i] == '\n' {
					inlineContent = string(append([]byte(inlineContent), text[i]))
					break
				} else {
					inlineContent = string(append([]byte(inlineContent), text[i]))
				}
			}

			if inlineContent != "" {
				if len(filename) > 59 {
					filename = filename[:59] + ".md"
				} else {
					filename += ".md"
				}
				inlineContent = "<!-- post: -->\n\n" + inlineContent
				for a := 0; a < len(inlineContent)-3; a++ {
					preInlineContent := ""
					var l, r int
					if inlineContent[a] == '`' && inlineContent[a+1] == '`' && inlineContent[a+2] == '`' {
						j := a
						l = j
						j += 6
						for ; j < len(inlineContent); j++ {
							if inlineContent[j] == '`' && inlineContent[j-1] == '`' && inlineContent[j-2] == '`' {
								r = j
								break
							}
						}
						if r != 0 {
							for k := l; k <= r+1; k++ {

								preInlineContent = string(append([]byte(preInlineContent), inlineContent[k]))
							}
							preInlineContent = preInlineContent[3:][:len(preInlineContent)-7]
							if len(preInlineContent) > 7 {
								code_filename := "code_" + filename[:len(filename)-3] + "-"
								if len(preInlineContent) > 15 {
									for k := 3; k <= 15; k++ {
										if unicode.IsLower(rune(preInlineContent[k])) {
											code_filename = string(append([]byte(code_filename), preInlineContent[k]))
										}
									}
								} else {
									code_filename += "code"
								}

								preInlineContent = "<!-- layout: code\npost: " + filename + " -->\n" + preInlineContent

								code_filename = strings.Trim(strings.Trim(code_filename, "\n"), "<") + ".md"

								_, er := os.Open(filename_base + "/" + string(code_filename))
								if er == nil {
									code_filename = code_filename[:len(code_filename)-3] + "-2.md"
								}
								_, er = os.Open(filename_base + "/" + string(code_filename))
								if er == nil {
									code_filename = code_filename[:len(code_filename)-3] + "-3.md"
								}
								_, er = os.Open(filename_base + "/" + string(code_filename))
								if er == nil {
									code_filename = code_filename[:len(code_filename)-3] + "-4.md"
								}
								_, er = os.Open(filename_base + "/" + string(code_filename))
								if er == nil {
									code_filename = code_filename[:len(code_filename)-3] + "-5.md"
								}

								if len(code_filename) > 59 {
									code_filename = code_filename[:59] + ".md"
								}
								include := filepath.Join(filename_base, code_filename)
								err := ioutil.WriteFile(filepath.Join(filename_base, code_filename), []byte(preInlineContent), 0644)
								if err != nil {
									fmt.Println("You fucked up: ", err, "\n\n\n")
								}
								fmt.Println(code_filename)
								inlineContent = inlineContent[:l] + "\n\n{%include _inlines/" + include + " %}\n\n" + inlineContent[r+1:]
								r = 0
								l = 0
							}
						}
					}
				}

				fmt.Println(filename)
				//common.PrintlnWarning(inlineContent)
				err := ioutil.WriteFile(filepath.Join(filename_base, filename), []byte(inlineContent), 0644)
				if err != nil {
					fmt.Println("You fucked up: ", err, "\n\n\n")
				}
			}

		}
	}
}

func getFileName(filename string) string {
	if !(unicode.IsLower(rune(filename[0])) || unicode.IsUpper(rune(filename[0]))) {
		//		filename = filename[1:]
	}
	filename = strings.ToLower(strings.Replace(filename, " ", "-", -1))
	filename = strings.Trim(filename, "\n")
	filename = strings.Trim(filename, "\\")
	for i := 0; i < len(filename); i++ {
		if isSpecial(filename[i]) {
			if i < len(filename)-1 {
				filename = filename[:i] + filename[i+1:]
			} else {
				filename = filename[:i]
			}
			i--
		}
	}

	return filename
}
func isSpecial(c byte) bool {
	special := []byte{':', '{', '}', '[', ']', ',', '&', '*', '#', '?', '|', '<', '>', '=', '!', '%', '@', '\\', '/', '\'', '(', ')', '"'}
	for i := 0; i < len(special); i++ {
		if special[i] == c {
			return true
		}
	}
	return false
}

func takeCareOfIncludes(file string, fpath string) {
	dirs := strings.Split(fpath, "/")
	p := 0;
	j := 0
	for p = 0; p < len(dirs); p++ {
		if dirs[p] == "_inlines" {
			break
		}
	}
	includePath := ""
	tempIncludePath := ""
	includedCode := ""
	for ; p < len(dirs)-1; p++ {
		includePath = filepath.Join(includePath, dirs[p])
	}
	lines := strings.Split(file, "\n")
	for i := 0; i < len(lines); i++ {
		includeLine, _ := regexp.MatchString("%include ", lines[i])
		if includeLine {
			includedCode = ""
			for k := len(lines[i]) - 3; k >= 0; k-- {
				if lines[i][k] == '/' {
					break
				} else {
					includedCode = string(append([]byte(includedCode), lines[i][k]))
				}
			}
			includedCode = reverse(includedCode)
			includedCode = getFileName(includedCode)
			//fmt.Println("THE CODE IN HERE IS: %s\n", includedCode)
			tempIncludePath = filepath.Join(includePath, includedCode[:len(includedCode)-1])
			for j = 0; j < len(lines[i]); j++ {
				if lines[i][j] == '_' {
					break
				}
			}
			lines[i] = lines[i][:j] + tempIncludePath + " %}\n"
		}
	}
	text := strings.Join(lines, "\n")
	err := ioutil.WriteFile(fpath, []byte(text), 0777)
	if err != nil {
		fmt.Errorf("YOU GOT AN ERROR: %s\n", err)
	}
}

func hasHeader(text string) bool {
	for i := 0; i < 15; i++ {
		if text[i] == '-' && text[i+1] == '-' && text[i+2] == '-' {
			return true
		}
	}
	return false
}

func reverse(s string) string {
	r := []rune(s)
	for i, j := 0, len(r)-1; i < len(r)/2; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}
