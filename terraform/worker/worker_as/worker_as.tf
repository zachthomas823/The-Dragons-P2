provider "aws" {
  #Two localfiles names as such. Each contains what they say, given to you from AWS.
  #DO NOT UPLOAD THESE FILES, make sure they are masked by the .gitignore
  #access_key = file("../../access_key")
  #secret_key = file("../../secret_key")
  region     = "us-east-2"
}

data "template_file" "user_data" {
  template = file("join.sh")
}

resource "aws_launch_template" "worker"{
  name_prefix   = "worker"
  image_id      = "ami-0920a73d71dd0ab71"
  instance_type = "t2.small"
  # user_data = base64encode(file("join.sh"))
  user_data = base64encode(data.template_file.user_data.rendered)

  #Generate your own Key_Name from AWS and use that here
  #DO NOT UPLOAD THESE FILES, make sure they are masked by the .gitignore
  key_name = "Temp"
  security_group_names = [aws_security_group.SSH.name]

  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_autoscaling_group" "worker_asg" {
  availability_zones = ["us-east-2a"]
  max_size           = 3
  desired_capacity   = 2
  min_size           = 1

  launch_template {
    id      = aws_launch_template.worker.id
    version = "$Latest"
  }
}

resource "aws_autoscaling_policy" "cpu_over" {
  name                   = "cpu_utilization_over"
  scaling_adjustment     = 1
  adjustment_type        = "ChangeInCapacity"
  cooldown               = 5
  autoscaling_group_name = aws_autoscaling_group.worker_asg.name
}

resource "aws_cloudwatch_metric_alarm" "cpu_over" {
  alarm_name          = "greater_than_cpu_usage"
  comparison_operator = "GreaterThanOrEqualToThreshold"
  evaluation_periods  = "1"
  metric_name         = "CPUUtilization"
  namespace           = "AWS/EC2"
  period              = "60"
  statistic           = "Average"
  threshold           = "90"

  dimensions = {
    AutoScalingGroupName = aws_autoscaling_group.worker_asg.name
  }

  alarm_description = "This metric monitors ec2 cpu utilization"
  alarm_actions     = [aws_autoscaling_policy.cpu_over.arn]
}

resource "aws_autoscaling_policy" "cpu_under" {
  name                   = "cpu_utilization_under"
  scaling_adjustment     = -1
  adjustment_type        = "ChangeInCapacity"
  cooldown               = 5
  autoscaling_group_name = aws_autoscaling_group.worker_asg.name
}

resource "aws_cloudwatch_metric_alarm" "cpu_under" {
  alarm_name          = "less_than_cpu_usage"
  comparison_operator = "LessThanThreshold"
  evaluation_periods  = "1"
  metric_name         = "CPUUtilization"
  namespace           = "AWS/EC2"
  period              = "60"
  statistic           = "Average"
  threshold           = "5"

  dimensions = {
    AutoScalingGroupName = aws_autoscaling_group.worker_asg.name
  }

  alarm_description = "This metric monitors ec2 cpu utilization"
  alarm_actions     = [aws_autoscaling_policy.cpu_under.arn]
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

# terraform {
#   backend "s3" {
#       bucket = "the-dragons-worker"
#       key    = "terraform.tfstate"
#       region = "us-east-2"
#   }
# }
