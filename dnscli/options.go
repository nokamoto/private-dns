package main

import (
	"google.golang.org/grpc"
	"fmt"
)

type Options struct {
	Host *string
	Port *int
}

func (o Options)Dial() (*grpc.ClientConn, error) {
	opts := []grpc.DialOption{grpc.WithInsecure()}

	return grpc.Dial(fmt.Sprintf("%s:%d", *o.Host, *o.Port), opts...)
}
