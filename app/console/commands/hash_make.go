package commands

import (
	"bufio"
	"fmt"
	"os"
	"strconv"

	"github.com/gookit/color"
	"github.com/goravel/framework/contracts/console"
	"github.com/goravel/framework/contracts/console/command"
	"github.com/spf13/cast"

	"weavatar/packages/helpers"
)

type HashMake struct {
}

// Signature The name and signature of the console command.
func (receiver *HashMake) Signature() string {
	return "hash:make"
}

// Description The console command description.
func (receiver *HashMake) Description() string {
	return "生成 MD5值 对应的 QQ邮箱地址"
}

// Extend The console command extend.
func (receiver *HashMake) Extend() command.Extend {
	return command.Extend{
		Category: "hash",
		Flags: []command.Flag{
			{
				Name:    "table",
				Value:   "4000",
				Aliases: []string{"t"},
				Usage:   "分表数量",
			},
			{
				Name:    "sum",
				Value:   "4000000000",
				Aliases: []string{"s"},
				Usage:   "生成的QQ号最大值",
			},
			{
				Name:    "dir",
				Value:   "hash",
				Aliases: []string{"d"},
				Usage:   "生成的文件存放目录",
			},
		},
	}
}

// Handle Execute the console command.
func (receiver *HashMake) Handle(ctx console.Context) error {
	// 要生成 MD5 值的 QQ 号的范围
	start := 10000
	end := cast.ToInt(ctx.Option("sum"))
	// 分表数量
	table := cast.ToInt(ctx.Option("table"))

	color.Warnf("分表数量: %d\n", table)
	color.Warnf("号最大值: %d\n", end)
	color.Warnf("存放目录: %s\n\n", ctx.Option("dir"))

	// 创建目录
	err := os.MkdirAll(ctx.Option("dir"), 0755)
	if err != nil {
		panic(err)
	}

	fileWriters := make(map[int64]*bufio.Writer)
	for j := 1; j <= table; j++ {
		fileName := fmt.Sprintf("%s/%d.csv", ctx.Option("dir"), j)
		file, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			panic(err)
		}

		writer := bufio.NewWriterSize(file, 4096*256)
		fileWriters[int64(j)] = writer
	}

	// 生成 MD5 值并写入对应的文件
	for num := start; num <= end; num++ {
		md5Sum := helpers.MD5(fmt.Sprintf("%d@qq.com", num))
		hashIndex, hashErr := strconv.ParseInt(md5Sum[:10], 16, 64)
		if hashErr != nil {
			return hashErr
		}
		tableIndex := (hashIndex % int64(table)) + 1

		writer := fileWriters[tableIndex]

		if writer.Available() <= 4096 || num == end {
			err = writer.Flush()
			if err != nil {
				panic(err)
			}
			color.Greenf("表 %d: %d - %d\n", tableIndex, num, end)
		}

		// 写入 MD5 值到对应的 writer
		_, err = fmt.Fprintf(writer, "%s,%d\n", md5Sum, num)
		if err != nil {
			panic(err)
		}
	}

	color.Greenln("生成完成")
	return nil
}
