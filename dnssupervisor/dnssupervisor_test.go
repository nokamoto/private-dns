package main

import (
	"testing"
	pb "github.com/nokamoto/private-dns/proto"
	"reflect"
	"io/ioutil"
	"os"
)

func expect(t *testing.T, s *DnsSupervisor, hosts ...*pb.Host) {
	hl := &pb.HostList{Hosts: hosts}
	sortHosts(hl.Hosts)

	if actual, err := s.hosts(); err != nil {
		t.Error(err.Error())
	} else {
		if !reflect.DeepEqual(actual, hl) {
			t.Errorf("expect %v but actual %v", hl, actual)
		}

		if res, err := s.Get(nil, nil); err != nil {
			t.Error(err)
		} else if !reflect.DeepEqual(res, hl) {
			t.Errorf("get: expect %v but actual %v", hl, res)
		}
	}
}

var hosts = []string{"my.test.host0", "my.test.host1", "my.test.hosts2"}

var ips = []string{"192.168.33.0", "192.168.33.1", "192.168.33.2"}

func newSupervisor(t *testing.T, f func(*DnsSupervisor)){
	s := new(DnsSupervisor)

	if temp, err := ioutil.TempFile("", "dnssupervisor"); err != nil {
		t.Error(err.Error())
	} else {
		s.hostsFile = temp.Name()

		defer os.Remove(temp.Name())

		f(s)
	}
}

func TestDnsSupervisor_Add(t *testing.T) {
	newSupervisor(t, func(s *DnsSupervisor){
		entry1 := &pb.Host{Hostname: hosts[0], Ip: ips[0]}
		t.Logf("Add new entry %v", entry1)
		s.Add(nil, entry1)
		expect(t, s, entry1)

		entry2 := &pb.Host{Hostname: hosts[0], Ip: ips[0]}
		t.Logf("Add duplicated entry %v", entry2)
		s.Add(nil, entry2)
		expect(t, s, entry1)

		entry3 := &pb.Host{Hostname: hosts[1], Ip: ips[1]}
		t.Logf("Add new entry %v", entry3)
		s.Add(nil, entry3)
		expect(t, s, entry1, entry3)

		entry4 := &pb.Host{Hostname: hosts[2], Ip: ips[0]}
		t.Logf("Add new entry %v", entry4)
		s.Add(nil, entry4)
		expect(t, s, entry1, entry3, entry4)

		entry5 := &pb.Host{Hostname: hosts[2], Ip: ips[2]}
		t.Logf("Add new entry %v", entry5)
		s.Add(nil, entry5)
		expect(t, s, entry1, entry3, entry4, entry5)
	})
}

func TestDnsSupervisor_Remove(t *testing.T) {
	newSupervisor(t, func(s *DnsSupervisor){
		entry1 := &pb.Host{Hostname: hosts[0], Ip: ips[0]}
		s.Add(nil, entry1)

		entry2 := &pb.Host{Hostname: hosts[0], Ip: ips[0]}
		s.Add(nil, entry2)

		entry3 := &pb.Host{Hostname: hosts[1], Ip: ips[1]}
		s.Add(nil, entry3)

		expect(t, s, entry1, entry3)

		t.Logf("Remove entry %v", entry1)
		s.Remove(nil, entry1)
		expect(t, s, entry3)

		t.Logf("Remove entry %v (already removed)", entry2)
		s.Remove(nil, entry2)
		expect(t, s, entry3)

		t.Logf("Remove entry %v", entry3)
		s.Remove(nil, entry3)
		expect(t, s)
	})
}
