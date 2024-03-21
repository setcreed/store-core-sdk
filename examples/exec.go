package examples

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/grpc"

	"github.com/setcreed/store-core-sdk/pkg/builder"
	v1 "github.com/setcreed/store-core/api/store_service/v1"
)

type UserAddResult struct {
	UserID       int `mapstructure:"user_id"`
	RowsAffected int `mapstructure:"_RowsAffected"`
}

func ExecTest() {
	client, _ := grpc.DialContext(context.Background(),
		"localhost:8080",
		grpc.WithInsecure(),
	)
	c := v1.NewDBServiceClient(client)

	//构建 参数
	paramBuilder := builder.NewParamBuilder().
		Add("user_name", "wufeng").
		Add("user_password", "123456")
	//构建API对象
	api := builder.NewApiBuilder("addUser", builder.APITYPE_EXEC)

	ret := &UserAddResult{}
	err := api.Invoke(context.Background(), paramBuilder, c, ret)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(ret)
}
