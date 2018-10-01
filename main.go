package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/rasecoiac03/clauda/pkg/config"
)

var lock sync.Mutex

var timeSample = flag.Int("time-sample", 0, "time sample")
var timeSampleUnit = flag.String("time-sample-unit", "", "time sample unit")
var sortCount = flag.Bool("sort-count", false, "sort by count")

var fileFolder = "files/"
var files = []string{
	// "lotomania",
	// "megasena",
	// "lotofacil",
	"quina",
}

// ValidateTime type
type ValidateTime func(hours float64, sample int) bool

var validations = map[string]ValidateTime{
	"year": func(hours float64, sample int) bool {
		t := int(hours / 24 / 30 / 12)
		return t < sample
	},
	"month": func(hours float64, sample int) bool {
		t := int(hours / 24 / 30)
		return t < sample
	},
	"day": func(hours float64, sample int) bool {
		t := int(hours / 24)
		return t < sample
	},
}

func main() {
	flag.Parse()
	downloadFiles()
	readFiles()
}

func downloadFiles() {
	executeCmd("files removed", exec.Command("rm", "-rf", fileFolder))
	executeCmd("files folder exists", exec.Command("mkdir", "-p", fileFolder))
	for _, file := range files {
		filePath := config.GetConfig(fmt.Sprintf("%s_file_endpoint", file))
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
	timeValidation := *timeSample != 0 && *timeSampleUnit != ""
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

		count := make(map[int]int)
		doc.Find("table tr").Each(func(i int, s *goquery.Selection) {
			var valid = !timeValidation
			s.Find("td").Each(func(j int, td *goquery.Selection) {
				if timeValidation {
					if t, err := time.Parse("02/01/2006", td.Text()); j == 1 && err == nil {
						validation := validations[*timeSampleUnit]
						valid = validation(time.Since(t).Hours(), *timeSample)
					}
				}
				if key, err := strconv.Atoi(td.Text()); valid && err == nil && j >= colIni && j <= colEnd {
					lock.Lock()
					count[key]++
					lock.Unlock()
					if j == colEnd {
						valid = false
					}
				}
			})
		})
		if *sortCount {
			countLogSortingByValue(f, count)
		} else {
			countLogSortingByKey(f, count)
		}
	}
}

func countLogSortingByValue(file string, count map[int]int) {
	countByValue := make(map[int][]int)
	for k, v := range count {
		if _, exists := countByValue[v]; exists {
			countByValue[v] = append(countByValue[v], k)
		} else {
			countByValue[v] = []int{k}
		}
	}
	var keys []int
	for k := range countByValue {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	log.Printf("Results for %s\n", strings.ToUpper(file))
	for _, k := range keys {
		for _, v := range countByValue[k] {
			log.Printf("%d: %d\n", k, v)
		}
	}
}

func countLogSortingByKey(file string, count map[int]int) {
	var keys []int
	for k := range count {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	log.Printf("Results for %s\n", strings.ToUpper(file))
	for _, k := range keys {
		log.Printf("%d: %d\n", k, count[k])
	}
}
