#!/bin/bash

HEARTBEAT=30
TTL=60

watch -n $HEARTBEAT sh etcd-set &
exec /consul-template -consul $CONSUL -template "/etcd-set.ctmpl:/etcd-set:sh /etcd-set" -retry 30s -wait 1s  --log-level debug $@