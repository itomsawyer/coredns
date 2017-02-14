drop database if exists iwg;
create database iwg;
use iwg;

drop view if exists route_view;
drop view if exists base_route_view;
drop view if exists policy_view;
drop view if exists clientset_view;
drop view if exists domain_view;
drop view if exists netlink_view;
drop view if exists netlinkset_view;
drop view if exists src_view;
drop view if exists filter_view;

drop table if exists viewer;
drop table if exists filter;
drop table if exists ipnet;
drop table if exists clientset;
drop table if exists domainlink;
drop table if exists domain;
drop table if exists domain_pool;
drop table if exists route;
drop table if exists routeset;
drop table if exists iptable;
drop table if exists netlink;
drop table if exists netlinkset;
drop table if exists policy_detail;
drop table if exists policy;
drop table if exists ldns;
drop table if exists outlink;
drop table if exists rr;
drop table if exists rrset;


create table clientset (
id int not null auto_increment,
name varchar(127) not null,
info varchar(255) not null default "",
primary key (id),
unique key (name)
)DEFAULT CHARSET=utf8 comment "client ipnet set";

insert into clientset (name, info) values("unknown", "src ipnet that igw has no idea where it belongs to");

create table ipnet (
id int not null auto_increment,
ip_start varchar(40) not null,
ip_end varchar(40) not null,
ipnet varchar(40) not null, 
mask  tinyint unsigned not null,
clientset_id int,
primary key(id),
unique key(ip_start, ip_end),
foreign key(clientset_id) references clientset(id) on delete cascade
)DEFAULT CHARSET=utf8 comment "client ipnet";

create table ipnet_wl (
id int not null auto_increment,
ip_start varchar(40) not null,
ip_end varchar(40) not null,
ipnet varchar(40) not null, 
mask  tinyint unsigned not null,
clientset_id int,
primary key(id),
unique key(ip_start, ip_end),
foreign key(clientset_id) references clientset(id) on delete cascade
)DEFAULT CHARSET=utf8 comment "client ipnet whitelist";

create table domain_pool (
id int not null auto_increment,
name varchar(127) not null ,
info varchar(255) not null default "",
primary key (id),
unique key (name)
)DEFAULT CHARSET=utf8 comment "serve domains set";

insert into domain_pool (name, info) values("global", "Base domain pool for all of domains which are not specifically configured");

create table domain (
id int not null auto_increment,
domain varchar(255) not null,
domain_pool_id int,
primary key(id),
unique key(domain, domain_pool_id),
foreign key(domain_pool_id) references domain_pool(id) on delete cascade
)DEFAULT CHARSET=utf8 comment "serve domains";


create table outlink (
id int not null auto_increment,
name varchar(127) not null,
addr varchar(40) not null,
typ varchar(32) not null default "normal",
enable bool not null default true,
unavailable smallint unsigned not null default 0 comment "if other than zero, outlink is unavailable, each bit indicate different reason",
primary key(id),
unique key(name)
)DEFAULT CHARSET=utf8 comment "network gateway, aka outlink";

create table netlinkset(
id int not null auto_increment,
name varchar(127) not null,
primary key(id),
unique key(name)
)DEFAULT CHARSET=utf8 comment "netlink set ";

create table netlink (
id int not null auto_increment,
isp  varchar(127) not null,
region varchar(127) not null,
typ varchar(32) not null default "normal",
primary key(id)
)DEFAULT CHARSET=utf8 comment "netlink (isp + province or CP) of a target ip";

insert into netlink (isp, region) values("unknown", "unknown");

create table iptable (
id int not null auto_increment,
ip_start varchar(40) not null,
ip_end varchar(40) not null,
ipnet varchar(40) not null, 
mask tinyint unsigned not null,
netlink_id int,
primary key(id),
unique key(ip_start, ip_end),
foreign key(netlink_id) references netlink(id) on delete restrict
)DEFAULT CHARSET=utf8 comment "IP to netlink";


create table iptable_wl (
id int not null auto_increment,
ip_start varchar(40) not null,
ip_end varchar(40) not null,
ipnet varchar(40) not null, 
mask tinyint unsigned not null,
netlink_id int,
primary key(id),
unique key(ip_start, ip_end),
foreign key(netlink_id) references netlink(id) on delete restrict
)DEFAULT CHARSET=utf8 comment "IP to netlink whitelist";


create table domainlink(
id int not null auto_increment,
domain_pool_id int not null,
netlink_id int not null, 
netlinkset_id int not null,
enable bool not null default true,
primary key(id),
unique key(domain_pool_id, netlink_id),
foreign key(domain_pool_id) references domain_pool(id) on delete restrict,
foreign key(netlink_id) references netlink(id) on delete restrict,
foreign key(netlinkset_id) references netlinkset(id) on delete cascade
)DEFAULT CHARSET=utf8 comment "bind domain_pool and netlink to a netlinkset";

create table routeset (
id int not null auto_increment,
name varchar(127) not null,
info varchar(255) not null default "",
primary key(id),
unique key(name)
)DEFAULT CHARSET=utf8;

create table route (
id int not null auto_increment,
outlink_id int not null,
netlinkset_id int not null,
routeset_id int not null,
enable bool not null default true,
priority smallint not null default 20,
score smallint not null default 50 comment "netlink performance index",
unavailable smallint unsigned not null default 0 comment "if other than zero, route is unavailable, each bit indicate different reason",
primary key(id),
unique key(netlinkset_id, outlink_id, routeset_id),
foreign key(outlink_id) references outlink(id) on delete restrict,
foreign key(netlinkset_id) references netlinkset(id) on delete restrict,
foreign key(routeset_id) references routeset(id) on delete cascade 
)DEFAULT CHARSET=utf8 comment "Performance of using gateway to serve paricular netlink";


create table policy (
id int not null auto_increment,
name varchar(127) not null,
primary key(id),
unique key(name)
)DEFAULT CHARSET=utf8 comment "policy index of choose ldns upstream forwarder";

create table ldns (
id int not null auto_increment,
name varchar(127) not null,
addr varchar(40) not null,
typ varchar(32) not null default "A",
enable bool not null default true,
unavailable smallint unsigned not null default 0 comment "if other than zero, ldns is unavailable with each bit indicate different reason",
primary key(id),
unique key(name)
)DEFAULT CHARSET=utf8 comment "upstream ldns info";

create table rrset(
id int not null auto_increment,
name varchar(64) not null default "",
enable bool not null default true,
primary key(id)
)DEFAULT CHARSET=utf8 comment "dns rrset(resource record)";

create table rr (
id int not null auto_increment,
rrtype int not null,
rrdata varchar(255) not null,
ttl int unsigned not null default 300,
rrset_id int,
enable bool not null default true,
primary key(id),
foreign key(rrset_id) references rrset(id) on delete cascade
)DEFAULT CHARSET=utf8 comment "dns rr(resource record)";

create table policy_detail (
id int not null auto_increment,
policy_id int not null,
policy_sequence int not null default 0,
enable bool not null default true,
priority smallint not null default 20,
weight smallint not null default 100,
op varchar(127) not null default "and",
op_typ varchar(32) not null default "builtin",
ldns_id int not null,
rrset_id int,
primary key(id),
foreign key(ldns_id) references ldns(id) on delete restrict,
foreign key(rrset_id) references rrset(id) on delete restrict,
foreign key(policy_id) references policy(id) on delete cascade
)DEFAULT CHARSET=utf8 comment "policy detail of choose ldns upstream forwarder";

create table viewer(
id int not null auto_increment,
clientset_id int not null,
domain_pool_id int not null,
routeset_id int not null, 
policy_id int not null,
enable bool not null default true,
primary key(id),
foreign key(clientset_id) references clientset(id) on delete restrict,
foreign key(domain_pool_id) references domain_pool(id) on delete restrict,
foreign key(routeset_id) references routeset(id) on delete restrict,
foreign key(policy_id) references policy(id) on delete restrict
)DEFAULT CHARSET=utf8 comment "map of <clientset , domain_pool> -> <policy, routeset>";

create table filter(
id int not null auto_increment,
src_ip_start varchar(40),
src_ip_end varchar(40),
clientset_id int,
domain_id int,
dst_ip varchar(40),
outlink_id int,
enable bool not null default true,
primary key(id),
foreign key(clientset_id) references clientset(id) on delete cascade,
foreign key(domain_id) references domain(id) on delete cascade,
foreign key(outlink_id) references outlink(id) on delete cascade
)DEFAULT CHARSET=utf8 comment "custom route strategy like iptables";

create ALGORITHM = MERGE view filter_view as
select
filter.id ,
src_ip_start ,
inet_aton(src_ip_start) as src_ip_start_int,
src_ip_end ,
inet_aton(src_ip_end) as src_ip_end_int,
clientset_id ,
domain_id ,
dst_ip ,
outlink_id ,
domain.domain,
outlink.name as outlink_name,
outlink.addr as outlink_addr
from domain, filter, outlink, clientset
where  filter.clientset_id = clientset.id and filter.outlink_id = outlink.id and filter.domain_id = domain.id and filter.enable = true;

create ALGORITHM = MERGE view clientset_view as
select ipnet.id as ipnet_id, ip_start, inet_aton(ip_start) as ip_start_int, ip_end, inet_aton(ip_end) as ip_end_int, ipnet, mask, clientset_id, name as clientset_name from ipnet join clientset on clientset.id = ipnet.clientset_id;

create ALGORITHM = MERGE view clientset_wl_view as
select ipnet_wl.id as ipnet_wl_id, ip_start, inet_aton(ip_start) as ip_start_int, ip_end, inet_aton(ip_end) as ip_end_int, ipnet, mask, clientset_id, name as clientset_name from ipnet_wl join clientset on clientset.id = ipnet_wl.clientset_id;


create ALGORITHM = MERGE view domain_view as
select domain.id as domain_id, domain, domain_pool_id, domain_pool.name as pool_name from domain join domain_pool on domain.domain_pool_id= domain_pool.id;

create ALGORITHM = MERGE view netlink_view as
select iptable.id as iptable_id, ip_start, inet_aton(ip_start) as ip_start_int, ip_end, inet_aton(ip_end) as ip_end_int, ipnet, mask, netlink_id, isp, region, typ
from netlink join iptable on iptable.netlink_id = netlink.id;

create ALGORITHM = MERGE view netlink_wl_view as
select iptable_wl.id as iptable_wl_id, ip_start, inet_aton(ip_start) as ip_start_int, ip_end, inet_aton(ip_end) as ip_end_int, ipnet, mask, netlink_id, isp, region, typ
from netlink join iptable_wl on iptable_wl.netlink_id = netlink.id;

create ALGORITHM = MERGE view dst_view as
select domainlink.domain_pool_id, domainlink.netlink_id, domainlink.netlinkset_id, domain_pool.name as domain_pool_name, isp, region
from  domain_pool, netlink, domainlink
where domain_pool.id = domainlink.domain_pool_id and netlink.id = domainlink.netlink_id
and domainlink.enable = true;

create ALGORITHM = MERGE view route_view as
select routeset_id, routeset.name as routeset_name,  netlinkset_id, netlinkset.name as netlinkset_name, route.id as route_id, min(route.priority) as route_priority, max(route.score) as route_score, outlink_id, outlink.name as outlink_name,  outlink.addr as outlink_addr, outlink.typ as outlink_typ
from netlinkset, outlink, route, routeset
where netlinkset.id = route.netlinkset_id and outlink.id = route.outlink_id and route.routeset_id = routeset.id
and outlink.enable = true and route.enable = true and outlink.unavailable = 0 and route.unavailable = 0
group by routeset_id;


create ALGORITHM = MERGE view base_route_view as
select routeset_id, routeset.name as routeset_name,  netlinkset_id, netlinkset.name as netlinkset_name, route.id as route_id, route.priority as route_priority, route.score as route_score, outlink_id, outlink.name as outlink_name,  outlink.addr as outlink_addr, outlink.typ as outlink_typ
from netlinkset, outlink, route, routeset
where netlinkset.id = route.netlinkset_id and outlink.id = route.outlink_id and route.routeset_id = routeset.id
and outlink.enable = true and route.enable = true and outlink.unavailable = 0 and route.unavailable = 0;


/*
create ALGORITHM = MERGE view base_route_view as
select
viewer.clientset_id,
clientset.name as clientset_name,
viewer.domain_pool_id,
domain_pool.name as domain_pool_name,
viewer.routeset_id, routeset.name as routeset_name,
netlinkset_id, netlinkset.name as netlinkset_name,
route.id as route_id, route.priority as route_priority, route.score as route_score,
outlink_id, outlink.name as outlink_name,  outlink.addr as outlink_addr, outlink.typ as outlink_typ
from clientset, domain_pool, viewer, routeset, netlinkset, route, outlink
where
clientset.id = viewer.clientset_id and domain_pool.id = viewer.domain_pool_id and viewer.routeset_id = routeset.id
and netlinkset.id = route.netlinkset_id and outlink.id = route.outlink_id and route.routeset_id = routeset.id
and outlink.enable = true and route.enable = true and outlink.unavailable = 0 and route.unavailable = 0 and viewer.enable = true;
*/


create ALGORITHM = MERGE view policy_view as
select policy.id as policy_id, policy.name as policy_name,
policy_detail.policy_sequence, policy_detail.priority, policy_detail.weight, policy_detail.op, policy_detail.op_typ,
policy_detail.ldns_id, ldns.name, ldns.addr, ldns.typ,
policy_detail.rrset_id
from  policy, policy_detail, ldns
where policy.id = policy_detail.policy_id and ldns.id = policy_detail.ldns_id 
and ldns.unavailable = 0 and ldns.enable = true and policy_detail.enable = true;

create ALGORITHM = MERGE view src_view as
select clientset_id, clientset.name as clientset_name, domain_pool_id, domain_pool.name as domain_pool_name, routeset_id, routeset.name as routeset_name, policy_id, policy.name as policy_name
from clientset, domain_pool, viewer, routeset, policy
where clientset.id = viewer.clientset_id and domain_pool.id = viewer.domain_pool_id and routeset.id = viewer.routeset_id and policy.id = viewer.policy_id
and viewer.enable = true
