#!/usr/bin/env bash

kops delete cluster testhacluster.k8s.local --yes --state "s3://test-k8s-kops-ha"