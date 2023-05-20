package commands

import (
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
	return "建立哈希数据表"
}

// Extend The console command extend.
func (receiver *HashInsert) Extend() command.Extend {
	return command.Extend{}
}

// Handle Execute the console command.
func (receiver *HashInsert) Handle(ctx console.Context) error {

	_, err := facades.Orm.Connection("hash").Query().Exec(`DROP TABLE IF EXISTS qq_mails;`)
	if err != nil {
		panic(err)
	}
	_, err = facades.Orm.Connection("hash").Query().Exec(`CREATE TABLE qq_mails (hash CHAR(32) NOT NULL, qq BIGINT NOT NULL, PRIMARY KEY ( hash ) CLUSTERED);`)
	if err != nil {
		panic(err)
	}
	color.Greenf("建表完成\n")

	color.Warnf("请使用TiDB Lightning完成导入操作，预留至少800G磁盘空间\n")

	return nil
}
