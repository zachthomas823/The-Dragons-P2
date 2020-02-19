provider "aws" {
  #Two localfiles names as such. Each contains what they say, given to you from AWS.
  #DO NOT UPLOAD THESE FILES, make sure they are masked by the .gitignore
  access_key = file("../../access_key")
  secret_key = file("../../secret_key")
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
    private_key = file("../../Temp.pem")
    host =  self.public_ip
    timeout = "20m"
  }

  provisioner "file" {
    source      = "../setup_master.sh"
    destination = "/tmp/setup_master.sh"
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