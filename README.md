### Required software

* go 1.9.3
* docker
* docker-compose

### Create pki ( from https://github.com/OpenVPN/easy-rsa )

``sh
$ ./easyrsa init-pki
$ ./easyrsa build-ca
$ ./easyrsa build-server-full sunjin.local nopass
$ ./easyrsa build-server-full caddy.local nopass
$ ./easyrsa build-client-full 'client0' nopass
``

### Run servers

``sh
$ sudo docker-compose up
``

### Send a request to caddy proxy

``sh
$ go run client/client.go
``
