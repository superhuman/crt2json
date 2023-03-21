crt2json prints a JSON formatted summary of an SSL certificate.

I made it because it was painful to use `openssl` to perform these tasks. I
chose JSON as the output format because it's relatively easy to read, and you
can pass it to other tools (like `jq` or `gron`) to easily extract the field
that you want.

### Installation

```
go install github.com/superhuman/crt2json@latest
```

### Usage

```
Usage: crt2json  [-sni HOSTNAME] [ SERVER | FILENAME ]

Prints out a JSON summary of an SSL certificate.

If the argument exists on disk then the file is assumed to be a certificate file, otherwise
it is interpreted as a URL or a hostname to connect to.
  -sni string
    	server name indication to instruct the server to return the correct certificate
```

### TODO

Right now it only prints out the information I needed for my use-cases. It would be nice to support more use-cases.
