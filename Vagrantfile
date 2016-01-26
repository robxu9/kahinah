# -*- mode: ruby -*-
# vi: set ft=ruby :

# Vagrantfile API/syntax version. Don't touch unless you know what you're doing!
VAGRANTFILE_API_VERSION = "2"

# Inline script to help us provision the box.
$script = <<SCRIPT
PACKAGE="github.com/robxu9/kahinah"
echo "Provisioning Vagrant VM:"

echo "> updating centos packages"
sudo yum check-update
sudo yum upgrade -y

echo "> setting timezone"
sudo timedatectl set-timezone America/New_York

echo "> installing golang dependencies"
sudo yum install curl git mercurial make bison gcc glibc-devel yum-updateonboot vim -y

echo "> enabling yum-updateonboot"
sudo chkconfig --add yum-updateonboot

echo "> installing gvm"
bash < <(curl -s -S -L https://raw.githubusercontent.com/moovweb/gvm/master/binscripts/gvm-installer)
source ~/.gvm/scripts/gvm

echo "> create checkgolang script"
cat > /tmp/checkgolang << "EOF"
#!/bin/bash

set -x

echo "> sourcing gvm"
source ~/.gvm/scripts/gvm

echo "> finding latest go"
version="$(git ls-remote -t https://go.googlesource.com/go | awk -F/ '{ print $NF }' | grep go | grep -v beta | grep -v rc | tail -n1)"
echo ">> latest is $version"

echo "> checking our go"
ourgo=$(gvm list | grep "=>" | sed -e 's/=> //')
echo ">> ours is $ourgo"

if [[ "$ourgo" == "" ]]; then
    echo ">> we don't have any go installed by default"
    echo ">> bootstrapping go1.4.3"
    gvm install go1.4.3
fi

gvm use go1.4.3

if [[ "$version" != "$ourgo" ]]; then
    echo "> installing $version"
    gvm install $version
    gvm use $version --default
else
    echo "> looks good to me"
fi
EOF
sudo mv /tmp/checkgolang /usr/bin
sudo chmod +x /usr/bin/checkgolang


echo "> create systemd unit for checking golang version"
cat > /tmp/check-golang.service << "EOF"
[Unit]
Description=Check for the latest Golang version
Requires=network-online.target

[Service]
Type=oneshot
User=vagrant
ExecStart=/usr/bin/checkgolang

[Install]
WantedBy=multi-user.target
EOF
sudo mv -f /tmp/check-golang.service /etc/systemd/system
sudo systemctl enable check-golang.service

sudo systemctl start check-golang.service

echo "> link /vagrant to $PACKAGE"
cd /vagrant
source ~/.gvm/scripts/gvm
gvm linkthis $PACKAGE

echo "> open firewall"
sudo iptables -I INPUT -j ACCEPT
sudo iptables-save

echo "done."
SCRIPT

Vagrant.configure(VAGRANTFILE_API_VERSION) do |config|
  config.vm.box = "puppetlabs/centos-7.0-64-nocm"

  config.vm.network "private_network", ip: "192.168.33.10"

  # don't agent-forward
  # config.ssh.forward_agent = true

  # sync and nfs
  config.vm.synced_folder ".", "/vagrant-nfs", type: "nfs"

  # then use bind mounts for correct permissions
  config.bindfs.bind_folder "/vagrant-nfs", "/vagrant"
  
  # script to provision
  config.vm.provision "shell" do |s|
    s.inline = $script
    s.privileged = false
  end

end
