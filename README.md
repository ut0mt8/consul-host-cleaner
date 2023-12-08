Implement a small daemon to reap left and failed nodes on consul ; also clean some scories from catalog.


```
Usage of ./cleaner:
  -consul-addresses="": go-netaddrs formated consul servers defintion [REQUIRED]
  -consul-grpc-port=8502: grpc port of consul server
  -consul-http-port=8500: http port of consul server
  -consul-http-timeout=5: http timeout for connecting to consul server
  -refresh-interval=20: interval between sync
```
