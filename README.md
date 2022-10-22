Balanced HTTP2 REST service with TLS mutual authentication
==========================================================

## Goals

This project realize two goals:

  - Build an **HTTP2** service and a client with **SSL/TLS [mutual authentication](https://github.com/jomoespe/go-tls-mutual-auth)** over **HTTP2** between them.
  - Be able to distribute workload across multiple services with a [layer 4 (TPC) load balancer](https://en.wikipedia.org/wiki/Load_balancing_(computing)).

The services are implemented in [Go language](https://golang.org/)



## Requisites

For SSL/TLS mutual authentication:

  - [Golang SDK](https://golang.org/dl/). Tested with version 1.5.1

For load balancing:

  - [Docker](http://www.docker.com/). Tested with version 1.9.1
  - [Docker compose](https://www.docker.com/docker-compose). Tested with version 1.5.1


## Introduction

Mutual authentication refers to two parties authenticating each other at the same time. That is a client authenticating itself to a server and that server authenticating itself to the client in such a way that both parties are assured of the others' identity. In adition to SSL, muutual authentication provides authentication and non-repudiation of the client, using using digital signatures. 

This process it performed with certificates interchange. That is both client and server send its own certificates in connection handshaking, the client validate if the server certificate is valid and then the server validates the client certificate validation. If all it's ok the connection is stablished. After this, the server can read client centificate information to perform client identification.

Because we are realizing client authentication and identification in the service process, we cannot put an HTTP/S (layer 7) reverse proxy/load balancer in front of a service instances. This is why we configure a TCP (layer 4) reverse proxy/load balancer.


## SSL/TSL mutual authentication

The project have three main components:

  - The server.
  - The client.
  - A certificate generation tool.

### Build

To build all components

```bash
make clean all 
```

There are make targets for each component.

```bash
make [cert] [serverd] [tlsclient]
```

### Certificate generation tool

Generate certificate:

```bash
    ./cert [-org <"Organization name">] [-name <"subject name">] [-duration <duration>] [-cert <certificate filename>] [-key <private key filename>] [-client [<true|false>]] [ip|servers....]
```

Example: generate a server certificate for 127.0.0.1 and localhost.localdomain

```bash
./cert -cert server.crt -key server.key 127.0.0.1 localhost.localdomain
```

Example: generate a client certificate with client_1 name

```bash
./cert -client -cert client.crt -key client.key -name=client_1
```

### The server

```bash
./serverd
```

### The client


```bash
./tlsclient
```

## References

  - [SSL/TLS Mutual Auth in GoLang](http://www.bite-code.com/2015/06/25/tls-mutual-auth-in-golang/)
