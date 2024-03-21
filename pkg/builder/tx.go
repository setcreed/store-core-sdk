package builder

import (
	"context"
	"github.com/mitchellh/mapstructure"
	v1 "github.com/setcreed/store-core/api/store_service/v1"
	"google.golang.org/grpc"
	"log"
)

// 事务API对象
type TxApi struct {
	ctx    context.Context
	cancel context.CancelFunc
	client v1.DBService_TxClient
}

func NewTxApi(ctx context.Context, client v1.DBServiceClient, opts ...grpc.CallOption) *TxApi {
	apiCtx, cancel := context.WithCancel(ctx)
	txClient, err := client.Tx(apiCtx, opts...)
	if err != nil {
		panic(err)
	}
	return &TxApi{ctx: ctx, client: txClient, cancel: cancel}
}
func (ta *TxApi) Exec(apiname string, paramBuilder *ParamBuilder, out interface{}) error {
	err := ta.client.Send(&v1.TxRequest{Api: apiname, Params: paramBuilder.Build(), Type: "exec"})
	//对于exec ,如果不出错， 会返回一个map，其中key=exec ,  值是一个interface切片 ，包含了两项
	//1、受影响的行  2 、selectkey(如果有的话)
	if err != nil {
		return err
	}
	rsp, err := ta.client.Recv() //接收消息
	if err != nil {
		return err
	}
	if out != nil {
		if execRet, ok := rsp.Result.AsMap()["exec"]; ok { //execRet 是一个 切片 []interface{受影响的行，select}  .select 可能是nil
			if execRet.([]interface{})[1] != nil { //代表select 有值
				m := execRet.([]interface{})[1].(map[string]interface{})
				m["_RowsAffected"] = execRet.([]interface{})[0]
				return mapstructure.Decode(m, out)
			} else { //没有select 情况 直接塞一个_RowsAffected 返回
				m := map[string]interface{}{"_RowsAffected": execRet.([]interface{})[0]}
				return mapstructure.Decode(m, out)
			}

		}
	}
	return nil

}
func (ta *TxApi) Query(apiName string, paramBuilder *ParamBuilder, out interface{}) error {
	err := ta.client.Send(&v1.TxRequest{Api: apiName, Params: paramBuilder.Build(), Type: "query"})
	// 对于查询，如果不出错，会返回一个map   其中key=query    值是查询结果
	if err != nil {
		return err
	}
	rsp, err := ta.client.Recv()
	if err != nil {
		return err
	}
	if out != nil {
		if queryRet, ok := rsp.Result.AsMap()["query"]; ok {
			return mapstructure.Decode(queryRet, out)
		}
	}
	return nil
}

// 模仿 gorm
func (ta *TxApi) Tx(fn func(tx *TxApi) error) error {
	err := fn(ta)
	if err != nil {
		log.Println("tx error:", err)
		ta.cancel() //取消
		return err
	}
	return ta.client.CloseSend() //协程不安全
}
