Syncs the service catalog of Consul to vulcand.

- based on consul-template, generating a template that executes curl calls to etcd
- syncs al services, using the service name as the backend name, and
  the service node ids as the backend server names.
- relies on TTL for removing backends/servers that are gone.

To use:

   docker run -e ETCD="192.168.59.103:4001" -e CONSUL="192.168.59.103:8500" elsdoerfer/consul2vulcand


If you have a consul service called "api-http", the backend in vulcand will
be called "api-http".