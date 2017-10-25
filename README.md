# GOHOST

A simple microservice to display a HTTP client's IP, hostname, and Owner.


## Example

```
$ curl gohost.domain.com
IP: 93.184.216.34
DNS: example.com
WHO: American Registry for Internet Numbers
```

## Routes

### /ip

return just the IP


### /host

return just the hostname


## Running

Run as a docker service `lanrat/gohost` or local binary.


### Configuration

The enviroment variable `LISTEN_ADDR` sets the IP and port to bind to. `0.0.0.0:8181` is used by default.

