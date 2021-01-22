# Mutual TLS - An example usage of 'approuvez' using mutual TLS auth{n,z} for the gRPC error

In this example, for simplicity purposes I will leverage a self-signed PKI

```bash
# You can clone the repo 
~$ git clone git@github.com:mvisonneau/approuvez.git; cd approuvez

# Generate certs using the available helper
~$ make certs
Generating a 4096 bit RSA private key
[..]

# Start the server using the certs
~$ approuvez \
     --tls-ca-cert certs/ca.crt \
     --tls-cert certs/server.crt \
     --tls-key certs/server.key \ 
     serve

# Use the client with the certs
~$ approuvez \
     --tls-ca-cert certs/ca.crt \
     --tls-cert certs/client.crt \
     --tls-key certs/client.key \ 
     ask -m"hey!" -u "foo@bar.baz"
```
