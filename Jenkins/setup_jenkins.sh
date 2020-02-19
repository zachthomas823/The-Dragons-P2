sudo apt-get update

#Make sure everything required is installed
sudo apt-get -y install \
    apt-transport-https \
    ca-certificates \
    curl \
    gnupg-agent \
    software-properties-common

#grab key
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -
apt-key fingerprint 0EBFCD88

#Add Download location
sudo add-apt-repository \
   "deb [arch=amd64] https://download.docker.com/linux/ubuntu \
   $(lsb_release -cs) \
   stable"

sudo apt-get update

#Install Docker
sudo apt-get -y install docker-ce docker-ce-cli containerd.io

#Start Jenkins
sudo docker run --name Jenkins --rm -d -p 80:8080 -p 50000:50000 jenkins/jenkins:lts

echo `sudo docker logs Jenkins`