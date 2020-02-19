provider "aws" {
  #Two localfiles names as such. Each contains what they say, given to you from AWS.
  #DO NOT UPLOAD THESE FILES, make sure they are masked by the .gitignore
  access_key = file("../access_key")
  secret_key = file("../secret_key")
  region     = "us-east-2"
}

resource "aws_instance" "master" {
  ami           = "ami-0fc20dd1da406780b"
  instance_type = "t2.medium"

  #Generate your own Key_Name from AWS and use that here
  #DO NOT UPLOAD THESE FILES, make sure they are masked by the .gitignore
  key_name = "Temp"
  security_groups = [aws_security_group.SSH.name]

  connection {
    user = "ubuntu"
    type = "ssh"
    private_key = file("../Temp.pem")
    host =  self.public_ip
    timeout = "20m"
  }

  provisioner "remote-exec" {
    inline = [
      "mkdir services",
      "mkdir terraform",
      "mkdir terraform/worker",
      "mkdir terraform/worker/worker_as",
      ]
  }

  provisioner "file" {
    source      = "setup_master.sh"
    destination = "/tmp/setup_master.sh"
  }
  provisioner "file" {
    source      = "../worker/worker_as/worker_as.tf"
    destination = "/home/ubuntu/terraform/worker/worker_as/worker_as.tf"
  }
  provisioner "file" {
    source      = "master.tf"
    destination = "/home/ubuntu/terraform/master/master.tf"
  }
  provisioner "file" {
    source      = "../secret_key"
    destination = "/home/ubuntu/terraform/secret_key"
  }
  provisioner "file" {
    source      = "../access_key"
    destination = "/home/ubuntu/terraform/access_key"
  }
  provisioner "file" {
    source      = "../Temp.pem"
    destination = "/home/ubuntu/terraform/Temp.pem"
  }
  provisioner "file" {
    source      = "../../kubernetes"
    destination = "/home/ubuntu/kubernetes"
  }
  provisioner "file" {
    source      = "../../sdn"
    destination = "/home/ubuntu/sdn"
  }
  provisioner "file" {
    source = "terraform"
    destination = "/home/ubuntu/terraform/terraform"
  }
  provisioner "file" {
    source = "join.sh"
    destination = "/home/ubuntu/terraform/worker/worker_as/join.sh"
  }

  provisioner "remote-exec" {
    inline = [
      "sudo /bin/bash /tmp/setup_master.sh",
      "cd terraform",
      "sudo chmod 777 terraform",
      "sudo mv terraform /usr/local/bin",
      "cd worker/worker_as",
      "sudo chmod 777 join.sh",
      "kubeadm token create --print-join-command >> join.sh",
      "terraform init",
      "terraform apply --auto-approve"
    ]
  }
}

  resource "aws_security_group" "SSH" {
  description = "Allow SSH traffic"


  ingress {
    from_port   = 0 
    to_port     = 0
    protocol =   "-1"

    cidr_blocks =  ["0.0.0.0/0"]
  }

  egress {
    from_port       = 0
    to_port         = 0
    protocol        = "-1"
    cidr_blocks     = ["0.0.0.0/0"]
  }
}