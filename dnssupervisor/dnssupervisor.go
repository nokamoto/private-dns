package main

import (
	"golang.org/x/net/context"
	"log"
	google_protobuf "github.com/golang/protobuf/ptypes/empty"
	pb "github.com/nokamoto/private-dns/proto"
	"net"
	"fmt"
	"google.golang.org/grpc"
	"sort"
	"sync"
	"reflect"
)

type DnsSupervisor struct {
	hosts pb.HostList
	mux sync.Mutex
}

func sortHosts(hosts *pb.HostList) {
	sort.Slice(hosts.Hosts, func(i, j int) bool {
		si := hosts.Hosts[i]
		sj := hosts.Hosts[j]
		if si.Hostname < sj.Hostname {
			return true
		} else if si.Hostname == sj.Hostname {
			return si.Ip < sj.Ip
		}
		return false
	})
}

func (s *DnsSupervisor)find(host *pb.Host) int {
	for i, h := range s.hosts.Hosts {
		if reflect.DeepEqual(h, host) {
			return i
		}
	}
	return -1
}

func (s *DnsSupervisor)Add(_ context.Context, host *pb.Host) (*google_protobuf.Empty, error) {
	log.Printf("add: %v", *host)

	s.mux.Lock()

	exists := s.find(host) != -1

	if !exists {
		s.hosts.Hosts = append(s.hosts.Hosts, host)

		sortHosts(&s.hosts)
	}

	s.mux.Unlock()

	return &google_protobuf.Empty{}, nil
}

func (s *DnsSupervisor)Remove(_ context.Context, host *pb.Host) (*google_protobuf.Empty, error) {
	log.Printf("remove: %v", *host)

	s.mux.Lock()

	if found := s.find(host); found != -1 {
		s.hosts.Hosts = append(s.hosts.Hosts[:found], s.hosts.Hosts[found + 1:]...)
	}

	s.mux.Unlock()

	return &google_protobuf.Empty{}, nil
}

func (s *DnsSupervisor)Get(context.Context, *google_protobuf.Empty) (*pb.HostList, error) {
	log.Println("get:")

	s.mux.Lock()

	hosts := s.hosts

	s.mux.Unlock()

	return &hosts, nil
}

func (s *DnsSupervisor)Run(port int) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	} else {
		log.Printf("listen %v port", port)
	}

	opts := []grpc.ServerOption{}

	log.Println("register dns service")
	ns := grpc.NewServer(opts...)
	pb.RegisterDnsServiceServer(ns, s)

	log.Println("start gRPC...")
	ns.Serve(lis)
}
