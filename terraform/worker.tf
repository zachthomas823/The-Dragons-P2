provider "aws" {
  #Two localfiles names as such. Each contains what they say, given to you from AWS.
  #DO NOT UPLOAD THESE FILES, make sure they are masked by the .gitignore
  access_key = "${file("./access_key")}"
  secret_key = "${file("./secret_key")}"
  region     = "us-east-2"
}

resource "aws_instance" "example" {
  ami           = "ami-0fc20dd1da406780b"
  instance_type = "t2.micro"

  #Generate your own Key_Name from AWS and use that here
  #DO NOT UPLOAD THESE FILES, make sure they are masked by the .gitignore
  key_name = "Temp"
  security_groups = ["${aws_security_group.SSH.name}"]

  connection {
    user = "ubuntu"
    type = "ssh"
    private_key = "${file("./Temp.pem")}"
    host =  self.public_ip
    timeout = "4m"
  }
  provisioner "file" {
    source      = "setup_node_env.sh"
    destination = "/tmp/setup_node_env.sh"
  }
  provisioner "remote-exec" {
    inline = [
      "sudo /bin/bash /tmp/setup_node_env.sh",
    ]
  }

}

resource "aws_security_group" "SSH" {
  name        = "allow_ssh"
  description = "Allow SSH traffic"


  ingress {
    from_port   = 22 
    to_port     = 22
    protocol =   "tcp"

    cidr_blocks =  ["64.189.196.114/32"]
  }

  egress {
    from_port       = 0
    to_port         = 0
    protocol        = "-1"
    cidr_blocks     = ["0.0.0.0/0"]
  }
}