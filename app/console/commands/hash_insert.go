package commands

import (
	"fmt"
	"os"
	"strconv"

	"github.com/gookit/color"
	"github.com/goravel/framework/contracts/console"
	"github.com/goravel/framework/contracts/console/command"
	"github.com/goravel/framework/facades"
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
			&command.StringFlag{
				Name:    "table",
				Value:   "500",
				Aliases: []string{"t"},
				Usage:   "分表数量",
			},
			&command.StringFlag{
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
	table, err := strconv.Atoi(ctx.Option("table"))
	if err != nil {
		color.Errorln("分表数量必须是数字")
		return err
	}
	dir := ctx.Option("dir")

	color.Warnf("分表数量: %d\n", table)
	color.Warnf("存放目录: %s\n", dir)

	for i := 1; i <= table; i++ {
		sql := fmt.Sprintf(`DROP TABLE IF EXISTS qq_%d;`, i)
		_, err = facades.Orm().Connection("hash").Query().Exec(sql)
		if err != nil {
			panic(err)
		}
		sql = fmt.Sprintf(`CREATE TABLE qq_%d (hash CHAR(32) NOT NULL, qq BIGINT NOT NULL);`, i)
		color.Greenf("正在创建表: %d\n", i)
		_, err = facades.Orm().Connection("hash").Query().Exec(sql)
		if err != nil {
			panic(err)
		}
		sql = fmt.Sprintf(`ALTER TABLE qq_%d OWNER TO hash;`, i)
		color.Greenf("正在设置表所有者: %d\n", i)
		_, err = facades.Orm().Connection("hash").Query().Exec(sql)
		if err != nil {
			panic(err)
		}
	}

	color.Greenln("建表完成")
	color.Warnln("正在导入数据")

	for i := 1; i <= table; i++ {
		_, err = facades.Orm().Connection("hash").Query().Exec(fmt.Sprintf(`COPY qq_%d FROM '%s/%d.csv' WITH DELIMITER ',';`, i, dir, i))
		if err != nil {
			panic(err)
		}
		color.Greenf("导入完成: %d\n", i)
		// 删除文件
		_ = os.Remove(fmt.Sprintf("%s/%d.csv", dir, i))
	}

	color.Warnln("导入完成")
	color.Warnln("正在创建索引")

	for i := 1; i <= table; i++ {
		_, err = facades.Orm().Connection("hash").Query().Exec(fmt.Sprintf(`CREATE INDEX hash_%d ON qq_%d USING hash (hash);`, i, i))
		if err != nil {
			panic(err)
		}
		color.Greenf("索引创建完成: %d\n", i)
	}

	color.Warnln("索引创建完成")

	_, err = facades.Orm().Connection("hash").Query().Exec(`VACUUM FULL ANALYZE;`)
	if err != nil {
		panic(err)
	}

	color.Warnln("运行结束")

	return nil
}
