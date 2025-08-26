package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

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

func main() {

	rootCmd := createRootCommand()
	addFlags(rootCmd)
	addSubCommands(rootCmd)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
	}

}

type Config struct {
	User string
}

var config Config

func initConfig() {
	viper.SetConfigName("config")
	viper.SetConfigFile("yaml")
	viper.AddConfigPath("./")
	if err := viper.ReadInConfig(); err != nil {
		currentPath, _ := os.Getwd()
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Fatal("找不到目录%s下的名为\"config.yaml\"的配置文件\n", currentPath)
		} else {
			log.Fatal("配置文件加载错误")
		}
	}
}

func createRootCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "gitLocalLog  [path]",
		Short: "Git 代码统计工具",
		Long:  `一个基于 go-git 的命令行工具，用于统计 Git 仓库的提交和代码活跃度。`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// 在这里写主逻辑
			if len(args) == 0 && !cmd.Flags().Changed("help") {
				return cmd.Help()
			}
			path := args[0]
			switch path {
			case "./":
				path, _ = os.Getwd()
			case "../":
				currentPath, _ := os.Getwd()
				path = filepath.Dir(currentPath)
			default:
			}
			runLog(path, offset, verbose, "", 0)
			return nil
		},
	}

	return rootCmd

}

func runLog(path string, monthOffset int, verbose bool, name string, metric stats.SortType) {
	repositories := scan.ScanPath(path)
	if len(repositories) == 0 {
		fmt.Println("此目录下不存在git仓库")
		return
	}
	stats.GenerateStats(repositories, offset, verbose, name, metric)
}

// 给rootCommand 添加参数
func addFlags(rootCmd *cobra.Command) {
	// rootCmd.PersistentFlags().StringVarP(&path, "path", "p", "./", "用于搜索的文件目录")
	rootCmd.PersistentFlags().IntVarP(&offset, "month", "m", 6, "查询月份数量")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "启用详细输出")
}

// 添加子命令
func addSubCommands(rootCmd *cobra.Command) {
	authorCmd := createAuthorCommand()
	contributionCmd := createContributionCommand()

	rootCmd.AddCommand(authorCmd, contributionCmd)
}

func createAuthorCommand() *cobra.Command {
	authorCmd := &cobra.Command{
		Use:   "author [path]",
		Short: "查询指定作者的状况",
		Run: func(cmd *cobra.Command, args []string) {
			path := args[0]
			switch path {
			case "./":
				path, _ = os.Getwd()
			case "../":
				currentPath, _ := os.Getwd()
				path = filepath.Dir(currentPath)
			default:
			}
			name, _ := cmd.Flags().GetString("name")
			// fmt.Println("ss")
			runLog(path, offset, verbose, name, stats.Default)
		},
	}
	defaulyName := config.User
	authorCmd.Flags().StringP("name", "n", defaulyName, "排序依据 (commits, lines)")
	return authorCmd

}

func createContributionCommand() *cobra.Command {
	contributionCmd := &cobra.Command{
		Use:   "contribution [path]",
		Short: "根据排名列出仓库所有作者",
		Run: func(cmd *cobra.Command, args []string) {
			path := args[0]
			switch path {
			case "./":
				path, _ = os.Getwd()
			case "../":
				currentPath, _ := os.Getwd()
				path = filepath.Dir(currentPath)
			default:
			}

			sortType, _ := cmd.Flags().GetInt("sort")
			runLog(path, offset, verbose, "", stats.SortType(sortType))
		},
	}

	contributionCmd.Flags().IntP("sort", "s", int(stats.CommitNum), "排序依据 (commits, lines)")
	return contributionCmd

}
