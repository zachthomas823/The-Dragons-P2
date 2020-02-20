provider "aws" {
  #Two localfiles names as such. Each contains what they say, given to you from AWS.
  #DO NOT UPLOAD THESE FILES, make sure they are masked by the .gitignore
  access_key = file("../../access_key")
  secret_key = file("../../secret_key")
  region     = "us-east-2"
}


resource "aws_instance" "worker"{
  ami           = "ami-0fc20dd1da406780b"
  instance_type = "t2.micro"

  #Generate your own Key_Name from AWS and use that here
  #DO NOT UPLOAD THESE FILES, make sure they are masked by the .gitignore
  key_name = "Temp"
  security_groups = [aws_security_group.SSH.name]

  connection {
    user = "ubuntu"
    type = "ssh"
    private_key = file("../../Temp.pem")
    host =  self.public_ip
    timeout = "20m"
  }
  # provisioner "remote-exec" {
  #   inline = [
  #     "sudo -s"
  #   ]
  # }
  provisioner "file" {
    source      = "setup_worker.sh"
    destination = "/tmp/setup_worker.sh"
  }
  # provisioner "file" {
  #   source      = "join.sh"
  #   destination = "/tmp/join.sh"
  # }
  provisioner "file" {
    source      = "join.config"
    destination = "/etc/init"
  }
  provisioner "remote-exec" {
    inline = [
      "/bin/bash /tmp/setup_worker.sh",
      "/bin/bash /tmp/join.sh"
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