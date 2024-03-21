package builder

import (
	"github.com/setcreed/store-core/api/store_service/v1"
	"google.golang.org/protobuf/types/known/structpb"
)

// 参数构建器
type ParamBuilder struct {
	param map[string]interface{}
}

func NewParamBuilder() *ParamBuilder {
	return &ParamBuilder{param: make(map[string]interface{})}
}
func (pb *ParamBuilder) Add(name string, value interface{}) *ParamBuilder {
	pb.param[name] = value
	return pb
}
func (pb *ParamBuilder) Build() *v1.SimpleParams {
	paramStruct, _ := structpb.NewStruct(pb.param)
	return &v1.SimpleParams{
		Params: paramStruct,
	}
}
