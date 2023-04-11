package commands

import (
	"fmt"
	"github.com/gookit/color"
	"github.com/goravel/framework/contracts/console"
	"github.com/goravel/framework/contracts/console/command"
	"github.com/goravel/framework/facades"
	"github.com/spf13/cast"
	"os"
)

type HashInsert struct {
}

// Signature The name and signature of the console command.
func (receiver *HashInsert) Signature() string {
	return "hash:insert"
}

// Description The console command description.
func (receiver *HashInsert) Description() string {
	return "将给定路径的 MD5 文件插入到数据库中"
}

// Extend The console command extend.
func (receiver *HashInsert) Extend() command.Extend {
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
				Name:    "dir",
				Value:   "hash",
				Aliases: []string{"d"},
				Usage:   "MD5 文件存放目录",
			},
		},
	}
}

// Handle Execute the console command.
func (receiver *HashInsert) Handle(ctx console.Context) error {
	table := cast.ToInt(ctx.Option("table"))
	dir := ctx.Option("dir")

	color.Warnf("分表数量: %d\n", table)
	color.Warnf("存放目录: %s\n", dir)

	for i := 1; i <= table; i++ {
		sql := fmt.Sprintf(`DROP TABLE IF EXISTS qq_%d;`, i)
		_, err := facades.Orm.Connection("hash").Query().Exec(sql)
		if err != nil {
			panic(err)
		}
		sql = fmt.Sprintf(`CREATE TABLE qq_%d (hash CHAR(32) NOT NULL, qq BIGINT NOT NULL);`, i)
		color.Greenf("正在创建表: %d\n", i)
		_, err = facades.Orm.Connection("hash").Query().Exec(sql)
		if err != nil {
			panic(err)
		}
	}

	color.Greenf("建表完成\n")
	color.Warnf("正在导入数据\n")

	for i := 1; i <= table; i++ {
		_, err := facades.Orm.Connection("hash").Query().Exec(fmt.Sprintf(`COPY qq_%d FROM '%s/%d.csv' WITH DELIMITER ',';`, i, dir, i))
		if err != nil {
			panic(err)
		}
		color.Greenf("导入完成: %d\n", i)
		// 删除文件
		_ = os.Remove(fmt.Sprintf("%s/%d.csv", dir, i))
	}

	color.Warnf("导入完成\n")
	color.Warnf("正在创建索引\n")

	for i := 1; i <= table; i++ {

		_, err := facades.Orm.Connection("hash").Query().Exec(fmt.Sprintf(`CREATE INDEX hash_%d ON qq_%d USING hash (hash);`, i, i))
		if err != nil {
			panic(err)
		}
	}

	color.Warnf("索引创建完成\n")
	color.Warnf("运行结束\n")

	return nil
}
