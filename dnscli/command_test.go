package main

import "testing"

func Test_ValidateIp(t *testing.T) {
	validIps := []string{"127.0.0.1", "0.0.0.0", "255.255.255.255"}

	invalidIps := []string{"256.0.0.0", "01.0.0.0", "0.0.0", "a.0.0.0", "0.0.0.0a"}

	for _, ip := range validIps {
		if _, err := validateIp(ip); err != nil {
			t.Errorf("%s expected valid ip but actual %s", ip, err.Error())
		}
	}

	for _, ip := range invalidIps {
		if _, err := validateIp(ip); err == nil {
			t.Errorf("%s expected invalid ip but actual valid", ip)
		}
	}
}

func Test_ValidateHost(t *testing.T) {
	validHosts := []string{"localhost", "example.com", "example.com-test", "example.com_test", "EXAMPLE.COM"}

	invalidHosts := []string{"0.0.0.0", "localhost?"}

	for _, host := range validHosts {
		if _, err := validateHost(host); err != nil {
			t.Errorf("%s expected valid host but actual %s", host, err.Error())
		}
	}

	for _, host := range invalidHosts {
		if _, err := validateHost(host); err == nil {
			t.Errorf("%s expected invalid host but actual valid", host)
		}
	}
}

func Test_AllCommands(t *testing.T) {
	commands := allCommands()
	names := []string{}

	t.Log("all commands should have unique names")
	for _, c := range commands {
		for _, n := range c.Name() {
			for _, x := range names {
				if x == n {
					t.Errorf("duplicated subcommand: %s", x)
				}
			}
			names = append(names, n)
		}
	}
}
