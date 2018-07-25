# Building a TCP proxy server

  - Ubuntu 16.04
  - Go 1.10.3
  - curl 7.60.0 (x86_64-pc-linux-gnu) libcurl/7.60.0 OpenSSL/1.0.2g zlib/1.2.8 libssh2/1.5.0 (build locally to allow support for proxies, not available in curl from the Ubuntu package repository)
  
  The implementation shows how to build a proxy server and a client going through the proxy to retrieve a URL using Go. The project helped me get more insight into the proxying details. I also came across some gotchas, which provided even more insight :) 
  

## Running a proxy server

go run main.go -key ../../certs/localhost.key -pem ../../certs/localhost.pem  -proto https

### Issues

go run main.go 
panic: Get https://google.com: proxyconnect tcp: x509: certificate signed by unknown authority

  This means the OS is rejecting self-signed certificate. Please look at the Certificates section.

## Running a client

### Run a Go client

go run main.go certs.go [--insecure=true] [--server-host=localhost:9999] [--url=https://monzo.com] [--cert=../../localhost.pem]

go build -o client main.go certs.go
./client  [--insecure=true] [--server-host=localhost:9999] [--url=https://monzo.com] [--cert=../../localhost.pem]

## Use curl

/usr/local/bin/curl -Lv --proxy https://localhost:8888 --proxy-cacert path-to-self-signed-cert  (localhost.pem etc)  https://docs.docker.com

### Issues

curl: option --proxy-cacert: is unknown
                or
curl is not compiles with SSL support

    - Make sure you have curl built with HTTPS-proxy support (https://daniel.haxx.se/blog/2016/11/26/https-proxy-with-curl/):
        - download the right curl version (tar.gz) 
        - unpack
        - install libssl-dev package
        - enable ssl (and ssh) during curl build:
              
              ./configure --enable-ssl --enable-libssh2

## Certificates

### Compile curl if required:

    download (tar.gz) curl version curl 7.60.0
    install libssl-dev package
    enable ssl (and ssh) during curl build:
       ./configure --enable-ssl --enable-libssh2

### Create certificates:

a) script to create the key pair (certs/selfsigned.sh):
      The script creates the keys and certificate for "localhost" domain. There is also selfsignedFQDN.sh for
      creating a self-seigned certificate for the server domain (hostname -f)

b) For the Operating System to accept the CA certificate (localhost.pem) when running the Go client going through the proxy, the client code needs to add the created cert to the CA certificate pool.

  Another way to make the OS accept the self-signed cert is to create a custom openssl config file (localhost.cnf) and use the certs.orig/createCertsOpenSSLConf.sh script
  OR
  use the certs.orig/createCertsSAN.sh script.

  This should remove the need to complicate the client code with working with the root certificates pool. Neither of these two methods worked for me.

## Credits

https://medium.com/@mlowicki/http-s-proxy-in-golang-in-less-than-100-lines-of-code-6a51c2f2c38c
https://forfuncsake.github.io/post/2017/08/trust-extra-ca-cert-in-go-app/

## Further reading

https://fabianlee.org/2017/02/21/ubuntu-creating-a-self-signed-certificate-using-openssl-on-ubuntu/
https://developer.mozilla.org/en-US/docs/Glossary/Proxy_server
https://daniel.haxx.se/blog/2016/11/26/https-proxy-with-curl/

https://fabianlee.org/2018/02/17/ubuntu-creating-a-self-signed-san-certificate-using-openssl/
