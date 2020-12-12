FROM  wxext-registry.101.com/sdp/centos:centos7
LABEL maintainer="LinQing"

ENV TZ=Asia/Shanghai

RUN curl -o /usr/local/lib/go1.15.6.linux-amd64.tar.gz http://localhost:57180/cicd/go1.15.6.linux-amd64.tar.gz \
&& tar -xvf /usr/local/lib/go1.15.6.linux-amd64.tar.gz -C /usr/local/lib \
&& mkdir -p /usr/local/lib/go/gopath \
&& yum -y install git lrzsz \
&& yum clean all \
&& rm -rf /usr/local/lib/go1.15.6.linux-amd64.tar.gz \
&& ln -snf /usr/share/zoneinfo/$TZ /etc/localtime \
&& echo $TZ > /etc/timezone

ENV GOPRIVATE=.gitlab.com,.gite.com
ENV GOPROXY=https://goproxy.cn,direct
ENV GOSUMDB=off
ENV GOROOT=/usr/local/lib/go
ENV GOPATH=${GOROOT}/gopath
ENV PATH=${PATH}:${GOROOT}/bin

CMD ["sleep","inf"]