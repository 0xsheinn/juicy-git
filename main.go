package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func writeFile(lst []string, fileName string) {
	file, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}
	file.Close()

	files, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	for _, data := range lst {
		fmt.Fprint(files, data)
	}

	file.Close()
}

func EncodeParam(s string) string {
	return url.QueryEscape(s)
}

func main() {
	arguments := os.Args[1:]
	if len(arguments) < 2 || arguments[0] == "-h" {
		fmt.Print("[!] Usage: ./juicy-git [organization name] [output filename]\n\n")
		os.Exit(0)
	}

	org := arguments[0]
	fileName := arguments[1]

	resp, err := http.Get("https://raw.githubusercontent.com/vsec7/gitdorkhelper/main/all_dorks.txt")
	if err != nil {
		fmt.Println("[!] Request time Out")
		os.Exit(0)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	dorks := string(body)
	dorks_list := strings.Split(dorks, "\r\n")

	lst := make([]string, 0)
	for _, dork := range dorks_list {
		github := "https://github.com/search?q="
		query := fmt.Sprintf("\"%s\" %s", org, dork)
		end_query := "&type=code"
		query_encode := EncodeParam(query)
		fmt.Printf("%s%s%s [\033[31m%s\033[0m]\n", github, query_encode, end_query, dork)
		out := github + query_encode + end_query + "\n"
		lst = append(lst, out)

	}
	writeFile(lst, fileName)

}
