FROM  wxext-registry.101.com/sdp/centos:centos7
LABEL maintainer="LinQing"

ENV TZ=Asia/Shanghai

RUN wget -P /usr/local/bin/ http://172.24.135.10:57180/cicd/sdp-app-cli \
&& chmod +x /usr/local/bin/sdp-app-cli \
&& ln -snf /usr/share/zoneinfo/$TZ /etc/localtime \
&& echo $TZ > /etc/timezone