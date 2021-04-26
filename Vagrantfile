# -*- mode: ruby -*-
# vi: set ft=ruby :

Vagrant.configure("2") do |config|
  config.vagrant.plugins = ["vagrant-vbguest"]
  config.vm.box = "generic/fedora33"
  config.vm.hostname = "ps.vagrant"
  config.vm.network "private_network", ip: "192.168.33.10"
  config.vm.provider "virtualbox" do |vb|
    vb.gui = false
    vb.cpus = 4
    vb.memory = 8192
    vb.name = "ps.vagrant"
    vb.customize ["modifyvm", :id, "--natdnshostresolver1", "on"]
    vb.customize ["modifyvm", :id, "--natdnsproxy1", "on"]
  end
  config.vm.provision "shell", path: "./vagrant/provisioners/root.sh"
  config.vm.provision "shell", path: "./vagrant/provisioners/vagrant.sh", privileged: false
end
