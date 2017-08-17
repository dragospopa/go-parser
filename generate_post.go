package main

import (
	"fmt"
	"strings"
	"errors"
	"os"
	"path/filepath"
)

func generatePost(path string, includes []string) {
	tree, cat, product, topic, err := getStuffFromPath(path)
	if err != nil {
		fmt.Errorf("you deffo are not in the right place to call this.\n")
	}

	fmt.Println(tree + "\n" + cat + "\n" + product + "\n" + topic + "\n")

	productPath, catPath, topicPath := generateTargetPath(tree, product, cat, topic)
	fmt.Println(productPath + "\n" + catPath + "\n" + topicPath + "\n")

	checkStructure(productPath, catPath, topicPath)

	
}

func getStuffFromPath(path string) (string, string, string, string, error) {
	var cat, product, topic, tree string

	dirs := strings.Split(path, "/")
	if len(dirs) < 5 {
		return "", "", "", "", errors.New("this is too close to the root of your os ffs")
	}

	topic = dirs[len(dirs)-1]
	product = dirs[len(dirs)-2]
	cat = dirs[len(dirs)-3]
	tree = "/"
	for i := 0; i < len(dirs)-5; i++ {
		tree = filepath.Join(tree, dirs[i])
	}

	return tree, cat, product, topic, nil
}

func generateTargetPath(tree, product, cat, topic string) (string, string, string) {
	product = "_" + product
	return filepath.Join(tree, product), filepath.Join(tree, product, cat), filepath.Join(tree, product, cat, topic)
}

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

func checkStructure(productPath, catPath, topicPath string){

	okProduct, _ := exists(productPath)
	okCat, _ := exists(catPath)
	okTopic, _ := exists(topicPath)
	if !okProduct {
		os.Mkdir(productPath, 0777)
		os.Mkdir(catPath, 0777)
		os.Mkdir(topicPath, 0777)
	} else if !okCat {
		os.Mkdir(catPath, 0777)
		os.Mkdir(topicPath, 0777)
	} else if !okTopic {
		os.Mkdir(topicPath, 0777)
	}
}
