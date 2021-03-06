provider "aws" {
  #Two localfiles names as such. Each contains what they say, given to you from AWS.
  #DO NOT UPLOAD THESE FILES, make sure they are masked by the .gitignore
  access_key = "${file("../../access_key")}"
  secret_key = "${file("../../secret_key")}"
  region     = "us-east-2"
}

resource "aws_ami_from_instance" "master_image" {
  name               = "master_image"
  source_instance_id = "i-0f8bc5a1298eb613d"
}
