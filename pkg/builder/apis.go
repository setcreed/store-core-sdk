package builder

import (
	"context"
	"fmt"

	"github.com/mitchellh/mapstructure"
	"github.com/setcreed/store-core-sdk/pkg/utils"
	v1 "github.com/setcreed/store-core/api/store_service/v1"
)

const (
	APITYPE_QUERY = iota
	APITYPE_EXEC
)

type ApiBuilder struct {
	name    string //api 名称
	apiType int
}

func NewApiBuilder(name string, apiType int) *ApiBuilder {
	return &ApiBuilder{name: name, apiType: apiType}
}

// 普通执行， 不是事务
func (ab *ApiBuilder) Invoke(ctx context.Context, paramBuilder *ParamBuilder, client v1.DBServiceClient, out interface{}) error {
	if ab.apiType == APITYPE_QUERY { //查询
		req := &v1.QueryRequest{Name: ab.name, Params: paramBuilder.Build()}
		rsp, err := client.Query(ctx, req)
		if err != nil {
			fmt.Println(err)
			return err
		}
		for _, item := range rsp.GetResult() {
			fmt.Println(item.AsMap())
		}
		if out != nil { //如果 out 没有传值 不做转换
			return mapstructure.Decode(utils.PbStructsToMapList(rsp.GetResult()), out)
		}
		return nil

	} else {
		req := &v1.ExecRequest{Name: ab.name, Params: paramBuilder.Build()}
		rsp, err := client.Exec(ctx, req)
		if err != nil {
			return err
		}
		if out != nil { //如果 out 没有传值 不做转换
			var m map[string]interface{}
			if rsp.Select != nil {
				m = rsp.Select.AsMap()
				m["_RowsAffected"] = rsp.RowsAffected
			} else {
				m = map[string]interface{}{"_RowsAffected": rsp.RowsAffected}
			}
			return mapstructure.Decode(m, out)
		}
		return nil
	}
}
