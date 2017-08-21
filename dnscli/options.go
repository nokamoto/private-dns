package main

import (
	"google.golang.org/grpc"
	"fmt"
	"golang.org/x/net/context"
)

type Options struct {
	Host *string
	Port *int
	ApiKey *string
}

func (o Options)GetRequestMetadata(_ context.Context, _ ...string) (map[string]string, error) {
	return map[string]string{"x-api-key": *o.ApiKey}, nil
}

func (o Options)RequireTransportSecurity() bool {
	return false
}

func (o Options)Dial() (*grpc.ClientConn, error) {
	opts := []grpc.DialOption{grpc.WithInsecure()}

	opts = append(opts, grpc.WithPerRPCCredentials(o))

	return grpc.Dial(fmt.Sprintf("%s:%d", *o.Host, *o.Port), opts...)
}
