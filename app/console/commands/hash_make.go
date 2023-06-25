package commands

import (
	"bufio"
	"fmt"
	"os"

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
			&command.StringFlag{
				Name:    "sum",
				Value:   "4000000000",
				Aliases: []string{"s"},
				Usage:   "生成的QQ号最大值",
			},
		},
	}
}

// Handle Execute the console command.
func (receiver *HashMake) Handle(ctx console.Context) error {
	// 要生成 MD5 值的 QQ 号的范围
	start := 10000
	end := cast.ToInt(ctx.Option("sum"))

	color.Warnf("号最大值: %d\n", end)
	color.Warnf("开始生成\n")

	// 生成 MD5 值并写入数据库
	file, err := os.OpenFile("hash.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	writer := bufio.NewWriterSize(file, 4096*256)

	for num := start; num <= end; num++ {
		if writer.Available() <= 4096 || num == end {
			err = writer.Flush()
			if err != nil {
				panic(err)
			}
			color.Greenf("当前: %d - %d\n", num, end)
		}

		// 写入 MD5 值到对应的 writer
		md5Sum := helpers.MD5(fmt.Sprintf("%d@qq.com", num))
		_, err = fmt.Fprintf(writer, "%s,%d\n", md5Sum, num)
		if err != nil {
			panic(err)
		}
	}

	color.Greenln("生成完成")
	return nil
}
