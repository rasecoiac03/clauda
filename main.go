package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/rasecoiac03/clauda/pkg/config"
)

var fileFolder = "files/"
var files = []string{
	"lotomania",
	"lotofacil",
}

func main() {
	downloadFiles()
	readFiles()
}

func downloadFiles() {
	executeCmd("files removed", exec.Command("rm", "-rf", fileFolder))
	executeCmd("files folder exists", exec.Command("mkdir", "-p", fileFolder))
	for _, file := range files {
		filePath := config.GetConfig(fmt.Sprintf("%s_file_enpoint", file))
		filePathParts := strings.Split(filePath, "/")
		fileName := filePathParts[len(filePathParts)-1]
		executeCmd(fmt.Sprintf("file downloaded: %s", fileName), exec.Command("wget", filePath))
		executeCmd(fmt.Sprintf("file unzipped: %s", fileName), exec.Command("unzip", "-u", fileName))
		executeCmd(fmt.Sprintf("file moved to: %s%s.HTM", fileFolder, fileName), exec.Command("/bin/bash", "-c", fmt.Sprintf("mv *.HTM %s%s.HTM", fileFolder, file)))
		executeCmd(fmt.Sprintf("file zip removed: %s", fileName), exec.Command("rm", "-f", fileName))
		executeCmd("cleaning", exec.Command("/bin/bash", "-c", "rm -f *.GIF"))
	}
}

func executeCmd(message string, cmd *exec.Cmd) {
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	log.Println(message)
}

func readFiles() {
	for _, f := range files {
		file, err := os.Open(fmt.Sprintf("%s%s.HTM", fileFolder, f))
		if err != nil {
			log.Fatal(err)
		}
		doc, qErr := goquery.NewDocumentFromReader(file)
		if qErr != nil {
			log.Fatal(err)
		}

		colIni := config.GetIntConfig(fmt.Sprintf("%s_col_ini", f))
		colEnd := config.GetIntConfig(fmt.Sprintf("%s_col_end", f))

		count := make(map[int]int64)
		doc.Find("table tr").Each(func(i int, s *goquery.Selection) {
			s.Find("td").Each(func(j int, td *goquery.Selection) {
				if key, err := strconv.Atoi(td.Text()); err == nil && j >= colIni && j <= colEnd {
					count[key]++
				}
			})
		})
		var keys []int
		for k := range count {
			keys = append(keys, k)
		}
		sort.Ints(keys)

		log.Printf("Results for %s\n", strings.ToUpper(f))
		for _, k := range keys {
			log.Printf("%d: %d\n", k, count[k])
		}
	}
}
