# gitLocalLog

一个用 Go 写的小工具，用于统计本地 Git 仓库的提交活跃度。  

## 功能

- 扫描目录下所有 Git 仓库
- 按作者统计提交/代码活跃度
- 排序展示贡献情况
- 可通过配置文件设定默认作者

## 用法

```bash
# 统计当前目录的仓库
./gitLocalLog ./

# 指定作者
./gitLocalLog author ./ -n "Alice"
> autho不指定-n时会自动加载配置文件

# 按提交数排序所有作者
./gitLocalLog contribution ./ -s 0
