#!/usr/bin/env sh


docker rmi wxext-registry.101.com/sdp/centos:centos7-v1.2
docker build -t wxext-registry.101.com/sdp/centos:centos7-v1.2 ./
docker push wxext-registry.101.com/sdp/centos:centos7-v1.2
