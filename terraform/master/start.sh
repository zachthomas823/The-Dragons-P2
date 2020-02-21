#! /bin/bash -xe
terraform apply -auto-approve
instance_ip=$(terraform output instance_ip)
ssh -o "StrictHostKeyChecking no" -i "../Temp.pem" ubuntu@$instance_ip 'cd sdn;/bin/bash ./run.sh'
