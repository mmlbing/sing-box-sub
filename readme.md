# sing-box-sub

Convert Clash-Meta subscriptions to sing-box config. Supports CLI and HTTP server modes.

## Install

Download the ipk for your architecture from [Releases](../../releases), or grab the binary directly:

```sh
# OpenWrt
opkg install sing-box-sub_*.ipk
/etc/init.d/sing-box-sub enable && /etc/init.d/sing-box-sub start
```

## Usage

**CLI mode:**

```sh
sing-box-sub -u https://example.com/sub -t templates/momo.json -o config.json
sing-box-sub -u https://sub1.com,https://sub2.com -o config.json  # multiple subs
```

**HTTP server mode:**

```sh
sing-box-sub -d   # listens on 0.0.0.0:40533
```

```
GET /convert?sub=https://sub1.com&sub=https://sub2.com&tpl=https://tpl.json
GET /health
```

## Supported protocols

hysteria2 / anytls / vmess → sing-box outbound

## Template

Templates are standard sing-box JSON with `{all}` placeholder for node insertion, `filter` for per-group filtering, and `nodeFilter` for global node exclusion. See [doc/template-guide.md](doc/template-guide.md).
