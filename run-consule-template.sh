#!/bin/bash

export HEARTBEAT=30

watch -n $HEARTBEAT sh etcd-set &
exec /consul-template -consul $CONSUL -template "/etcd-set.ctmpl:/etcd-set:sh /etcd-set" -retry 30s -wait 1s $@