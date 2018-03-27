#!/bin/bash

echo "Provisioning VOIP virtual machine..."

#echo "deb http://archive.debian.org/debian/ squeeze main contrib non-free" > /etc/apt/sources.list


echo "Running apt-get"

apt-get update

apt-get install -y subversion
apt-get install -y autoconf
apt-get install -y build-essential
apt-get install -y sox
apt-get install -y postgresql
apt-get install -y postgresql-client
apt-get install -y openssl
apt-get install -y libh323-1.24.0v5
apt-get install -y libh323plus-dev
apt-get install -y speex
apt-get install -y libgsm-tools
apt-get install -y libpq5
apt-get install -y libpq-dev
apt-get install -y default-libmysqlclient-dev
apt-get install -y rsync
apt-get install -y vim
apt-get install -y doxygen

#echo "Setup postgres host access"
echo "listen_addresses = '*'" >> /etc/postgresql/9.6/main/postgresql.conf
echo "host    all         all         172.28.128.1/24           trust" >> /etc/postgresql/9.6/main/pg_hba.conf
service postgresql restart

echo "Downloading Yate"

cd /usr/src/
mkdir yate
cd yate
svn checkout http://voip.null.ro/svn/yate/tags/RELEASE_6_0_0
cd RELEASE_6_0_0

echo "Building Yate"

./autogen.sh
./configure --prefix=/opt/yate/
make
make install

echo "/opt/yate/lib" > /etc/ld.so.conf.d/yate.conf
ldconfig

mkdir /var/log/yate

ln -s /vagrant/cmd/inline/inline /opt/yate/share/yate/scripts/yatego-inline
ln -s /vagrant/cmd/callflow-static/callflow-static /opt/yate/share/yate/scripts/yatego-callflow-static
ln -s /vagrant/cmd/callflow-json/callflow-json /opt/yate/share/yate/scripts/yatego-callflow-json
ln -s /vagrant/cmd/callflow-vars/callflow-vars /opt/yate/share/yate/scripts/yatego-callflow-vars

echo "Yate config"

rsync -avz /vagrant/deployments/vagrant/configs/yate/ /opt/yate/

printf "\n[41587000201]\npassword=milan" >> /opt/yate/etc/yate/regfile.conf

printf "\n^900$=tone/congestion\n^920$=external/nodata/yatego-inline" >> /opt/yate/etc/yate/regexroute.conf
printf "\n^921$=external/nodata/yatego-callflow-static" >> /opt/yate/etc/yate/regexroute.conf
printf "\n^922$=external/nodata/yatego-callflow-json" >> /opt/yate/etc/yate/regexroute.conf
printf "\n^923$=external/nodata/yatego-callflow-vars" >> /opt/yate/etc/yate/regexroute.conf



