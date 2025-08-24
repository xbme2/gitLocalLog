package stats

import (
	"fmt"
	"log"
	"time"

	"github.com/go-git/go-git/v6"
	"github.com/go-git/go-git/v6/plumbing/object"
)

// with go modules disabled

func calculateDaysSince(date time.Time) int {
	days := 0
	now := time.Now()
	for date.Before(now) {
		days++
		date = date.Add(24 * time.Hour)
	}
	return days
}

func GenerateStats(repositories []string) {
	// Clone the given repository to the given directory

}

func fillMap(repositories []string, commits map[int]int) map[int]int {

	for _, repo := range repositories {
		r, err := git.PlainOpen(repo)
		if err != nil {
			fmt.Println("not exist the repo——", repo)
			log.Fatal(err.Error())
		}
		ref, err := r.Head()
		if err != nil {
			log.Fatal(err.Error())
		}
		iterator, err := r.Log(&git.LogOptions{From: ref.Hash()})
		if err != nil {
			log.Fatal(err.Error())
		}
		// offset := calc
		err = iterator.ForEach(func(c *object.Commit) error {
			fmt.Println(c.Author.When)
			
			return nil
		})
	}

}
