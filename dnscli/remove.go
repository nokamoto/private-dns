package main

import (
	"errors"
	pb "github.com/nokamoto/private-dns/proto"
	"golang.org/x/net/context"
	"fmt"
)

type Remove struct {}

func (c Remove)Name() []string {
	return []string{"remove", "r"}
}

func (c Remove)Run(opts Options, args []string) error {
	if len(args) != 2 {
		return errors.New("Usage: remove ip hostname")
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

	fmt.Printf("remove: ip=%s hostname=%s\n", ip, host)

	_, err = client.Remove(context.Background(), &req)
	if err != nil {
		return err
	}

	return nil
}
