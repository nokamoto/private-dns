package main

import (
	"fmt"
	"errors"
	pb "github.com/nokamoto/private-dns/proto"
	"golang.org/x/net/context"
	"github.com/golang/protobuf/ptypes/empty"
)

type Get struct {}

func (c Get)Name() []string {
	return []string{"get", "g"}
}

func (c Get)Run(opts Options, args []string) error {
	if len(args) != 0 {
		return errors.New("Usage: get")
	}

	con, err := opts.Dial()
	if err != nil {
		return err
	}

	defer con.Close()

	client := pb.NewDnsServiceClient(con)
	req := empty.Empty{}

	res, err := client.Get(context.Background(), &req)
	if err != nil {
		return err
	}

	for _, host := range res.Hosts {
		fmt.Printf("%s %s\n", host.Ip, host.Hostname)
	}

	return nil
}
