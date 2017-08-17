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
)

func init() {
	flag.StringVar(&flagPath, "p", "", "project path")
	flag.BoolVar(&flagMove, "move", false, "call it in the exact folder where you want stuff to happen")
	flag.BoolVar(&flagParse, "parse", false, "runs starter to do nothing related to starter")
}

func main() {
	flag.Parse()
	var ok int

	if flagMove {
		// OUTPUT: CREATES FILE NAMED BY THE NAME OF THE FOLDER IN THE RIGHT PLACE IN STRUCTURE OF THE ACTUAL POSTS(OR AS CLOSE AS POSSBILE)
		// THE FILE HAS THE HEADER COMPLETED AS MUCH AS HUMANLY POSSIBLE
		// {% assign product=""[NAME_OF_THE_PRODUCT]" %} - taken from the path
		// {% list of includes that matches the number of files that are note code %} - path should be completed

		var includes []string
		folderPath, _ := os.Getwd()
		filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
			ok ++
			if info.IsDir() && ok > 1 {
				childFolderPath := filepath.Join(folderPath, info.Name())
				fmt.Println(childFolderPath)
				filepath.Walk(childFolderPath, func(path string, info os.FileInfo, err error) error {
					if !info.IsDir() {
						_, _ = os.Open(info.Name())
						codeFile, _ := regexp.Match("code", []byte(info.Name()[:5]))
						if !codeFile {
							includes = append(includes, info.Name())
						}
					}
					return nil
				})
				generatePost(childFolderPath, includes)
			}
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
		if filename[i] == '/' || filename[i] == '$' || filename[i] == '*' || filename[i] == ':' || filename[i] == '?' || filename[i] == '!' {
			filename = filename[:i]
			if i < len(filename)-1 {
				filename += filename[i+1:]
			}
		}
	}
	filename = strings.Trim(filename, "/")
	filename = strings.Trim(filename, "$")
	filename = strings.Trim(filename, "*")
	filename = strings.Trim(filename, ":")

	return filename
}

func hasHeader(text string) bool {
	for i := 0; i < 15; i++ {
		if text[i] == '-' && text[i+1] == '-' && text[i+2] == '-' {
			return true
		}
	}
	return false
}
