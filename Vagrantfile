Vagrant.configure("2") do |config|
	
    # Specify the base box
    config.vm.box = "kaorimatz/ubuntu-16.04-amd64"

    config.vm.provision :shell, path: "Vagrant-provision.sh"
    config.vm.provision :shell, path: "Vagrant-startup.sh", run: "always", privileged: true

    config.vm.network "private_network", ip: "172.28.128.3"
    
    config.vm.provider :virtualbox do |vb|
        vb.name = "ivr-yatego"
    end
end
