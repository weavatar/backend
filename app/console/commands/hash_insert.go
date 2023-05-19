package commands

import (
	"fmt"
	"github.com/gookit/color"
	"github.com/goravel/framework/contracts/console"
	"github.com/goravel/framework/contracts/console/command"
	"github.com/goravel/framework/facades"
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
				Name:    "file",
				Value:   "hash.csv",
				Aliases: []string{"d"},
				Usage:   "MD5 文件",
			},
		},
	}
}

// Handle Execute the console command.
func (receiver *HashInsert) Handle(ctx console.Context) error {
	file := ctx.Option("file")

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
	_, err = facades.Orm.Connection("hash").Query().Exec(fmt.Sprintf(`LOAD DATA INFILE '%s' INTO TABLE qq_mails FIELDS TERMINATED BY ',' LINES TERMINATED BY '\n';`, file))
	if err != nil {
		panic(err)
	}
	color.Greenf("导入完成\n")
	// 删除文件
	_ = os.Remove(file)

	color.Warnf("运行结束\n")

	return nil
}
