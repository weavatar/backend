package commands

import (
	"fmt"
	"os"
	"slices"
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
	return "将给定路径的哈希文件插入到数据库中"
}

// Extend The console command extend.
func (receiver *HashInsert) Extend() command.Extend {
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
				Name:    "dir",
				Value:   "hash",
				Aliases: []string{"d"},
				Usage:   "哈希文件存放目录",
			},
			&command.StringFlag{
				Name:    "type",
				Value:   "md5",
				Aliases: []string{"t"},
				Usage:   "哈希类型 md5, sha256",
			},
		},
	}
}

// Handle Execute the console command.
func (receiver *HashInsert) Handle(ctx console.Context) error {
	number, err := strconv.Atoi(ctx.Option("number"))
	if err != nil {
		color.Errorln("分表数量必须是数字")
		return err
	}
	dir := ctx.Option("dir")
	hashType := ctx.Option("type")
	hashSlices := []string{"md5", "sha256"}
	if !slices.Contains(hashSlices, hashType) {
		color.Errorln("哈希类型只能是 md5 或 sha256")
		return err
	}

	color.Warnf("分表数量: %d\n", number)
	color.Warnf("存放目录: %s\n", dir)
	color.Warnf("哈希类型: %s\n\n", hashType)

	for i := 1; i <= number; i++ {
		sql := fmt.Sprintf(`DROP TABLE IF EXISTS qq_%s_%d;`, hashType, i)
		_, err = facades.Orm().Connection("hash").Query().Exec(sql)
		if err != nil {
			panic(err)
		}
		sql = fmt.Sprintf(`CREATE TABLE qq_%s_%d (hash TEXT NOT NULL, qq BIGINT NOT NULL);`, hashType, i)
		color.Greenf("正在创建表: %d\n", i)
		_, err = facades.Orm().Connection("hash").Query().Exec(sql)
		if err != nil {
			panic(err)
		}
		sql = fmt.Sprintf(`ALTER TABLE qq_%s_%d OWNER TO hash;`, hashType, i)
		color.Greenf("正在设置表所有者: %d\n", i)
		_, err = facades.Orm().Connection("hash").Query().Exec(sql)
		if err != nil {
			panic(err)
		}
	}

	color.Greenln("建表完成")
	color.Warnln("正在导入数据")

	for i := 1; i <= number; i++ {
		_, err = facades.Orm().Connection("hash").Query().Exec(fmt.Sprintf(`COPY qq_%s_%d FROM '%s/%d.csv' WITH DELIMITER ',';`, hashType, i, dir, i))
		if err != nil {
			panic(err)
		}
		color.Greenf("导入完成: %d\n", i)
		// 删除文件
		_ = os.Remove(fmt.Sprintf("%s/%d.csv", dir, i))
	}

	color.Warnln("导入完成")
	color.Warnln("正在创建索引")

	for i := 1; i <= number; i++ {
		_, err = facades.Orm().Connection("hash").Query().Exec(fmt.Sprintf(`CREATE INDEX hash_%s_%d ON qq_%s_%d USING hash (hash);`, hashType, i, hashType, i))
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
