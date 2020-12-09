#!/usr/bin/env sh

rm -rf /var/www/html/cicd/sdp-app-cli
cd ../../
go build -o bin/sdp-app-cli cmd/appinfo/sdp_appinfo.go
cp bin/sdp-app-cli /var/www/html/cicd/
cd -
docker rmi wxext-registry.101.com/sdp/centos:centos7-v1
docker build -t wxext-registry.101.com/sdp/centos:centos7-v1 ./
docker push wxext-registry.101.com/sdp/centos:centos7-v1
