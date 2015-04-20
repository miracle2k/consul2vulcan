FROM gliderlabs/alpine:3.1
RUN apk --update add curl

RUN cd /tmp && curl -L https://github.com/hashicorp/consul-template/releases/download/v0.8.0/consul-template_0.8.0_linux_amd64.tar.gz| tar xz
RUN mv /tmp/consul*/consul-template /consul-template

ADD etcd-set.ctmpl /etcd-set.ctmpl
ADD run-consule-template.sh /start

ENTRYPOINT /bin/sh /start
