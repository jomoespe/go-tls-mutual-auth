all: cert serverd clean

clean:
	@ rm -f serverd tlsclient cert

cert:
	@ CGO_ENABLED=0 go build -a -installsuffix cgo -o cert certgen/certgen.go

serverd:
	@ CGO_ENABLED=0 go build -a -installsuffix cgo -o serverd server/server.go 

tlsclient:
	@ CGO_ENABLED=0 go build -a -installsuffix cgo -o tlsclient client/client.go 

