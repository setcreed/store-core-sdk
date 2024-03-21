package utils

import "google.golang.org/protobuf/types/known/structpb"

func PbStructsToMapList(in []*structpb.Struct) []map[string]interface{} {
	list := make([]map[string]interface{}, len(in))
	for i, item := range in {
		list[i] = item.AsMap()
	}
	return list
}
