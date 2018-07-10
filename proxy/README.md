# Building of a proxy server

  - Ubuntu 16.04
  - Go 1.10.3
  - curl 7.60.0 (x86_64-pc-linux-gnu) libcurl/7.60.0 OpenSSL/1.0.2g zlib/1.2.8 libssh2/1.5.0 (build locally to allow support for proxies, not available in curl from the Ubuntu package repository)

## Run the proxy server

go run main.go -key ../../certs/localhost.key -pem ../../certs/localhost.pem  -proto https

### Issues

go run main.go 
panic: Get https://google.com: proxyconnect tcp: x509: certificate signed by unknown authority

  Means the OS is rejecting self-signed certificate:

    Solution: Create a custom OpenSSL config that will be used when creating the certificates

## Run the client

go run main.go certs.go [--insecure=true] [--server-host=localhost:9999] [--url=https://monzo.com] [--cert=../../localhost.pem]

go build -o client main.go certs.go
./client  [--insecure=true] [--server-host=localhost:9999] [--url=https://monzo.com] [--cert=../../localhost.pem]

## curl

/usr/local/bin/curl -Lv --proxy https://localhost:8888 --proxy-cacert path-to-self-signed-cert  (localhost.pem etc)  https://docs.docker.com

### Issues

    - Make sure you have curl built with HTTPS-proxy support (https://daniel.haxx.se/blog/2016/11/26/https-proxy-with-curl/):
        - download the right curl version (tar.gz) 
        - unpack
        - install libssl-dev package
        - enable ssl (and ssh) during curl build:
              
              ./configure --enable-ssl --enable-libssh2

## Certificates

download (tar.gz) curl version curl 7.60.0
install libssl-dev package
enable ssl (and ssh) during curl build:
   ./configure --enable-ssl --enable-libssh2

create certificates:

create openssh config to make the OS accept self-signed certificates (for running curl):

a) openssl config for custom domain (localhost etc)
      https://fabianlee.org/2018/02/17/ubuntu-creating-a-self-signed-san-certificate-using-openssl/

b) script to create the key pair (certs/create_certs.sh):


## Credits

https://medium.com/@mlowicki/http-s-proxy-in-golang-in-less-than-100-lines-of-code-6a51c2f2c38c
https://forfuncsake.github.io/post/2017/08/trust-extra-ca-cert-in-go-app/

## Further reading

https://fabianlee.org/2017/02/21/ubuntu-creating-a-self-signed-certificate-using-openssl-on-ubuntu/
https://developer.mozilla.org/en-US/docs/Glossary/Proxy_server
https://daniel.haxx.se/blog/2016/11/26/https-proxy-with-curl/

https://fabianlee.org/2018/02/17/ubuntu-creating-a-self-signed-san-certificate-using-openssl/
