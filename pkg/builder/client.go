package builder

import (
	"context"

	v1 "github.com/setcreed/store-core/api/store_service/v1"
	"google.golang.org/grpc"
)

type ClientBuilder struct {
	url  string
	opts []grpc.DialOption
}

func NewClientBuilder(url string) *ClientBuilder {
	return &ClientBuilder{url: url}
}
func (cb *ClientBuilder) WithOption(opts ...grpc.DialOption) *ClientBuilder {
	cb.opts = append(cb.opts, opts...)
	return cb
}
func (cb *ClientBuilder) Build() (v1.DBServiceClient, error) {
	client, err := grpc.DialContext(context.Background(),
		cb.url,
		cb.opts...,
	)
	if err != nil {
		return nil, err
	}
	return v1.NewDBServiceClient(client), nil
}
