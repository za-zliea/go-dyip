# go-dyip

A ddns client and server [GitHub](https://github.com/za-zliea/go-dyip)

## Feature

### Server

- Base on [atreugo](https://github.com/savsgio/atreugo) Web Server Framework.
- Use HTTP API to sync IP and DNS.

### Client

- Use time.Ticker to call Server sync API.

### Support DNS Provider

- NONE(Just record IP, Provide Domain end with .internal)
- Aliyun
- Tencent
- Godaddy
- Google(Dynamic DNS)
- Cloudflare

*Note*

- Google have no api to query DNS, use net.LookupIP instead.
- Google use Dynamic DNS which has HTTP API.

## Build

### Prerequirements

- golang
- make

### Without Docker

```shell
make all
```

### With Docker

```shell
make image VERSION=[GIT TAG]
make image-apline VERSION=[GIT TAG]
```

## Usage

### Server Usage

```
Usage:
  server startup:
    dyip-server [-c config file]
  server startup in background:
    nohup dyip-server [-c config file] &
  generate demo config file:
    dyip-server -g [-c config file]
  print usage:
    dyip-server -h
Options:
  -c string
    	config file path, default server.conf (default "server.conf")
  -g	generate config, default server.conf
  -h	print usage
```

### Client Usage

```shell
Usage:
  client startup:
    dyip-client [-c config file]
  client startup in background:
    nohup dyip-client [-c config file] &
  generate demo config file:
    dyip-client -g [-c config file]
  print usage:
    dyip-client -h
Options:
  -c string
    	config file path, default client.conf (default "client.conf")
  -g	generate config, default client.conf
```

## Docker

### Server

[Docker Hub](https://hub.docker.com/r/zliea/dyip-server)

```shell
docker run -d -p 8080:8080 --name dyip-server -v ./:/etc/dyip zliea/dyip-server:latest
```

### Client

[Docker Hub](https://hub.docker.com/r/zliea/dyip-client)

```shell
docker run -d --name dyip-client -v ./:/etc/dyip zliea/dyip-client:latest
```

## Config

### Server Config

```yaml
address: 127.0.0.1                   # Listen address
port: 8080                           # Listen port
realip: x-real-ip                    # IP to sync from header, use remote address if empty
token: your-token-abcde12345         # Client and server auth token
ips:
- provider: your-provider            # Support Provider (NONE/TENCENT/ALIYUN/GODADDY/GOOGLE)
  ak: abcde12345                     # Provider ak (USERNAME/AccessKey ID ...)
  sk: abcde12345                     # Provider sk (PASSWORD/AccessKey Secret ...)
  domain: your-doamin                # Domain
  subdomain: your-subdomain          # Subdomain
  auth: your-doamin-token-abce12345  # Client and server domain auth token
  protocol: IPV4                     # IPV4/IPV6 protocol
```

### Client Config

```yaml
server: http://127.0.0.1:8080/       # Server url [Format http(s)://ip:port/prefix/]
token: your-token-abcde12345         # Client and server auth token
domain: your-subdomain.your-doamin   # Full domain
auth: your-doamin-token-abce12345    # Client and server domain auth token
interval: 300                        # Sync interval (second)
protocol: IPV4                       # IPV4/IPV6 protocol
```