# go-build-servers
Various implementations of servers in Go

This repo is intended to be a a collection of implementations of different purpose servers.

## Proxy server

### Implementation
  
      Provides a proxy server and a Go client that chooses or needs to go through a proxy using a self-signed certificate.
      The proxy can be either HTTP or HTTPS. It creats a tunnel, through which the request is channeled. It does not alter the request.
      
      A script for creating a self-signed certificate.
      
      The client ensures the OS (Operating System) accepts the self-signed cert.
      
### Credits
https://medium.com/@mlowicki/http-s-proxy-in-golang-in-less-than-100-lines-of-code-6a51c2f2c38c https://forfuncsake.github.io/post/2017/08/trust-extra-ca-cert-in-go-app/

## URL shortener

A proxy that redirects a request to a previously stored short url to its real world equivalent.