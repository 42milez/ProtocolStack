# -*- mode: ruby -*-
# vi: set ft=ruby :

Vagrant.configure("2") do |config|
  config.vm.box = "generic/fedora33"
  config.vm.define "ps.vagrant" do |instance|
    instance.vm.provider :virtualbox do |vb|
      vb.cpus = 4
      vb.customize ["modifyvm", :id, "--natdnshostresolver1", "on"]
      vb.customize ["modifyvm", :id, "--natdnsproxy1", "on"]
      vb.gui = false
      vb.memory = 8192
      vb.name = "ps.vagrant"
    end
  end
  config.vm.hostname = "ps.vagrant"
  config.vm.network :private_network, ip: "192.168.33.10"
  config.vm.provision "shell", path: "./vagrant/provisioners/root.sh"
  config.vm.provision "shell", path: "./vagrant/provisioners/vagrant.sh", privileged: false
end
