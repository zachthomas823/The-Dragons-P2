provider "aws" {
  #Two localfiles names as such. Each contains what they say, given to you from AWS.
  #DO NOT UPLOAD THESE FILES, make sure they are masked by the .gitignore
  access_key = file("../../access_key")
  secret_key = file("../../secret_key")
  region     = "us-east-2"
}

data "template_file" "user_data" {
  template = file("join.sh")
}

resource "aws_launch_template" "worker"{
  name_prefix   = "worker"
  image_id      = "ami-0920a73d71dd0ab71"
  instance_type = "t2.micro"
  # user_data = base64encode(file("join.sh"))
  user_data = base64encode(data.template_file.user_data.rendered)

  #Generate your own Key_Name from AWS and use that here
  #DO NOT UPLOAD THESE FILES, make sure they are masked by the .gitignore
  key_name = "Temp"
  security_group_names = [aws_security_group.SSH.name]
}

resource "aws_autoscaling_group" "worker_asg" {
  availability_zones = ["us-east-2a"]
  desired_capacity   = 1
  max_size           = 5
  min_size           = 1

  launch_template {
    id      = aws_launch_template.worker.id
    version = "$Latest"
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