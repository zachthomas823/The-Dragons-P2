#! /bin/bash -xe
sudo kubeadm join 172.31.5.246:6443 --token z3coug.0yviajl3b25d3ato     --discovery-token-ca-cert-hash sha256:4cd766e79b72e4e908ba9e68d0406c94f5180ead24d1ea906158693f118a2aca
mkdir success