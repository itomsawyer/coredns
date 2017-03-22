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

drop table if exists _locker;
drop table if exists viewer;
drop table if exists filter;
drop table if exists ipnet;
drop table if exists ipnet_wl;
drop table if exists clientset;
drop table if exists domainlink;
drop table if exists domain;
drop table if exists domain_pool;
drop table if exists route;
drop table if exists routeset;
drop table if exists iptable;
drop table if exists iptable_wl;
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

insert into clientset (name, info) values("default", "global default");

create table ipnet (
id int not null auto_increment,
ip_start varchar(40) not null,
ip_end varchar(40) not null,
ipnet varchar(40) not null, 
mask  tinyint unsigned not null,
clientset_id int,
primary key(id),
unique key(ip_start, ip_end),
unique key(ipnet, mask),
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
unique key(ipnet, mask),
foreign key(clientset_id) references clientset(id) on delete cascade
)DEFAULT CHARSET=utf8 comment "client ipnet whitelist";

create table domain_pool (
id int not null auto_increment,
name varchar(127) not null ,
info varchar(255) not null default "",
enable bool not null default true,
unavailable smallint unsigned not null default 0,
domain_monitor bool not null default false,
primary key (id),
unique key (name)
)DEFAULT CHARSET=utf8 comment "serve domains set";

insert into domain_pool (name, info) values("default", "global default");

create table domain (
id int not null auto_increment,
domain varchar(255) not null,
domain_pool_id int,
primary key(id),
unique key(domain),
foreign key(domain_pool_id) references domain_pool(id) on delete cascade
)DEFAULT CHARSET=utf8 comment "serve domains";


create table outlink (
id int not null auto_increment,
name varchar(127) not null,
addr varchar(255) not null,
typ varchar(32) not null default "normal",
enable bool not null default true,
unavailable smallint unsigned not null default 0 comment "if other than zero, outlink is unavailable, each bit indicate different reason",
primary key(id),
unique key(name)
)DEFAULT CHARSET=utf8 comment "network gateway, aka outlink";

insert into outlink (name, addr) values ("default", "0.0.0.0");

create table netlinkset(
id int not null auto_increment,
name varchar(127) not null,
primary key(id),
unique key(name)
)DEFAULT CHARSET=utf8 comment "netlink set ";

insert into netlinkset(name) values("default");

create table netlink (
id int not null auto_increment,
isp  varchar(127) not null,
region varchar(127) not null default "",
typ varchar(32) not null default "normal",
primary key(id),
unique key(isp, region)
)DEFAULT CHARSET=utf8 comment "netlink (isp + province or CP) of a target ip";

insert into netlink (isp, region) values("default", "global default");

create table iptable (
id int not null auto_increment,
ip_start varchar(40) not null,
ip_end varchar(40) not null,
ipnet varchar(40) not null, 
mask tinyint unsigned not null,
netlink_id int,
primary key(id),
unique key(ip_start, ip_end),
unique key(ipnet, mask),
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
unique key(ipnet, mask),
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

insert into domainlink (domain_pool_id, netlink_id, netlinkset_id) values (1,1,1);

create table routeset (
id int not null auto_increment,
name varchar(127) not null,
info varchar(255) not null default "",
primary key(id),
unique key(name)
)DEFAULT CHARSET=utf8;

insert into routeset(name, info) values("default", "global default");

create table route (
id int not null auto_increment,
routeset_id int not null,
netlinkset_id int not null,
outlink_id int not null,
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

insert into route (netlinkset_id, routeset_id, outlink_id) values (1,1,1);

create table policy (
id int not null auto_increment,
name varchar(127) not null,
primary key(id),
unique key(name)
)DEFAULT CHARSET=utf8 comment "policy index of choose ldns upstream forwarder";

insert into policy (name) values ("default");

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

insert into ldns (name, addr) values("default", "223.5.5.5");

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

insert into policy_detail (policy_id, ldns_id) values(1, 1);

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

insert into viewer (clientset_id, domain_pool_id, routeset_id, policy_id) values (1,1,1,1);

create table filter(
id int not null auto_increment,
src_ip_start varchar(40),
src_ip_end   varchar(40),
clientset_id int,
domain_id    int,
domain_pool_id int,
dst_ip_start varchar(40),
dst_ip_end   varchar(40),
netlink_id   int,
target_ip    varchar(40),
outlink_id   int,
enable bool not null default true,
primary key(id),
foreign key(clientset_id) references clientset(id) on delete cascade,
foreign key(netlink_id) references netlink(id) on delete cascade,
foreign key(domain_id) references domain(id) on delete cascade,
foreign key(domain_pool_id) references domain_pool(id) on delete cascade,
foreign key(outlink_id) references outlink(id) on delete cascade
)DEFAULT CHARSET=utf8 comment "custom route strategy like iptables";

create table _locker(
id int not null auto_increment,
name varchar(32) not null default "",
clientset_id   int,
domain_pool_id int,
netlink_id     int,

viewer_id      int,
routeset_id    int,

domainlink_id  int,
netlinkset_id  int,

route_id       int,
outlink_id     int,

policy_id      int,
policy_detail_id int,
ldns_id        int,

primary key(id),
foreign key(clientset_id) references clientset(id) on delete restrict,
foreign key(domain_pool_id) references domain_pool(id) on delete restrict,
foreign key(netlink_id) references netlink(id) on delete restrict,
foreign key(viewer_id) references viewer(id) on delete restrict,
foreign key(routeset_id) references routeset(id) on delete restrict,
foreign key(domainlink_id) references domainlink(id) on delete restrict,
foreign key(netlinkset_id) references netlinkset(id) on delete restrict,
foreign key(route_id) references route(id) on delete restrict,
foreign key(outlink_id) references outlink(id) on delete restrict,
foreign key(policy_id) references policy(id) on delete restrict,
foreign key(policy_detail_id) references policy_detail(id) on delete restrict,
foreign key(ldns_id) references ldns(id) on delete restrict
)DEFAULT CHARSET=utf8 comment "Internal lock for default data, DO NOT MODIFY unless you known what you are doing";

insert into _locker values(0, "global default data locker", 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1);

create table oplog (
id bigint not null auto_increment,
opr varchar(16),
action  varchar(6),
tbl varchar(16),
row_id  int,
time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
primary key(id)
)DEFAULT CHARSET=utf8 comment "Log of operations to other tables in this db";

-- ALTER TABLE clientset
-- ADD COLUMN opr VARCHAR(16) NOT NULL DEFAULT 'unknown';
-- ALTER TABLE domain
-- ADD COLUMN opr VARCHAR(16) NOT NULL DEFAULT 'unknown';
-- ALTER TABLE domain_pool
-- ADD COLUMN opr VARCHAR(16) NOT NULL DEFAULT 'unknown';
-- ALTER TABLE domainlink
-- ADD COLUMN opr VARCHAR(16) NOT NULL DEFAULT 'unknown';
-- ALTER TABLE filter
-- ADD COLUMN opr VARCHAR(16) NOT NULL DEFAULT 'unknown';
-- ALTER TABLE ipnet
-- ADD COLUMN opr VARCHAR(16) NOT NULL DEFAULT 'unknown';
-- ALTER TABLE ipnet_wl
-- ADD COLUMN opr VARCHAR(16) NOT NULL DEFAULT 'unknown';
-- ALTER TABLE iptable
-- ADD COLUMN opr VARCHAR(16) NOT NULL DEFAULT 'unknown';
-- ALTER TABLE iptable_wl
-- ADD COLUMN opr VARCHAR(16) NOT NULL DEFAULT 'unknown';
-- ALTER TABLE ldns
-- ADD COLUMN opr VARCHAR(16) NOT NULL DEFAULT 'unknown';
-- ALTER TABLE netlink
-- ADD COLUMN opr VARCHAR(16) NOT NULL DEFAULT 'unknown';
-- ALTER TABLE netlinkset
-- ADD COLUMN opr VARCHAR(16) NOT NULL DEFAULT 'unknown';
-- ALTER TABLE outlink
-- ADD COLUMN opr VARCHAR(16) NOT NULL DEFAULT 'unknown';
-- ALTER TABLE policy
-- ADD COLUMN opr VARCHAR(16) NOT NULL DEFAULT 'unknown';
-- ALTER TABLE policy_detail
-- ADD COLUMN opr VARCHAR(16) NOT NULL DEFAULT 'unknown';
-- ALTER TABLE route
-- ADD COLUMN opr VARCHAR(16) NOT NULL DEFAULT 'unknown';
-- ALTER TABLE routeset
-- ADD COLUMN opr VARCHAR(16) NOT NULL DEFAULT 'unknown';
-- ALTER TABLE rr
-- ADD COLUMN opr VARCHAR(16) NOT NULL DEFAULT 'unknown';
-- ALTER TABLE rrset
-- ADD COLUMN opr VARCHAR(16) NOT NULL DEFAULT 'unknown';
-- ALTER TABLE viewer
-- ADD COLUMN opr VARCHAR(16) NOT NULL DEFAULT 'unknown';

create ALGORITHM = MERGE view filter_view as
select
filter.id ,
src_ip_start ,
inet_aton(src_ip_start) as src_ip_start_int,
src_ip_end ,
inet_aton(src_ip_end) as src_ip_end_int,
clientset_id ,
clientset.name as clientset_name, 
domain_id ,
domain.domain,
filter.domain_pool_id, 
domain_pool.name as domain_pool_name,
dst_ip_start,
inet_aton(dst_ip_start) as dst_ip_start_int,
dst_ip_end,
inet_aton(dst_ip_end) as dst_ip_end_int,
netlink_id,
CONCAT_WS(':', netlink.isp , netlink.region) as netlink_name,
target_ip,
inet_aton(target_ip) as target_ip_int,
outlink_id ,
outlink.name as outlink_name,
outlink.addr as outlink_addr
from filter, domain, domain_pool, clientset, netlink, outlink
where  filter.netlink_id = netlink.id and  filter.clientset_id = clientset.id and filter.outlink_id = outlink.id and  filter.domain_pool_id = domain_pool.id and filter.domain_id = domain.id and filter.enable = true;

create ALGORITHM = MERGE view clientset_view as
select ipnet.id as ipnet_id, ip_start, inet_aton(ip_start) as ip_start_int, ip_end, inet_aton(ip_end) as ip_end_int, ipnet, mask, clientset_id, name as clientset_name from ipnet join clientset on clientset.id = ipnet.clientset_id;

create ALGORITHM = MERGE view clientset_wl_view as
select ipnet_wl.id as ipnet_wl_id, ip_start, inet_aton(ip_start) as ip_start_int, ip_end, inet_aton(ip_end) as ip_end_int, ipnet, mask, clientset_id, name as clientset_name from ipnet_wl join clientset on clientset.id = ipnet_wl.clientset_id;

create ALGORITHM = MERGE view domain_view as
select
domain.id as domain_id,
domain,
domain_pool_id,
domain_pool.name as pool_name,
domain_pool.domain_monitor
from domain join domain_pool on domain.domain_pool_id= domain_pool.id
where domain_pool.enable = true and domain_pool.unavailable = 0 ;

create ALGORITHM = MERGE view netlink_view as
select iptable.id as iptable_id, ip_start, inet_aton(ip_start) as ip_start_int, ip_end, inet_aton(ip_end) as ip_end_int, ipnet, mask, netlink_id, isp, region, typ
from netlink join iptable on iptable.netlink_id = netlink.id;

create ALGORITHM = MERGE view netlink_wl_view as
select iptable_wl.id as iptable_wl_id, ip_start, inet_aton(ip_start) as ip_start_int, ip_end, inet_aton(ip_end) as ip_end_int, ipnet, mask, netlink_id, isp, region, typ
from netlink join iptable_wl on iptable_wl.netlink_id = netlink.id;

create ALGORITHM = MERGE view dst_view as
select domainlink.domain_pool_id, domainlink.netlink_id, domainlink.netlinkset_id, domain_pool.name as domain_pool_name, isp, region, netlinkset.name as netlinkset_name
from  domain_pool, netlink, domainlink, netlinkset
where domain_pool.id = domainlink.domain_pool_id and netlink.id = domainlink.netlink_id and domain_pool.enable = true and domain_pool.unavailable = 0
and domainlink.enable = true and netlinkset.id = domainlink.netlinkset_id;

create ALGORITHM = MERGE view route_view as
select routeset_id, routeset.name as routeset_name,  netlinkset_id, netlinkset.name as netlinkset_name, route.id as route_id, min(route.priority) as route_priority, max(route.score) as route_score, outlink_id, outlink.name as outlink_name,  outlink.addr as outlink_addr, outlink.typ as outlink_typ
from netlinkset, outlink, route, routeset
where netlinkset.id = route.netlinkset_id and outlink.id = route.outlink_id and route.routeset_id = routeset.id
and outlink.enable = true and route.enable = true and outlink.unavailable = 0 and route.unavailable = 0 and route.score != 0
group by routeset_id, netlinkset_id;

create ALGORITHM = MERGE view base_route_view as
select routeset_id, routeset.name as routeset_name,  netlinkset_id, netlinkset.name as netlinkset_name, route.id as route_id, route.priority as route_priority, route.score as route_score, outlink_id, outlink.name as outlink_name,  outlink.addr as outlink_addr, outlink.typ as outlink_typ
from netlinkset, outlink, route, routeset
where netlinkset.id = route.netlinkset_id and outlink.id = route.outlink_id and route.routeset_id = routeset.id
and outlink.enable = true and route.enable = true and outlink.unavailable = 0 and route.unavailable = 0 and route.score != 0;

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
and domain_pool.enable = true and domain_pool.unavailable = 0 and viewer.enable = true;


-- delimiter !
-- create trigger oplog_clientset_insert after insert on clientset for each row
-- insert into oplog values(0, opr, "clientset", "insert", id, null);
-- !
-- create trigger oplog_clientset_update after update on clientset for each row
-- insert into oplog values(0, opr, "clientset", "update", id, null);
-- !
-- create trigger oplog_clientset_delete after delete on clientset for each row
-- insert into oplog values(0, opr, "clientset", "delete", id, null);
-- !
-- create trigger oplog_domain_insert after insert on domain for each row
-- insert into oplog values(0, opr, "domain", "insert", id, null);
-- !
-- create trigger oplog_domain_update after update on domain for each row
-- insert into oplog values(0, opr, "domain", "update", id, null);
-- !
-- create trigger oplog_domain_delete after delete on domain for each row
-- insert into oplog values(0, opr, "domain", "delete", id, null);
-- !
-- create trigger oplog_domain_pool_insert after insert on domain_pool for each row
-- insert into oplog values(0, opr, "domain_pool", "insert", id, null);
-- !
-- create trigger oplog_domain_pool_update after update on domain_pool for each row
-- insert into oplog values(0, opr, "domain_pool", "update", id, null);
-- !
-- create trigger oplog_domain_pool_delete after delete on domain_pool for each row
-- insert into oplog values(0, opr, "domain_pool", "delete", id, null);
-- !
-- create trigger oplog_domainlink_insert after insert on domainlink for each row
-- insert into oplog values(0, opr, "domainlink", "insert", id, null);
-- !
-- create trigger oplog_domainlink_update after update on domainlink for each row
-- insert into oplog values(0, opr, "domainlink", "update", id, null);
-- !
-- create trigger oplog_domainlink_delete after delete on domainlink for each row
-- insert into oplog values(0, opr, "domainlink", "delete", id, null);
-- !
-- create trigger oplog_filter_insert after insert on filter for each row
-- insert into oplog values(0, opr, "filter", "insert", id, null);
-- !
-- create trigger oplog_filter_update after update on filter for each row
-- insert into oplog values(0, opr, "filter", "update", id, null);
-- !
-- create trigger oplog_filter_delete after delete on filter for each row
-- insert into oplog values(0, opr, "filter", "delete", id, null);
-- !
-- create trigger oplog_ipnet_insert after insert on ipnet for each row
-- insert into oplog values(0, opr, "ipnet", "insert", id, null);
-- !
-- create trigger oplog_ipnet_update after update on ipnet for each row
-- insert into oplog values(0, opr, "ipnet", "update", id, null);
-- !
-- create trigger oplog_ipnet_delete after delete on ipnet for each row
-- insert into oplog values(0, opr, "ipnet", "delete", id, null);
-- !
-- create trigger oplog_ipnet_wl_insert after insert on ipnet_wl for each row
-- insert into oplog values(0, opr, "ipnet_wl", "insert", id, null);
-- !
-- create trigger oplog_ipnet_wl_update after update on ipnet_wl for each row
-- insert into oplog values(0, opr, "ipnet_wl", "update", id, null);
-- !
-- create trigger oplog_ipnet_wl_delete after delete on ipnet_wl for each row
-- insert into oplog values(0, opr, "ipnet_wl", "delete", id, null);
-- !
-- create trigger oplog_iptable_insert after insert on iptable for each row
-- insert into oplog values(0, opr, "iptable", "insert", id, null);
-- !
-- create trigger oplog_iptable_update after update on iptable for each row
-- insert into oplog values(0, opr, "iptable", "update", id, null);
-- !
-- create trigger oplog_iptable_delete after delete on iptable for each row
-- insert into oplog values(0, opr, "iptable", "delete", id, null);
-- !
-- create trigger oplog_iptable_wl_insert after insert on iptable_wl for each row
-- insert into oplog values(0, opr, "iptable_wl", "insert", id, null);
-- !
-- create trigger oplog_iptable_wl_update after update on iptable_wl for each row
-- insert into oplog values(0, opr, "iptable_wl", "update", id, null);
-- !
-- create trigger oplog_iptable_wl_delete after delete on iptable_wl for each row
-- insert into oplog values(0, opr, "iptable_wl", "delete", id, null);
-- !
-- create trigger oplog_ldns_insert after insert on ldns for each row
-- insert into oplog values(0, opr, "ldns", "insert", id, null);
-- !
-- create trigger oplog_ldns_update after update on ldns for each row
-- insert into oplog values(0, opr, "ldns", "update", id, null);
-- !
-- create trigger oplog_ldns_delete after delete on ldns for each row
-- insert into oplog values(0, opr, "ldns", "delete", id, null);
-- !
-- create trigger oplog_netlink_insert after insert on netlink for each row
-- insert into oplog values(0, opr, "netlink", "insert", id, null);
-- !
-- create trigger oplog_netlink_update after update on netlink for each row
-- insert into oplog values(0, opr, "netlink", "update", id, null);
-- !
-- create trigger oplog_netlink_delete after delete on netlink for each row
-- insert into oplog values(0, opr, "netlink", "delete", id, null);
-- !
-- create trigger oplog_netlinkset_insert after insert on netlinkset for each row
-- insert into oplog values(0, opr, "netlinkset", "insert", id, null);
-- !
-- create trigger oplog_netlinkset_update after update on netlinkset for each row
-- insert into oplog values(0, opr, "netlinkset", "update", id, null);
-- !
-- create trigger oplog_netlinkset_delete after delete on netlinkset for each row
-- insert into oplog values(0, opr, "netlinkset", "delete", id, null);
-- !
-- create trigger oplog_outlink_update after update on outlink for each row
-- insert into oplog values(0, opr, "outlink", "update", id, null);
-- !
-- create trigger oplog_outlink_delete after delete on outlink for each row
-- insert into oplog values(0, opr, "outlink", "delete", id, null);
-- !
-- create trigger oplog_policy_insert after insert on policy for each row
-- insert into oplog values(0, opr, "policy", "insert", id, null);
-- !
-- create trigger oplog_policy_update after update on policy for each row
-- insert into oplog values(0, opr, "policy", "update", id, null);
-- !
-- create trigger oplog_policy_delete after delete on policy for each row
-- insert into oplog values(0, opr, "policy", "delete", id, null);
-- !
-- create trigger oplog_policy_detail_insert after insert on policy_detail for each row
-- insert into oplog values(0, opr, "policy_detail", "insert", id, null);
-- !
-- create trigger oplog_policy_detail_update after update on policy_detail for each row
-- insert into oplog values(0, opr, "policy_detail", "update", id, null);
-- !
-- create trigger oplog_policy_detail_delete after delete on policy_detail for each row
-- insert into oplog values(0, opr, "policy_detail", "delete", id, null);
-- !
-- create trigger oplog_route_insert after insert on route for each row
-- insert into oplog values(0, opr, "route", "insert", id, null);
-- !
-- create trigger oplog_route_update after update on route for each row
-- insert into oplog values(0, opr, "route", "update", id, null);
-- !
-- create trigger oplog_route_delete after delete on route for each row
-- insert into oplog values(0, opr, "route", "delete", id, null);
-- !
-- create trigger oplog_routeset_insert after insert on routeset for each row
-- insert into oplog values(0, opr, "routeset", "insert", id, null);
-- !
-- create trigger oplog_routeset_update after update on routeset for each row
-- insert into oplog values(0, opr, "routeset", "update", id, null);
-- !
-- create trigger oplog_routeset_delete after delete on routeset for each row
-- insert into oplog values(0, opr, "routeset", "delete", id, null);
-- !
-- create trigger oplog_rr_insert after insert on rr for each row
-- insert into oplog values(0, opr, "rr", "insert", id, null);
-- !
-- create trigger oplog_rr_update after update on rr for each row
-- insert into oplog values(0, opr, "rr", "update", id, null);
-- !
-- create trigger oplog_rr_delete after delete on rr for each row
-- insert into oplog values(0, opr, "rr", "delete", id, null);
-- !
-- create trigger oplog_rrset_insert after insert on rrset for each row
-- insert into oplog values(0, opr, "rrset", "insert", id, null);
-- !
-- create trigger oplog_rrset_update after update on rrset for each row
-- insert into oplog values(0, opr, "rrset", "update", id, null);
-- !
-- create trigger oplog_rrset_delete after delete on rrset for each row
-- insert into oplog values(0, opr, "rrset", "delete", id, null);
-- !
-- create trigger oplog_viewer_insert after insert on viewer for each row
-- insert into oplog values(0, opr, "viewer", "insert", id, null);
-- !
-- create trigger oplog_viewer_update after update on viewer for each row
-- insert into oplog values(0, opr, "viewer", "update", id, null);
-- !
-- create trigger oplog_viewer_delete after delete on viewer for each row
-- insert into oplog values(0, opr, "viewer", "delete", id, null);
-- !
