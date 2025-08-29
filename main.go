package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"main.go/scan"
	"main.go/stats"
)

var (
	path    string
	offset  int
	verbose bool
)

type Config struct {
	User string
}

var config Config

// viper 加载配置文件 "config.yaml"
func initConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./")
	if err := viper.ReadInConfig(); err != nil {
		currentPath, _ := os.Getwd()
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Fatalf("找不到目录 %s 下的名为 config.yaml 的配置文件\n", currentPath)
		} else {
			log.Fatalf("配置文件加载错误: %v", err)
		}
	}
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("配置文件解析错误: %v", err)
	}

}

func preRunLogic(cmd *cobra.Command, args []string) {
	if path == "./" {
		path, _ = os.Getwd()
	} else if path == "../" {
		currentPath, _ := os.Getwd()
		path = filepath.Dir(currentPath)
	}
}

// 创建根名令
func createRootCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:              "./gitLocalLog  ",
		Short:            "Git 代码统计工具",
		Long:             `一个基于 go-git 的命令行工具，用于统计 Git 仓库的提交和代码活跃度。`,
		PersistentPreRun: preRunLogic,
		RunE: func(cmd *cobra.Command, args []string) error {
			// 在这里写主逻辑

			runLog(path, offset, verbose, "", 0)
			return nil
		},
	}

	return rootCmd

}

// 给rootCommand 添加参数
func addFlags(rootCmd *cobra.Command) {
	rootCmd.PersistentFlags().StringVarP(&path, "path", "p", "./", "用于搜索的文件目录")
	rootCmd.PersistentFlags().IntVarP(&offset, "month", "m", 6, "查询月份数量")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "启用详细输出")
}

// 添加子命令
func addSubCommands(rootCmd *cobra.Command) {
	authorCmd := createAuthorCommand()
	contributionCmd := createContributionCommand()

	rootCmd.AddCommand(authorCmd, contributionCmd)
}

// 创建Author 子命令
func createAuthorCommand() *cobra.Command {
	authorCmd := &cobra.Command{
		Use:   "author ",
		Short: "查询指定作者的状况",
		RunE: func(cmd *cobra.Command, args []string) error {

			name, _ := cmd.Flags().GetString("name")
			if len(name) == 0 {
				initConfig()
				name = config.User
				// fmt.Println("config :", name)
			}
			runLog(path, offset, verbose, name, stats.Default)
			return nil
		},
	}
	defaulyName := config.User
	authorCmd.Flags().StringP("name", "n", defaulyName, "作者名称")
	return authorCmd

}

// 创建Contribution 子命令
func createContributionCommand() *cobra.Command {
	contributionCmd := &cobra.Command{
		Use:   "contribution ",
		Short: "根据排名列出仓库所有作者",
		RunE: func(cmd *cobra.Command, args []string) error {

			sortType, _ := cmd.Flags().GetInt("sort")

			runLog(path, offset, verbose, "", stats.SortType(sortType))
			return nil
		},
	}

	contributionCmd.Flags().IntP("sort", "s", int(stats.CommitNum), "排序依据 (commits, lines)")
	return contributionCmd

}

// 主函数入口
func runLog(path string, monthOffset int, verbose bool, name string, metric stats.SortType) {
	repositories := scan.ScanPath(path)
	if len(repositories) == 0 {
		fmt.Println("此目录下不存在git仓库")
		return
	}
	stats.GenerateStats(repositories, monthOffset, verbose, name, metric)
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	rootCmd := createRootCommand()
	addFlags(rootCmd)
	addSubCommands(rootCmd)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
	}
}
