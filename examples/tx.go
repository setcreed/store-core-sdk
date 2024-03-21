package examples

import (
	"context"
	"fmt"
	"github.com/setcreed/store-core-sdk/pkg/builder"
	"google.golang.org/grpc"
	"log"
)

//- name: addUser
//sql: "insert into users(user_name,user_pass) values(@userName,@userPass)"
//select:
//sql: "SELECT LAST_INSERT_ID() as user_id"
//- name: addUserOrder
//sql: "insert into user_order(user_id) values(@user_id)"

type ExecResult struct {
	RowsAffected int `mapstructure:"_RowsAffected"`
}

func TxTest() {
	client, cerr := builder.NewClientBuilder("localhost:8080").WithOption(grpc.WithInsecure()).Build()
	if cerr != nil {
		log.Fatal(cerr)
	}
	//创建 事务API
	txApi := builder.NewTxApi(context.Background(), client)
	// 模仿了 Gorm
	err := txApi.Tx(func(tx *builder.TxApi) error {
		//构建 用户实体
		user := &User{UserName: "zhangsan", UserPassword: "123"}
		addUserParam := builder.NewParamBuilder().
			Add("user_name", user.UserName).
			Add("user_password", user.UserPassword)

		//执行新增 用户，user.UserID 会自动赋值
		err := tx.Exec("addUser", addUserParam, user)
		if err != nil {
			return err
		}

		log.Println("新增用户成功，用户ID是:", user.UserId)

		//给用户下订单
		addUserOrderParam := builder.NewParamBuilder().
			Add("user_id", user.UserId)

		execRet := &ExecResult{}
		err = tx.Exec("addUserOrder", addUserOrderParam, execRet)
		if err != nil || execRet.RowsAffected != 1 {
			return fmt.Errorf("addUserOrder error")
		}
		return nil
	})
	log.Println(err)
	select {}
}
