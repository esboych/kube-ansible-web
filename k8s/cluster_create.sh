#!/usr/bin/env bash

# run kops
kops create cluster \
       --state "s3://test-k8s-kops-ha" \
       --zones "eu-west-1a,eu-west-1b"  \
       --master-count 3 \
       --master-size=t2.micro \
       --node-count 2 \
       --node-size=t2.micro \
       --name testhacluster.k8s.local

# create real cluster
kops update cluster \
            --state "s3://test-k8s-kops-ha" \
            testhacluster.k8s.local --yes

# validate cluster
kops validate cluster --state "s3://test-k8s-kops-ha"