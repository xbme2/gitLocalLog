package scan

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/briandowns/spinner"
)

type scanResult struct {
	repositories []string
	fileNum      int
}

var ignoreDirs = map[string]bool{
	"vendor":       true,
	"node_modules": true,
	".svn":         true,
	".hg":          true,
	"build":        true,
	"dist":         true,
	"__pycache__":  true,
	".idea":        true,
	".vscode":      true,
}

func ScanPath(path string) []string {
	// safety check in case of too many files
	if path == "~/" || path == "~/home" {
		fmt.Println("too many files to scan")
		return nil
	}
	s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
	s.Suffix = fmt.Sprintf("扫描路径\"%s\"中\n", path)
	start := time.Now()
	repositories, fileNum := scanResurive(path)
	s.Stop()
	scanResult := &scanResult{
		repositories: repositories,
		fileNum:      fileNum,
	}
	fmt.Printf("共扫描%d 个文件,耗时 %v \n", scanResult.fileNum, time.Since(start))
	// saveFile := getSaveFile()
	// for _, repo := range repositories {
	// 	fmt.Println("repo :", repo)
	// }
	return scanResult.repositories

}

func scanResurive(path string) ([]string, int) {
	repositories, fileNum := make([]string, 0), 0
	entries, err := os.ReadDir(path)
	if err != nil {
		fmt.Println("fail to open path", path)
		log.Fatal(err.Error())
	}
	for _, entry := range entries {
		fileNum++
		if !entry.IsDir() {
			continue
		}
		name := entry.Name()

		if ignoreDirs[name] {
			continue
		}

		if name == ".git" {
			repositories = append(repositories, filepath.Join(path, name))
			continue
		}
		newRepositories, newFileNum := scanResurive(filepath.Join(path, name))
		repositories = append(newRepositories, repositories...)
		fileNum += newFileNum
	}

	return repositories, fileNum
}
