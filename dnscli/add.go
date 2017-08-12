package main

import (
	"fmt"
	"errors"
	pb "github.com/nokamoto/private-dns/proto"
	"golang.org/x/net/context"
)

type Add struct {}

func (c Add)Name() []string {
	return []string{"add", "a"}
}

func (c Add)Run(opts Options, args []string) error {
	if len(args) != 2 {
		return errors.New("Usage: add ip hostname")
	}

	ip, err := validateIp(args[0])
	if err != nil {
		return err
	}

	host, err := validateHost(args[1])
	if err != nil {
		return err
	}

	con, err := opts.Dial()
	if err != nil {
		return err
	}

	defer con.Close()

	client := pb.NewDnsServiceClient(con)
	req := pb.Host{Hostname: host, Ip: ip}

	fmt.Printf("add: ip=%s hostname=%s\n", ip, host)

	_, err = client.Add(context.Background(), &req)
	if err != nil {
		return err
	}

	return nil
}
