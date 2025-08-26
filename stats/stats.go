package stats

import (
	"fmt"
	"log"
	"time"

	"github.com/briandowns/spinner"
	"github.com/fatih/color"
	"github.com/go-git/go-git/v6"
	"github.com/go-git/go-git/v6/plumbing/object"
)

// with go modules disabled

var NumberMonth = map[int]string{
	1:  "Jan",
	2:  "Feb",
	3:  "Mar",
	4:  "Apr",
	5:  "May",
	6:  "Jun",
	7:  "Jul",
	8:  "Aug",
	9:  "Sep",
	10: "Oct",
	11: "Nov",
	12: "Dec",
}

type Stats struct {
	Commits        map[int]int
	TotalAdditions int
	TotalDeletions int
	TotalComitNum  int
	RepoNum        int
	vervose        bool
}

// var NumberWeek = map[int]string{
// 	0: "Sun",
// 	1: "Mon",
// 	2:"Fer"
// }

func getStartOfDay(date time.Time) time.Time {
	return time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
}

func calculateDaysSince(date time.Time) int {
	days := 0
	now := getStartOfDay(time.Now())
	for date.Before(now) {
		days++
		date = date.Add(24 * time.Hour)
	}
	return days
}

func PrintMonths(now time.Time, monthOffset int) {
	start, _ := startDayRecentMonths(now, monthOffset)

	month := start.Month()
	fmt.Println(month.String())
	fmt.Printf("       ")
	// fmt.Printf(month.String()[:3])

	for {
		if start.Month() != month {
			month = start.Month()
			fmt.Printf(month.String()[:3])
		} else {
			fmt.Printf("   ")
		}
		start = start.Add(time.Hour * 24 * 7)
		if start.After(now) {
			break
		}
	}

	fmt.Print("\n")
}

func PrintDays(commits map[int]int) {

}

func PrintWeek(weekIndex int) {
	switch weekIndex {
	case 1:
		fmt.Print("Mon    ")
	case 3:
		fmt.Print("Wed    ")
	case 5:
		fmt.Print("Fri    ")
	default:
		fmt.Print("       ")
	}
}

func printDay(commitNum int) {
	switch {
	case commitNum == 0:
		color.New(color.FgBlack).Print("--") // 空格子
		fmt.Printf(" ")
	case commitNum < 2:
		// fmt.Printf(" ")
		color.New(color.BgBlack).Add(color.BgWhite).Print("+", commitNum)
		// fmt.Printf(" ")
	case commitNum < 5:
		// fmt.Printf(" ")
		color.New(color.BgHiGreen).Add(color.BgHiGreen).Print("+", commitNum)
		// fmt.Printf(" ")

	case commitNum < 10:
		color.New(color.BgGreen).Add(color.BgWhite).Print("+", commitNum)
	default:
		color.New(color.BgGreen, color.FgHiWhite).Print(commitNum) // 深绿
		// color.New(color.BgGreen, color.FgWhite).Print("██") // 最深绿
	}
}

func detailShow(stat *Stats, start time.Time) {
	fmt.Printf("分析报告  (%s to %s)\n", formatTime(start), formatTime(time.Now()))
	fmt.Printf("============================\n")
	fmt.Printf("共扫描%d个仓库\n", stat.RepoNum)
	fmt.Printf("共计提交%d次\n", stat.RepoNum)

	fmt.Printf("总计增加行数: +%d\n", stat.TotalAdditions)
	fmt.Printf("总计删除行数: +%d\n", stat.TotalDeletions)

}

func printLegend() {
	fmt.Printf("Less  ")
	color.New(color.FgBlack).Print("--")
	fmt.Printf("  ")
	color.New(color.BgBlack).Add(color.BgWhite).Print("  ")
	fmt.Printf("  ")
	color.New(color.BgHiGreen).Add(color.BgHiGreen).Print("  ")
	fmt.Printf("  ")
	color.New(color.BgGreen).Add(color.BgGreen).Print("  ")
	fmt.Printf("More \n")
}

func uiShow(commits map[int]int, offset int, daysTotal int, start time.Time) {
	now := time.Now()
	PrintMonths(now, offset)
	for i := 0; i < 7; i++ {
		PrintWeek(i)
		// tempStart := start.AddDate(0, 0, i)
		// month := tempStart.Month()
		for j := daysTotal - i; j > 0; j -= 7 {
			// if tempStart.Month() != month {
			// 	fmt.Printf()
			// }
			printDay(commits[j])
		}
		fmt.Print("\n")

		fmt.Printf("\n")
	}
	printLegend()
	// for offset, num := range commits {
	// 	if num > 2 {

	// 	}
	// }
}

func GenerateStats(repositories []string, month int, verbose bool) {
	stat := &Stats{
		Commits:        make(map[int]int),
		TotalAdditions: 0,
		TotalDeletions: 0,
		RepoNum:        0,
		vervose:        verbose,
	}
	s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
	s.Suffix = "正在分析仓库中..."
	s.Start()
	start, days := startDayRecentMonths(time.Now(), month)
	fillMap(repositories, stat, start)
	uiShow(stat.Commits, month, days, start)
	if verbose {
		detailShow(stat, start)
	}
	s.Stop()
	// PrintContribGraph(commits, days)

}

func fillMap(repositories []string, stat *Stats, since time.Time) {
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
		iterator, err := r.Log(&git.LogOptions{From: ref.Hash(), Since: &since})
		if err != nil {
			log.Printf("警告：无法获取仓库 %s 的日志: %v", repo, err)
			continue
		}
		// offset := calc
		err = iterator.ForEach(func(c *object.Commit) error {
			// fmt.Println(calculateDaysSince(c.Author.When))
			stat.RepoNum++
			stat.Commits[calculateDaysSince(c.Author.When)]++
			stat.TotalComitNum++
			if stat.vervose {
				commitStats, err := c.Stats()
				if err != nil {
					log.Fatal(err.Error())
				}
				for _, comcommitStat := range commitStats {
					addition, delection := comcommitStat.Addition, comcommitStat.Deletion
					stat.TotalAdditions += addition
					stat.TotalDeletions += delection
				}
			}
			return nil
		})
	}
}
