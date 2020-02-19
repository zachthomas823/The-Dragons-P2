#!/bin/bash

#Remove Bad Docker Installs
apt-get remove docker docker-engine docker.io containerd runc

apt-get update

#Make sure everything required is installed
apt-get -y install \
    apt-transport-https \
    ca-certificates \
    curl \
    gnupg-agent \
    software-properties-common

#grab key
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -
apt-key fingerprint 0EBFCD88

#Add Download location
add-apt-repository \
   "deb [arch=amd64] https://download.docker.com/linux/ubuntu \
   $(lsb_release -cs) \
   stable"

apt-get update

#Install Docker Proper
apt-get -y install docker-ce docker-ce-cli containerd.io

#Setup Docker.service
systemctl enable docker
systemctl start docker

setenforce 0
sed -i --follow-symlinks 's/^SELINUX=enforcing/SELINUX=disabled/' /etc/sysconfig/selinux

systemctl disable firewalld
systemctl stop firewalld

#Kuber requirements
sed -i '/swap/d' /etc/fstab
swapoff -a

#Get Kubelet kubeadm kubectl
apt-get update && apt-get install -y apt-transport-https curl
curl -s https://packages.cloud.google.com/apt/doc/apt-key.gpg | sudo apt-key add -
cat <<EOF | sudo tee /etc/apt/sources.list.d/kubernetes.list
deb https://apt.kubernetes.io/ kubernetes-xenial main
EOF
apt-get update
apt-get install -y kubelet kubeadm kubectl
apt-mark hold kubelet kubeadm kubectl

#Start Kubelet for overlay connections
systemctl enable kubelet
systemctl start kubelet