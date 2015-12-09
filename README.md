Autenticación mutua TLS en Go
=============================

Siguiendo el artículo [TLS Mutual Auth in GoLang](http://www.bite-code.com/2015/06/25/tls-mutual-auth-in-golang/)


## Ejemplo

Generar claves

    go run certgen/certgen.go -cert server.crt -key server.key 127.0.0.1
    go run certgen/certgen.go -client -cert client.crt -key client.key -name=client_1


Ejecutar el servidor

    go run server/server.go


Ejecutar el cliente

    go run client/client.go
