# go-dyip

A ddns client and server [GitHub](https://github.com/za-zliea/go-dyip)

## Feature

### Server

- Base on [atreugo](https://github.com/savsgio/atreugo) Web Server Framework.
- Use HTTP API to sync IP and DNS.
- Embedded web admin UI: log in with an auto-generated admin account to view all domains, check live DNS status, and trigger a manual sync from the browser.

### Client

- Use time.Ticker to call Server sync API.

### Support DNS Provider

- NONE(Just record IP, Provide Domain end with .internal)
- Aliyun
- Tencent
- Godaddy
- GoogleDomain(Not Google Cloud DNS, with Dynamic DNS, not test for ipv6)
- Cloudflare

*Note*

- Google have no api to query DNS, use net.LookupIP instead.
- Google use Dynamic DNS which has HTTP API.
- Cloudflare authenticate with the API Token as `sk` (Bearer token); `ak` is unused.
- Google authenticate with HTTP Basic Auth: `ak` is the username, `sk` is the password.

### Support Protocol
- IPV4
- IPV6

### SyncType
Each server record (`ips[].synctype`) and the client (`synctype`) carries a two-digit string:
- **First digit — console update**: `1` lets the web UI push an IP to this record, `0` blocks it (403). Server-side only; ignored by the client.
- **Last digit — IP source**: `1` ⇒ the client reads the local network interface and submits the IP (`localip`); `0` ⇒ the server derives the IP from the request (`realip` header / remote address). Set `interface` on the client when this is `1`.

## Build

### Prerequirements

- golang
- make
- node + pnpm (only to build the embedded web UI; see Frontend below)

### Frontend (embedded web UI)

The admin SPA source lives in `src/frontend/` (vue-element-plus-admin, mini branch). Its build output `src/frontend/dist/` is embedded into `dyip-server` via `//go:embed`, so **`go build src/server.go` fails unless `src/frontend/dist/` exists**. Build it once before building the server:

```shell
make frontend-build      # = cd src/frontend && pnpm install --frozen-lockfile=false && pnpm build
```

Build config lives in `vite.config.ts` (there is no `.env`). The build is **root-anchored**: `base: '/'` emits absolute asset URLs (`/assets/...`), and the API is called via absolute paths (`/front/api/...`). The SPA must be served at the site root — a sub-path deploy (e.g. behind a `/xxx-dyip/` reverse proxy) requires reconfiguring `base` and the API path prefixes.

For local development with hot reload, run the dev server alongside `dyip-server` (vite on `:5173`, proxies `/front` **and** `/api` to `$DYIP_DEV_SERVER` or `http://localhost:8080`):

```shell
cd src/frontend && pnpm dev
```

### Without Docker

```shell
make frontend-build   # build the web UI into src/frontend/dist (first run only / after UI changes)
make all
```

### With Docker

```shell
make image VERSION=[GIT TAG]
make image-alpine VERSION=[GIT TAG]
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
    	config file path, default server.yml (default "server.yml")
  -g	generate config, default server.yml
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
    	config file path, default client.yml (default "client.yml")
  -g	generate config, default client.yml
```

### Admin Web UI

The server serves an embedded web UI at `/`. On first start (or when generating a config with `-g`), an admin account is auto-generated — username `admin`, and a random 16-char password printed once to the logs:

```
generated admin account on first run username=admin password=xxxxxxxxxxxxxxxx
```

Log in at `http://<server>:<port>/`. From the UI you can:

- View all configured domains with their last synced IP and update time.
- Inspect a single record: its live DNS IP vs. the recorded IP, and its history.
- Trigger a manual sync to push a chosen IP to DNS — only allowed for records whose `synctype` first digit is `1` (others return 403; see SyncType).

The web UI authenticates with a JWT (HS256, 8h) signed with the server `token`. The machine sync API (`/api/sync`) is unaffected and still uses the global `token` + per-domain `auth`.

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
- provider: your-provider            # Support Provider (NONE/TENCENT/ALIYUN/GODADDY/GOOGLE/CLOUDFLARE)
  ak: abcde12345                     # Provider ak (USERNAME/AccessKey ID ...)
  sk: abcde12345                     # Provider sk (PASSWORD/AccessKey Secret ...)
  domain: your-doamin                # Domain
  subdomain: your-subdomain          # Subdomain
  auth: your-doamin-token-abce12345  # Client and server domain auth token
  protocol: IPV4                     # IPV4/IPV6 protocol
  synctype: "00"                     # Two digits: [console 0/1][local-ip 0/1] (see SyncType)
```

> The `admin` field is intentionally absent from the sample. On first start (and with `-g`), the server auto-generates an admin account: the password is bcrypt-hashed into `admin.password` and the **plaintext is printed once to the logs**. To reset the password, delete the `admin` block and restart.

### Client Config

```yaml
server: http://127.0.0.1:8080/       # Server url [Format http(s)://ip:port/prefix/]
token: your-token-abcde12345         # Client and server auth token
domain: your-subdomain.your-doamin   # Full domain
auth: your-doamin-token-abce12345    # Client and server domain auth token
interval: 300                        # Sync interval (second)
protocol: IPV4                       # IPV4/IPV6 protocol
synctype: "00"                       # Two digits: client only uses the last (local-ip 0/1)
interface: eth0                      # Local Interface Name (used when last digit is 1)
```