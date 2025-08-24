package scan

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func ScanPath(path string) []string {
	// safety check in case of too many files
	if !strings.Contains(path, "/home/xbme2/") {
		fmt.Println("too many filesto scan")
		return nil
	}
	fmt.Println("begin scaning", path, "--------------------------")
	repositories := scanResurive(path)
	// saveFile := getSaveFile()
	// for _, repo := range repositories {
	// 	fmt.Println("repo :", repo)
	// }

	fmt.Println("end scaning", path, "--------------------------")
	return repositories

}

func scanResurive(path string) (repositories []string) {
	entries, err := os.ReadDir(path)
	if err != nil {
		fmt.Println("fail to open path")
		log.Fatal(err.Error())
	}
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		name := entry.Name()

		if name == "vendor" || name == "node_modules" {
			continue
		}

		if name == ".git" {
			repositories = append(repositories, filepath.Join(path, name))
			continue
		}
		repositories = append(scanResurive(filepath.Join(path, name)), repositories...)
	}

	return
}
