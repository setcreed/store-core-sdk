package examples

import (
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"

	"github.com/setcreed/store-core-sdk/pkg/builder"
)

func QueryTest() {
	// 客户端构建器
	c, _ := builder.NewClientBuilder("localhost:8080").WithOption(grpc.WithInsecure()).Build()
	// 参数构建器
	paramBuilder := builder.NewParamBuilder().Add("id", 1)
	// api 构建器
	api := builder.NewApiBuilder("userList", builder.APITYPE_QUERY)

	// 查询结果集
	users := make([]*User, 0)

	//执行 和调用  API
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()
	err := api.Invoke(ctx, paramBuilder, c, &users)
	if err != nil {
		log.Fatal(err)
	}
	for _, dept := range users {
		fmt.Println(dept)
	}
}
