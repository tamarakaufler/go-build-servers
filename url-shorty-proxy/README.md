# URL shortener proxy
The proxy translates a URL shortcut into its full real world equivalent. The mappings are stored in a datastore. The implementation makes it easy to substitute one datastore for another. Current stores include Posgres and MongoDB.

The databases run in a container and are bootstrapped for use with the url-shorty-proxy (see https://github.com/tamarakaufler/go-build-servers/tree/master/bootstrapping).

To swap one datastore for another, make a change in the handler.go. To make a swap, replace 'psql' with
'mongo' and 'PSQL' with 'MGO', and visa versa.

## Usage
a) Run a datastore where the short URL/full URL mappings are stored
b) Run the proxy, either directly or in a container
c) 
  In Postman:

      POST http://localhost:8888/ggl.c/google.com
      DELETE http://localhost:8888/ggl.c

  In the browser:

      http://localhost:8888/ggl.c
      http://localhost:8888/info/ggl.c
      http://localhost:8888/list


## Reading

Inspired by:

https://github.com/douglasmakey/ursho

https://www.youtube.com/watch?v=SWKuYLqouIY&index=20&t=0s&list=PL64wiCrrxh4Jisi7OcCJIUpguV_f5jGnZ
