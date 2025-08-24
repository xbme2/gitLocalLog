package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"main.go/scan"
	"main.go/stats"
)

var (
	path  string
	month string
)

func main() {

	flag.StringVar(&path, "path", "", " 用于搜索的文件目录")
	flag.StringVar(&month, "month", "6", "查询月份数量")
	flag.Parse()
	var repositories []string
	if path != "" {
		if path == "./" {
			path, _ = os.Getwd()
		} else if path == "../" {
			currentPath, _ := os.Getwd()
			path = filepath.Dir(currentPath)
		}
		repositories = scan.ScanPath(path)
	} else {
		fmt.Println("your paramter is wrong . Enter -h for help")
	}
	if len(repositories) > 0 {
		stats.GenerateStats(repositories)
	}
}
