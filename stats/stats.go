package stats

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"time"

	"github.com/briandowns/spinner"
	"github.com/fatih/color"
	"github.com/go-git/go-git/v6"
	"github.com/go-git/go-git/v6/plumbing/object"
	"github.com/olekukonko/tablewriter"
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

type SortType int

const (
	Default SortType = iota
	CommitNum
	LineNum
	DeletionNum
)

type AuthorStat struct {
	Name           string
	Commits        int
	TotalAdditions int
	TotalDeletions int
}

type Stats struct {
	CommitsByDay map[int]int
	// CommitsByAuthor map[string]int
	Authors        map[string]*AuthorStat
	TotalAdditions int
	TotalDeletions int
	TotalComitNum  int
	RepoNum        int
	Verbose        bool
	Find           bool // 是否找到name 作者
}

// var NumberWeek = map[int]string{
// 	0: "Sun",
// 	1: "Mon",
// 	2:"Fer"
// }

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
		fmt.Printf(" ")
	case commitNum < 5:
		// fmt.Printf(" ")
		color.New(color.BgHiGreen).Add(color.BgHiGreen).Print("+", commitNum)
		fmt.Printf(" ")

	case commitNum < 10:
		color.New(color.BgGreen).Add(color.BgWhite).Print("+", commitNum)
		fmt.Printf(" ")
	default:
		color.New(color.BgGreen, color.FgHiWhite).Print(commitNum) // 深绿
		fmt.Printf(" ")
		// color.New(color.BgGreen, color.FgWhite).Print("██") // 最深绿
	}
}

func detailShow(stat *Stats, start time.Time) {
	fmt.Printf("分析报告  (%s to %s)\n", formatTime(start), formatTime(time.Now()))
	fmt.Printf("============================\n")
	fmt.Printf("共扫描%d个仓库\n", stat.RepoNum)
	fmt.Printf("共计提交%d次\n", stat.TotalComitNum)

	fmt.Printf("总计增加行数: +%d\n", stat.TotalAdditions)
	fmt.Printf("总计删除行数: -%d\n", stat.TotalDeletions)

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
	fmt.Printf("  More \n")
}

func uiShow(stat *Stats, offset int, daysTotal int, start time.Time, name string, metric SortType) {
	switch metric {
	case Default:
		printCol(stat, offset, daysTotal, start, name) // 打印贡献图
	case CommitNum, LineNum, DeletionNum:
		printContribution(stat, metric) // 打印贡献排名
	default:
		fmt.Println("未知的显示模式。")
	}
}

func printContribution(stat *Stats, metric SortType) {
	if stat == nil || len(stat.Authors) == 0 {
		fmt.Println("没有找到作者贡献数据。")
		return
	}

	authorStats := make([]*AuthorStat, 0, len(stat.Authors))
	for _, authorData := range stat.Authors {
		authorStats = append(authorStats, authorData)
	}

	metricName := "提交次数" // 默认
	sort.Slice(authorStats, func(i, j int) bool {
		switch metric {
		case LineNum:
			metricName = "增加行数"
			return authorStats[i].TotalAdditions > authorStats[j].TotalAdditions
		case DeletionNum:
			metricName = "删除行数"
			return authorStats[i].TotalDeletions > authorStats[j].TotalDeletions
		case CommitNum:
			fallthrough
		default:
			metricName = "提交次数"
			return authorStats[i].Commits > authorStats[j].Commits
		}
	})

	fmt.Println()
	color.Cyan("作者贡献排行榜 (按 %s 排序)", metricName)

	table := tablewriter.NewWriter(os.Stdout)
	header := []string{"排名", "作者", "提交次数", "增加行数(+)", "删除行数(-)"}
	table.Header(header)

	limit := 15
	if len(authorStats) < limit {
		limit = len(authorStats)
	}

	for i := 0; i < limit; i++ {
		s := authorStats[i]
		row := []string{
			strconv.Itoa(i + 1),
			s.Name,
			strconv.Itoa(s.Commits),
			color.GreenString(strconv.Itoa(s.TotalAdditions)),
			color.RedString(strconv.Itoa(s.TotalDeletions)),
		}
		table.Append(row)
	}

	table.Render()
}

func printCol(stat *Stats, offset int, daysTotal int, start time.Time, name string) {
	if len(name) > 0 && stat.TotalComitNum == 0 {
		fmt.Printf("作者%s 在仓库中尚未有提交\n", name)
		return
	}
	now := time.Now()
	PrintMonths(now, offset)
	for i := 0; i < 7; i++ {
		PrintWeek(i)
		// tempStart := start.AddDate(0, 0, i)
		// month := tempStart.Month()
		for j := daysTotal - i; j >= 0; j -= 7 {
			// if tempStart.Month() != month {
			// 	fmt.Printf()
			// }
			printDay(stat.CommitsByDay[j])

		}
		fmt.Print("\n")
		fmt.Printf("\n")
	}
	printLegend()
	if stat.Verbose {
		detailShow(stat, start)
	}
}

func GenerateStats(repositories []string, month int, verbose bool, name string, metric SortType) {
	stat := &Stats{
		CommitsByDay:   make(map[int]int),
		Authors:        make(map[string]*AuthorStat),
		TotalAdditions: 0,
		TotalDeletions: 0,
		RepoNum:        0,
		Verbose:        verbose,
	}
	s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
	s.Suffix = "正在分析git仓库中...\n"
	s.Start()
	start, days := startDayRecentMonths(time.Now(), month)
	fillMap(repositories, stat, start, name)
	uiShow(stat, month, days, start, name, metric)

	s.Stop()
	// PrintContribGraph(commits, days)

}

func fillMap(repositories []string, stat *Stats, since time.Time, name string) {
	// commitsByDay := make(map[int]int)
	// commitsByAuthor := make(map[string]int)
	if stat.Authors == nil {
		stat.Authors = make(map[string]*AuthorStat)
	}
	for _, repo := range repositories {
		stat.RepoNum++

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
			authorName := c.Author.Name

			if len(name) > 0 && name != authorName {
				return nil
			}

			if _, ok := stat.Authors[authorName]; !ok {
				stat.Authors[authorName] = &AuthorStat{Name: authorName}
			}

			stat.Authors[authorName].Commits++
			stat.CommitsByDay[calculateDaysSince(c.Author.When)]++
			stat.TotalComitNum++

			commitStats, err := c.Stats()
			if err != nil {
				return nil
			}
			for _, commitStat := range commitStats {
				addition := commitStat.Addition
				deletion := commitStat.Deletion

				stat.Authors[authorName].TotalAdditions += addition
				stat.Authors[authorName].TotalDeletions += deletion

				stat.TotalAdditions += addition
				stat.TotalDeletions += deletion
			}
			return nil
		})
	}
}
