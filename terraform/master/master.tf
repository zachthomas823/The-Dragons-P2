provider "aws" {
  #Two localfiles names as such. Each contains what they say, given to you from AWS.
  #DO NOT UPLOAD THESE FILES, make sure they are masked by the .gitignore
  access_key = "${file("../access_key")}"
  secret_key = "${file("../secret_key")}"
  region     = "us-east-2"
}

resource "aws_instance" "master" {
  ami           = "ami-0fc20dd1da406780b"
  instance_type = "t2.medium"

  #Generate your own Key_Name from AWS and use that here
  #DO NOT UPLOAD THESE FILES, make sure they are masked by the .gitignore
  key_name = "Temp"
  security_groups = ["${aws_security_group.SSH.name}"]

  connection {
    user = "ubuntu"
    type = "ssh"
    private_key = "${file("../Temp.pem")}"
    host =  self.public_ip
    timeout = "20m"
  }

  provisioner "remote-exec" {
    inline = [
      "mkdir pods",
      "mkdir services",
    ]
  }

  provisioner "file" {
    source      = "setup_master.sh"
    destination = "/tmp/setup_master.sh"
  }
  provisioner "file" {
    source      = "../../kubernetes/services/client-server.yaml"
    destination = "/home/ubuntu/services/client-server.yaml"
  }
  provisioner "file" {
    source      = "../../kubernetes/pods/html-server.yaml"
    destination = "/home/ubuntu/pods/html-server.yaml"
  }
  provisioner "file" {
    source      = "../worker/worker.tf"
    destination = "/home/ubuntu/terraform/worker.tf"
  }
  provisioner "file" {
    source      = "master.tf"
    destination = "/home/ubuntu/terraform/master.tf"
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

  provisioner "remote-exec" {
    inline = [
      "sudo /bin/bash /tmp/setup_master.sh",
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