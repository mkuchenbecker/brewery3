#!/bin/bash

sudo kubeadm init --pod-network-cidr=10.244.0.0/16 &&
  rm -rf $HOME/.kube &&
  mkdir -p $HOME/.kube &&
  sudo cp -i /etc/kubernetes/admin.conf $HOME/.kube/config &&
  sudo chown $(id -u):$(id -g) $HOME/.kube/config &&
  kubectl apply -f "https://cloud.weave.works/k8s/net?k8s-version=$(kubectl version | base64 | tr -d '\n')"
