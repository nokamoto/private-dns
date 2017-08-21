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
	"reflect"
	"io/ioutil"
	"strings"
	"os"
	"path"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type DnsSupervisor struct {
	hostsFile string
	apiKey string
}

func sortHosts(hosts []*pb.Host) {
	sort.Slice(hosts, func(i, j int) bool {
		si := hosts[i]
		sj := hosts[j]
		if si.Hostname < sj.Hostname {
			return true
		} else if si.Hostname == sj.Hostname {
			return si.Ip < sj.Ip
		}
		return false
	})
}

func find(hosts *pb.HostList, host *pb.Host) int {
	for i, h := range hosts.Hosts {
		if reflect.DeepEqual(h, host) {
			return i
		}
	}
	return -1
}

func (s *DnsSupervisor)hosts() (*pb.HostList, error) {
	log.Printf("load %s", s.hostsFile)

	bytes, err := ioutil.ReadFile(s.hostsFile)
	if err != nil {
		return nil, err
	}

	hl := new(pb.HostList)

	for i, line := range strings.Split(string(bytes), "\n") {
		xs := strings.Split(line, " ")

		if len(xs) == 1 && len(xs[0]) == 0 {
			continue
		} else if len(xs) != 2 {
			log.Printf("%s:line%d: illegal line %v (%d)", s.hostsFile, i, xs, len(xs))
		}

		if len(xs) == 2 {
			hl.Hosts = append(hl.Hosts, &pb.Host{Hostname: xs[1], Ip: xs[0]})
		}
	}

	log.Printf("%d host(s): %v", len(hl.Hosts), hl.Hosts)

	return hl, err
}

func atomicWrite(filename string, data []byte, perm os.FileMode) error {
	dir, name := path.Split(filename)
	f, err := ioutil.TempFile(dir, name)
	if err != nil {
		return err
	}
	log.Printf("write temporary file %s", f.Name())
	_, err = f.Write(data)
	if err == nil {
		err = f.Sync()
	}
	if closeErr := f.Close(); err == nil {
		err = closeErr
	}
	if permErr := os.Chmod(f.Name(), perm); err == nil {
		err = permErr
	}
	if err == nil {
		log.Printf("rename %s to %s", f.Name(), filename)
		err = os.Rename(f.Name(), filename)
	}
	// Any err should result in full cleanup.
	if err != nil {
		os.Remove(f.Name())
	}
	return err
}

func (s *DnsSupervisor)sync(hosts *pb.HostList) error {
	log.Printf("write %d host(s): %v", len(hosts.Hosts), hosts.Hosts)
	str := ""
	for _, h := range hosts.Hosts {
		str += fmt.Sprintf("%s %s\n", h.Ip, h.Hostname)
	}
	return atomicWrite(s.hostsFile, []byte(str), 0644)
}

func (s *DnsSupervisor)Add(_ context.Context, host *pb.Host) (*google_protobuf.Empty, error) {
	log.Printf("add: %v", *host)

	hl, err := s.hosts()
	if err != nil {
		return nil, err
	}

	exists := find(hl, host) != -1

	if !exists {
		hl.Hosts = append(hl.Hosts, host)

		sortHosts(hl.Hosts)

		err = s.sync(hl)
	}

	return &google_protobuf.Empty{}, err
}

func (s *DnsSupervisor)Remove(_ context.Context, host *pb.Host) (*google_protobuf.Empty, error) {
	log.Printf("remove: %v", *host)

	hl, err := s.hosts()
	if err != nil {
		return nil, err
	}

	if found := find(hl, host); found != -1 {
		hl.Hosts = append(hl.Hosts[:found], hl.Hosts[found + 1:]...)

		sortHosts(hl.Hosts)

		err = s.sync(hl)
	}

	return &google_protobuf.Empty{}, err
}

func (s *DnsSupervisor)Get(context.Context, *google_protobuf.Empty) (*pb.HostList, error) {
	log.Println("get:")

	return s.hosts()
}

func (s *DnsSupervisor)newInterceptor() grpc.UnaryServerInterceptor {
	f := func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Errorf(codes.Unauthenticated, "authentication required")
		}

		value, ok := md["x-api-key"]
		if !ok || len(value) != 1 {
			return nil, status.Errorf(codes.Unauthenticated, "authentication required")
		}

		if value[0] != s.apiKey {
			return nil, status.Errorf(codes.Unauthenticated, "authentication failed")
		}

		return handler(ctx, req)
	}
	return grpc.UnaryServerInterceptor(f)
}


func (s *DnsSupervisor)Run(port int) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	} else {
		log.Printf("listen %v port", port)
	}

	opts := []grpc.ServerOption{}

	opts = append(opts, grpc.UnaryInterceptor(s.newInterceptor()))

	log.Println("register dns service")
	ns := grpc.NewServer(opts...)
	pb.RegisterDnsServiceServer(ns, s)

	log.Println("start gRPC...")
	ns.Serve(lis)
}
