provider "aws" {
  #Two localfiles names as such. Each contains what they say, given to you from AWS.
  #DO NOT UPLOAD THESE FILES, make sure they are masked by the .gitignore
  access_key = "AKIAWUGRQXZRTQRQX3YE"
  secret_key = "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"
  region     = "us-east-2"
}
resource "aws_instance" "Jenkins" {
    ami           = "ami-0fc20dd1da406780b"
    instance_type   = "t2.micro"
    key_name        = "basekey"
    security_groups = ["${aws_security_group.Jenkins_Group.name}"]

    connection {
    user = "ubuntu"
    type = "ssh"
    private_key = file("./basekey.pem")
    host =  self.public_ip
    timeout = "10m"
    }
}


resource "aws_security_group" "Jenkins_Group" {
    name        = "Jenkins_Group"
    description = "Allows traffic on port 22 for ssh and 80 for tcp"

    ingress {
        from_port = 22
        to_port = 22
        protocol = "tcp"

        cidr_blocks = ["0.0.0.0/0"]
    }
        ingress {
        from_port = 80
        to_port = 80
        protocol = "tcp"

        cidr_blocks = ["0.0.0.0/0"]
    }

    egress {
        from_port       = 0
        to_port         = 0
        protocol        = "-1"
        cidr_blocks     = ["0.0.0.0/0"]
  }
}