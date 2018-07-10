# go-build-servers
Various implementations of servers in Go

This repo is intended to be a a collection of implementations of different purpose servers.

  - proxy implementation:
  
      Provides a proxy server and a Go client that chooses or needs to go through a proxy using a self-signed certificate.
      The proxy can be either HTTP or HTTPS.
      
      A script for creating a self-signed certificate.
      
      The client ensure the OS (Operating System) accepts the self-signed cert.
      
