Syncs the service catalog of Consul to vulcand.

- based on consul-template, generating a template that executes curl calls to etcd
- syncs al services, using the service name as the backend name, and
  the service node ids as the backend server names.
- relies on TTL for removing backends/servers that are gone.

To use:

   docker run -e ETCD="http://192.168.59.103:4001" -e CONSUL="192.168.59.103:8500" -e TTL=60 elsdoerfer/consul2vulcand

If you have a consul service called "api-http", the backend in vulcand will
be called "api-http".


To test this locally
--------------------

Run consul, etcd:

    $ docker run -p 8500:8500 progrium/consul --server --bootstrap
    $ docker run -p 4001:4001 coreos/etcd --addr 192.168.59.103:4001

Run this:

    $ docker run -e ETCD="192.168.59.103:4001" -e CONSUL="192.168.59.103:8500" -e TTL=60 -e VERBOSE=1 elsdoerfer/consul2vulcand

Get some services into consul:

    # docker run -v /var/run/docker.sock:/tmp/docker.sock gliderlabs/registrator -ip 192.168.59.103 consul://192.168.59.103:8500

Define a backend:

    $ etcdctl --peers http://192.168.59.103:4001 mkdir /vulcand/backends/etcd-4001/backend