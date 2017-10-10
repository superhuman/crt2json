crt2json prints a JSON formatted summary of an SSL certificate.

I made it because it was painful to use `openssl` to perform these tasks. I
chose JSON as the output format because it's relatively easy to read, and you
can pass it to other tools (like `jq` or `gron`) to easily extract the field
that you want.

### Installation

```
go install github.com/superhuman/crt2json
```

### Usage

To decode a local certificate

```
crt2json <a.crt
```

To print certificate information for a website:

```
crt2json https://mail.superhuman.com
```

### TODO

Right now it only prints out the information I needed for my use-case. It would be nice to support more use-cases.
