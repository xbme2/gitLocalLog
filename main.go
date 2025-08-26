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
	path    string
	offset  int
	verbose bool
)

func main() {

	flag.StringVar(&path, "path", "./", " 用于搜索的文件目录")
	flag.IntVar(&offset, "month", 6, "查询月份数量")
	flag.BoolVar(&verbose, "verbose", true, "启用详细输出")
	flag.Parse()
	var repositories []string
	fmt.Println(offset)
	// stats.PrintMonths(time.Now(), offsetNum)
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
		stats.GenerateStats(repositories, offset, verbose)
	} else {
		fmt.Println("no git repo under this path")
	}
}
