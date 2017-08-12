# -*- mode: ruby -*-
# vi: set ft=ruby :

Vagrant.configure("2") do |config|
  config.vm.box = "centos/7"

  config.vm.network "private_network", ip: "192.168.33.10"

  config.vm.provision "docker" do |d|
    d.run "private-dns",
      image: "nokamotohub/private-dns",
      args: "-p 53:53/tcp -p 53:53/udp -p 9999:9999/tcp --cap-add=NET_ADMIN"
  end
end
