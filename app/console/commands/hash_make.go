package commands

import (
	"fmt"

	"github.com/gookit/color"
	"github.com/goravel/framework/contracts/console"
	"github.com/goravel/framework/contracts/console/command"
	"github.com/goravel/framework/facades"
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

	_, err := facades.Orm.Connection("hash").Query().Exec(`DROP TABLE IF EXISTS qq_mails;`)
	if err != nil {
		panic(err)
	}
	_, err = facades.Orm.Connection("hash").Query().Exec(`CREATE TABLE qq_mails (hash CHAR(32) NOT NULL, qq BIGINT NOT NULL, PRIMARY KEY ( hash ) CLUSTERED);`)
	if err != nil {
		panic(err)
	}

	color.Greenf("建表完成\n")
	color.Warnf("正在导入数据\n")

	// 生成 MD5 值并写入数据库
	for num := start; num <= end; num++ {
		md5Sum := helpers.MD5(fmt.Sprintf("%d@qq.com", num))
		_, insertErr := facades.Orm.Connection("hash").Query().Exec(fmt.Sprintf(`INSERT INTO qq_mails (hash, qq) VALUES ('%s', '%d');`, md5Sum, num))
		if insertErr != nil {
			panic(insertErr)
		}
	}

	color.Greenln("生成完成")
	return nil
}
