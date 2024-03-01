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

	"weavatar/pkg/helper"
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
			&command.StringFlag{
				Name:    "number",
				Value:   "500",
				Aliases: []string{"n"},
				Usage:   "分表数量",
			},
			&command.StringFlag{
				Name:    "sum",
				Value:   "5000000000",
				Aliases: []string{"s"},
				Usage:   "生成的QQ号最大值",
			},
			&command.StringFlag{
				Name:    "dir",
				Value:   "hash",
				Aliases: []string{"d"},
				Usage:   "生成的文件存放目录",
			},
			&command.StringFlag{
				Name:    "type",
				Value:   "md5",
				Aliases: []string{"t"},
				Usage:   "生成的哈希类型 md5, sha256",
			},
		},
	}
}

// Handle Execute the console command.
func (receiver *HashMake) Handle(ctx console.Context) error {
	start := 10000
	end := cast.ToInt(ctx.Option("sum"))
	number := cast.ToInt(ctx.Option("number"))
	dir := ctx.Option("dir")
	hashType := ctx.Option("type")

	color.Warnf("分表数量: %d\n", number)
	color.Warnf("号最大值: %d\n", end)
	color.Warnf("存放目录: %s\n", dir)
	color.Warnf("哈希类型: %s\n\n", hashType)

	// 创建目录
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		panic(err)
	}

	fileWriters := make(map[int64]*bufio.Writer)
	for j := 1; j <= number; j++ {
		fileName := fmt.Sprintf("%s/%d.csv", dir, j)
		file, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			panic(err)
		}

		writer := bufio.NewWriterSize(file, 4096*256)
		fileWriters[int64(j)] = writer
	}

	// 生成 MD5 值并写入对应的文件
	for num := start; num <= end; num++ {
		var sum string
		if hashType == "sha256" {
			sum = helper.SHA256(fmt.Sprintf("%d@qq.com", num))
		} else {
			sum = helper.MD5(fmt.Sprintf("%d@qq.com", num))
		}
		hashIndex, hashErr := strconv.ParseInt(sum[:10], 16, 64)
		if hashErr != nil {
			return hashErr
		}
		tableIndex := (hashIndex % int64(number)) + 1

		writer := fileWriters[tableIndex]

		if writer.Available() <= 4096 || num == end {
			err = writer.Flush()
			if err != nil {
				panic(err)
			}
			color.Greenf("表 %d: %d - %d\n", tableIndex, num, end)
		}

		// 写入 MD5 值到对应的 writer
		_, err = fmt.Fprintf(writer, "%s,%d\n", sum, num)
		if err != nil {
			panic(err)
		}
	}

	color.Greenln("生成完成")
	return nil
}
