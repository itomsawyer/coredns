# CoreDNS

[![Documentation](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](https://godoc.org/github.com/coredns/coredns)
[![Build Status](https://img.shields.io/travis/coredns/coredns.svg?style=flat-square&label=build)](https://travis-ci.org/coredns/coredns)
[![Code Coverage](https://img.shields.io/codecov/c/github/coredns/coredns/master.svg?style=flat-square)](https://codecov.io/github/coredns/coredns?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/coredns/coredns?style=flat-square)](https://goreportcard.com/report/coredns/coredns)

CoreDNS is a DNS server that started as a fork of [Caddy](https://github.com/mholt/caddy/). It has the
same model: it chains middleware. In fact it's so similar that CoreDNS is now a server type plugin for
Caddy.

CoreDNS is the successor to [SkyDNS](https://github.com/skynetservices/skydns). SkyDNS is a thin
layer that exposes services in etcd in the DNS. CoreDNS builds on this idea and is a generic DNS
server that can talk to multiple backends (etcd, kubernetes, etc.).

CoreDNS aims to be a fast and flexible DNS server. The keyword here is *flexible*: with CoreDNS you
are able to do what you want with your DNS data. And if not: write some middleware!

Currently CoreDNS is able to:

* Serve zone data from a file; both DNSSEC (NSEC only) and DNS are supported (*file*).
* Retrieve zone data from primaries, i.e., act as a secondary server (AXFR only) (*secondary*).
* Sign zone data on-the-fly (*dnssec*).
* Load balancing of responses (*loadbalance*).
* Allow for zone transfers, i.e., act as a primary server (*file*).
* Automatically load zone files from disk (*auto*)
* Caching (*cache*).
* Health checking endpoint (*health*).
* Use etcd as a backend, i.e., a 101.5% replacement for
  [SkyDNS](https://github.com/skynetservices/skydns) (*etcd*).
* Use k8s (kubernetes) as a backend (*kubernetes*).
* Serve as a proxy to forward queries to some other (recursive) nameserver (*proxy*).
* Provide metrics (by using Prometheus) (*metrics*).
* Provide query (*log*) and error (*error*) logging.
* Support the CH class: `version.bind` and friends (*chaos*).
* Profiling support (*pprof*).
* Rewrite queries (qtype, qclass and qname) (*rewrite*).
* Echo back the IP address, transport and port number used (*whoami*).

Each of the middlewares has a README.md of its own.

## Status

CoreDNS can be used as a authoritative nameserver for your domains, and should be stable enough to
provide you with good DNS(SEC) service.

There are still few [issues](https://github.com/coredns/coredns/issues), and work is ongoing on making
things fast and to reduce the memory usage.

All in all, CoreDNS should be able to provide you with enough functionality to replace parts of BIND
9, Knot, NSD or PowerDNS and SkyDNS. Most documentation is in the source and some blog articles can
be [found here](https://blog.coredns.io). If you do want to use CoreDNS in production, please
let us know and how we can help.

<https://caddyserver.com/> is also full of examples on how to structure a Corefile (renamed from
Caddyfile when forked).

## Compilation

CoreDNS (as a servertype plugin for Caddy) has a dependency on Caddy, but this is not different than
any other Go dependency. If you have the source of CoreDNS, get all dependencies:

    go get ./...

And then `go build` as you would normally do:

    go build

This should yield a `coredns` binary.

## Add Custom Middleware Plugin

### Include plugin package

Add your own plugin into core/coredns.go import block

### Edit your plugin priority in Middleware chain

Append your plugin name into ./middleware.cfg

### Auto generate related code

```
make gen
```

## Examples

When starting CoreDNS without any configuration, it loads the `whoami` middleware and starts
listening on port 53 (override with `-dns.port`), it should show the following:

~~~ txt
.:53
2016/09/18 09:20:50 [INFO] CoreDNS-001
CoreDNS-001
~~~

Any query send to port 53 should return some information; your sending address, port and protocol
used.

If you have a Corefile without a port number specified it will, by default, use port 53, but you
can override the port with the `-dns.port` flag:

~~~ txt
.: {
    proxy . 8.8.8.8:53
    log stdout
}
~~~

`./coredns -dns.port 1053`, runs the server on port 1053.

Start a simple proxy, you'll need to be root to start listening on port 53.

`Corefile` contains:

~~~ txt
.:53 {
    proxy . 8.8.8.8:53
    log stdout
}
~~~

Just start CoreDNS: `./coredns`.
And then just query on that port (53). The query should be forwarded to 8.8.8.8 and the response
will be returned. Each query should also show up in the log.

Serve the (NSEC) DNSSEC-signed `example.org` on port 1053, with errors and logging sent to stdout.
Allow zone transfers to everybody, but specically mention 1 IP address so that CoreDNS can send
notifies to it.

~~~ txt
example.org:1053 {
    file /var/lib/coredns/example.org.signed {
        transfer to *
        transfer to 2001:500:8f::53
    }
    errors stdout
    log stdout
}
~~~

Serve `example.org` on port 1053, but forward everything that does *not* match `example.org` to a recursive
nameserver *and* rewrite ANY queries to HINFO.

~~~ txt
.:1053 {
    rewrite ANY HINFO
    proxy . 8.8.8.8:53

    file /var/lib/coredns/example.org.signed example.org {
        transfer to *
        transfer to 2001:500:8f::53
    }
    errors stdout
    log stdout
}
~~~

### Zone Specification

The following Corefile fragment is legal, but does not explicitly define a zone to listen on:

~~~ txt
{
   # ...
}
~~~

This defaults to `.:53` (or whatever `-dns.port` is).

The next one only defines a port:
~~~ txt
:123 {
    # ...
}
~~~
This defaults to the root zone `.`, but can't be overruled with the `-dns.port` flag.

Just specifying a zone, default to listening on port 53 (can still be overridden with `-dns.port`:

~~~ txt
example.org {
    # ...
}
~~~

## Blog and Contact

Website: <https://coredns.io>
Twitter: [@corednsio](https://twitter.com/corednsio)
Docs: <https://miek.nl/tags/coredns/>
Github: <https://github.com/coredns/coredns>


## Systemd Service File

Use this as a systemd service file. It defaults to a coredns with a homedir of /home/coredns
and the binary lives in /opt/bin and the config in `/etc/coredns/Corefile`:

~~~ txt
[Unit]
Description=CoreDNS DNS server
Documentation=https://coredns.io
After=network.target

[Service]
PermissionsStartOnly=true
LimitNOFILE=8192
User=coredns
WorkingDirectory=/home/coredns
ExecStartPre=/sbin/setcap cap_net_bind_service=+ep /opt/bin/coredns
ExecStart=/opt/bin/coredns -conf=/etc/coredns/Corefile
ExecReload=/bin/kill -SIGUSR1 $MAINPID
Restart=on-failure

[Install]
WantedBy=multi-user.target
~~~

## Code structure
.
├── CONTRIBUTING.md
├── Corefile                            //configuration file
├── Corefile.vane
├── Dockerfile
├── LICENSE
├── Makefile
├── Makefile.bak
├── Makefile.release
├── README.md
├── core
│   ├── coredns.go
│   ├── dnsserver                       //dns server class
│   │   ├── address.go
│   │   ├── address_test.go
│   │   ├── config.go
│   │   ├── directives.go
│   │   ├── register.go
│   │   ├── server.go
│   │   └── zdirectives.go
│   └── zmiddleware.go
├── coredns
├── coredns.go                           //main
├── coremain
│   ├── run.go
│   ├── run_test.go
│   └── version.go
├── diff
├── diff1
├── directives_generate.go
├── error.log
├── glide.lock
├── glide.yaml
├── middleware                            //coredns available modules runs one by one
│   ├── auto
│   │   ├── README.md
│   │   ├── auto.go
│   │   ├── regexp.go
│   │   ├── regexp_test.go
│   │   ├── setup.go
│   │   ├── setup_test.go
│   │   ├── walk.go
│   │   ├── walk_test.go
│   │   ├── watcher_test.go
│   │   └── zone.go
│   ├── backend.go
│   ├── backend_lookup.go
│   ├── bind
│   │   ├── README.md
│   │   ├── bind.go
│   │   ├── bind_test.go
│   │   └── setup.go
│   ├── cache                               //cache module
│   │   ├── README.md
│   │   ├── cache.go                        
│   │   ├── cache_test.go
│   │   ├── handler.go                      //cache main procedure
│   │   ├── item.go
│   │   ├── item_test.go
│   │   ├── setup.go
│   │   └── setup_test.go
│   ├── chaos
│   │   ├── README.md
│   │   ├── chaos.go
│   │   ├── chaos_test.go
│   │   ├── setup.go
│   │   └── setup_test.go
│   ├── dnssec
│   │   ├── README.md
│   │   ├── black_lies.go
│   │   ├── black_lies_test.go
│   │   ├── cache.go
│   │   ├── cache_test.go
│   │   ├── dnskey.go
│   │   ├── dnssec.go
│   │   ├── dnssec_test.go
│   │   ├── handler.go
│   │   ├── handler_test.go
│   │   ├── responsewriter.go
│   │   ├── rrsig.go
│   │   ├── setup.go
│   │   └── setup_test.go
│   ├── erratic
│   │   ├── README.md
│   │   ├── erratic.go
│   │   ├── erratic_test.go
│   │   ├── setup.go
│   │   └── setup_test.go
│   ├── errors                              //error log module
│   │   ├── README.md
│   │   ├── errors.go
│   │   ├── errors_test.go
│   │   ├── setup.go
│   │   └── setup_test.go
│   ├── etcd
│   │   ├── README.md
│   │   ├── cname_test.go
│   │   ├── debug_test.go
│   │   ├── etcd.go
│   │   ├── group_test.go
│   │   ├── handler.go
│   │   ├── lookup_test.go
│   │   ├── msg
│   │   ├── multi_test.go
│   │   ├── other_test.go
│   │   ├── proxy_lookup_test.go
│   │   ├── setup.go
│   │   ├── setup_test.go
│   │   ├── stub.go
│   │   ├── stub_handler.go
│   │   └── stub_test.go
│   ├── file
│   │   ├── README.md
│   │   ├── closest.go
│   │   ├── closest_test.go
│   │   ├── cname_test.go
│   │   ├── delegation_test.go
│   │   ├── dnssec_test.go
│   │   ├── dnssex_test.go
│   │   ├── ds_test.go
│   │   ├── ent_test.go
│   │   ├── example_org.go
│   │   ├── file.go
│   │   ├── file_test.go
│   │   ├── glue_test.go
│   │   ├── lookup.go
│   │   ├── lookup_test.go
│   │   ├── notify.go
│   │   ├── nsec3_test.go
│   │   ├── reload_test.go
│   │   ├── secondary.go
│   │   ├── secondary_test.go
│   │   ├── setup.go
│   │   ├── setup_test.go
│   │   ├── tree
│   │   ├── wildcard.go
│   │   ├── wildcard_test.go
│   │   ├── xfr.go
│   │   ├── xfr_test.go
│   │   ├── zone.go
│   │   └── zone_test.go
│   ├── health
│   │   ├── README.md
│   │   ├── health.go
│   │   ├── health_test.go
│   │   └── setup.go
│   ├── kubernetes
│   │   ├── DEV-README.md
│   │   ├── README.md
│   │   ├── SkyDNS.md
│   │   ├── controller.go
│   │   ├── coredns.yaml.sed
│   │   ├── deploy.sh
│   │   ├── handler.go
│   │   ├── kubernetes.go
│   │   ├── kubernetes_test.go
│   │   ├── lookup.go
│   │   ├── setup.go
│   │   ├── setup_test.go
│   │   ├── subzone.go
│   │   └── subzone_test.go
│   ├── loadbalance
│   │   ├── README.md
│   │   ├── handler.go
│   │   ├── loadbalance.go
│   │   ├── loadbalance_test.go
│   │   └── setup.go
│   ├── log                                     //coredns query log module
│   │   ├── README.md
│   │   ├── log.go
│   │   ├── log_test.go
│   │   ├── setup.go
│   │   └── setup_test.go
│   ├── metrics
│   │   ├── README.md
│   │   ├── handler.go
│   │   ├── metrics.go
│   │   ├── metrics_test.go
│   │   ├── setup.go
│   │   ├── setup_test.go
│   │   ├── test
│   │   └── vars
│   ├── middleware.go
│   ├── middleware_test.go
│   ├── normalize.go
│   ├── normalize_test.go
│   ├── pkg
│   │   ├── debug
│   │   ├── dmtree
│   │   ├── dnsrecorder
│   │   ├── dnsutil
│   │   ├── edns
│   │   ├── nettree
│   │   ├── rcode
│   │   ├── replacer
│   │   ├── response
│   │   ├── singleflight
│   │   ├── storage
│   │   ├── strings
│   │   └── tls
│   ├── pprof
│   │   ├── README.md
│   │   ├── pprof.go
│   │   ├── setup.go
│   │   └── setup_test.go
│   ├── proxy
│   │   ├── README.md
│   │   ├── dns.go
│   │   ├── exchanger.go
│   │   ├── google.go
│   │   ├── google_rr.go
│   │   ├── google_test.go
│   │   ├── grpc.go
│   │   ├── lookup.go
│   │   ├── metrics.go
│   │   ├── pb
│   │   ├── policy.go
│   │   ├── policy_test.go
│   │   ├── proxy.go
│   │   ├── response.go
│   │   ├── setup.go
│   │   ├── upstream.go
│   │   ├── upstream_test.go
│   │   └── upstreamhost_wrapper.go
│   ├── reverse
│   │   ├── README.md
│   │   ├── network.go
│   │   ├── network_test.go
│   │   ├── reverse.go
│   │   ├── setup.go
│   │   └── setup_test.go
│   ├── rewrite
│   │   ├── README.md
│   │   ├── class.go
│   │   ├── condition.go
│   │   ├── condition_test.go
│   │   ├── field.go
│   │   ├── name.go
│   │   ├── reverter.go
│   │   ├── rewrite.go
│   │   ├── rewrite_test.go
│   │   ├── setup.go
│   │   ├── testdata
│   │   └── type.go
│   ├── root
│   │   ├── README.md
│   │   ├── root.go
│   │   └── root_test.go
│   ├── secondary
│   │   ├── README.md
│   │   ├── secondary.go
│   │   ├── setup.go
│   │   └── setup_test.go
│   ├── test
│   │   ├── doc.go
│   │   ├── file.go
│   │   ├── file_test.go
│   │   ├── helpers.go
│   │   ├── helpers_test.go
│   │   ├── responsewriter.go
│   │   └── server.go
│   ├── trace
│   │   ├── README.md
│   │   ├── setup.go
│   │   ├── setup_test.go
│   │   └── trace.go
│   ├── vane                                    //coredns proxy and query module
│   │   ├── dns.go                              //dns exchanger
│   │   ├── dns_test.go
│   │   ├── engine                              //coredns db loader module
│   │   ├── exchange_helper.go
│   │   ├── models                              //coredb definination models
│   │   ├── setup.go
│   │   ├── setup_test.go
│   │   ├── types.go
│   │   └── vane.go                             // query main procedure
│   └── whoami
│       ├── README.md
│       ├── setup.go
│       ├── setup_test.go
│       ├── whoami.go
│       └── whoami_test.go
├── middleware.cfg                              //coredns module configure, a list of which module will be compiled and detemine the order of running of modules
├── middleware.md
├── misc                                        //coredb init sql file
│   ├── gen.sh
│   ├── iwg.sql
│   ├── iwg_dump.sql                            //coredb init sql file
│   ├── iwg_nat_def.sql
│   └── local_test.sql
├── request
│   ├── request.go
│   └── request_test.go
├── rpm                                         //rpm make
│   ├── Corefile
│   ├── coredns.spec
│   ├── coredns.version
│   └── rpm-coredns.sh                          //under rpm director, run `./rpm-coredns.sh` will create a rpm package for coredns
├── stdout
├── test
│   ├── auto_test.go
│   ├── cache_test.go
│   ├── doc.go
│   ├── ds_file_test.go
│   ├── etcd_cache_debug_test.go
│   ├── etcd_test.go
│   ├── example_test.go
│   ├── file.go
│   ├── file_reload_test.go
│   ├── file_test.go
│   ├── kubernetes_test.go
│   ├── metrics_test.go
│   ├── middleware_dnssec_test.go
│   ├── middleware_test.go
│   ├── miek_test.go
│   ├── proxy_health_test.go
│   ├── proxy_test.go
│   ├── reverse_test.go
│   ├── server.go
│   ├── server_test.go
│   └── wildcard_test.go
├── utils                                   //useful tools
│   ├── digger                              //digger is a tool to dig @server with listed domains
│   │   ├── cm_zj_sh_domain
│   │   ├── digger
│   │   ├── digger.go
│   │   ├── domains.conf                    //digger input configuration file with json formated domains
│   │   ├── grep.sh
│   │   ├── ok_domain
│   │   └── result
│   ├── domain_add_script
│   │   ├── domain.sh
│   │   └── domainCM_HN
│   ├── ipnet_add_script
│   │   ├── hn_test
│   │   ├── ipnet.sh
│   │   └── ipnet_wl.sh
│   ├── ipnet_find_clientset
│   │   └── find_clientset.sh
│   ├── iptable_add_script
│   │   ├── guangdong_cm
│   │   ├── hn_cm
│   │   ├── huashudiaodu
│   │   ├── huashudiaodu.sql
│   │   ├── huashudiaodu1
│   │   ├── iptable.sh
│   │   ├── iptable_wl.sh
│   │   ├── new.sql
│   │   └── zj_cm_ctt
│   ├── iptable_find_netlink
│   │   └── find_netlink.sh
│   ├── mkconf
│   │   ├── auto.conf
│   │   ├── auto.example.conf
│   │   ├── conf
│   │   ├── main.go
│   │   ├── mkconf
│   │   └── pretty
│   ├── nginx_log
│   │   ├── http_traffic.sh
│   │   └── top_domain_traffic.sh
│   └── top_domain_from_Log
│       └── get.sh
├── vane.log
└── vendor                              //third party packages
