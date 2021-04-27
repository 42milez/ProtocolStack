# -*- mode: ruby -*-
# vi: set ft=ruby :

BOX = "generic/fedora33"
CPU = 4
HOSTNAME = "ps.vagrant"
MEM = 8192
STATIC_IP = "192.168.33.10"

Vagrant.configure("2") do |config|
  config.vm.box = BOX
  config.vm.define HOSTNAME do |instance|
    instance.vm.provider :virtualbox do |vb|
      vb.cpus = CPU
      vb.customize ["modifyvm", :id, "--natdnshostresolver1", "on"]
      vb.customize ["modifyvm", :id, "--natdnsproxy1", "on"]
      vb.gui = false
      vb.memory = MEM
      vb.name = HOSTNAME
    end
  end
  config.vm.hostname = HOSTNAME
  config.vm.network :forwarded_port, guest: 2345, host: 2345
  config.vm.network :private_network, ip: STATIC_IP
  config.vm.provision "shell", path: "./vagrant/provisioners/root.sh"
  config.vm.provision "shell", path: "./vagrant/provisioners/vagrant.sh", privileged: false
end
