#!/usr/bin/env bash

mkdir -p certs/
pushd certs/

rm -f *.{crt,key,csr,crl,cnf} 

# CA
openssl req -x509 \
  -newkey rsa:4096 \
  -days 30 \
  -nodes \
  -keyout ca.key \
  -out ca.crt \
  -subj "/"

# Server
openssl req \
  -newkey rsa:4096 \
  -nodes \
  -keyout server.key \
  -out server.csr \
  -subj "/" \

cat > server.cnf <<EOF
[v3_ca]
subjectAltName = IP:127.0.0.1
EOF

openssl x509 -req \
  -in server.csr \
  -days 7 \
  -CA ca.crt \
  -CAkey ca.key \
  -CAcreateserial \
  -out server.crt \
  -extensions v3_ca \
  -extfile server.cnf

# Client
openssl req \
  -newkey rsa:4096 \
  -nodes \
  -keyout client.key \
  -out client.csr \
  -subj "/"

openssl x509 -req \
  -in client.csr \
  -days 7 \
  -CA ca.crt \
  -CAkey ca.key \
  -CAcreateserial \
  -out client.crt

popd
