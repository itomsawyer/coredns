drop view if exists gwlink_view;
drop view if exists gwlink_summary_view;
drop view if exists policy_view;
drop view if exists ipset_view;
drop view if exists domain_view;
drop view if exists netlink_view;
drop view if exists netlinkset_view;
drop view if exists route_view;
drop view if exists the_view;

drop table if exists viewer;
drop table if exists ipnet;
drop table if exists ipset;
drop table if exists domainlink;
drop table if exists domain;
drop table if exists domain_pool;
drop table if exists gwlink;
drop table if exists gwlinkset;
drop table if exists iptable;
drop table if exists netlink;
drop table if exists netlinkset;
drop table if exists gw;
drop table if exists policy_detail;
drop table if exists policy;
drop table if exists ldns;
drop table if exists rr;
drop table if exists rrset;



create table ipset (
id int(64) not null auto_increment,
name varchar(255) not null,
info text not null default "",
primary key (id),
unique key (name)
)DEFAULT CHARSET=utf8 comment "client ipnet set";

create table ipnet (
id int(64) not null auto_increment,
ip_start varchar(40) not null,
ip_end varchar(40) not null,
ipnet varchar(40) not null, 
mask int(8) not null,
priority int not null default 0,
ipset_id int(64),
primary key(id),
unique key(ip_start, ip_end , priority),
foreign key(ipset_id) references ipset(id) on delete cascade
)DEFAULT CHARSET=utf8 comment "client ipnet";

create table domain_pool (
id int(64) not null auto_increment,
name varchar(255) not null ,
info text not null default "",
primary key (id),
unique key (name)
)DEFAULT CHARSET=utf8 comment "serve domains set";

insert into domain_pool (name, info) values("global", "Base domain pool for all of domains which are not specifically configured");

create table domain (
id int(64) not null auto_increment,
domain varchar(255) not null,
domain_pool_id int(64),
primary key(id),
unique key(domain, domain_pool_id),
foreign key(domain_pool_id) references domain_pool(id) on delete cascade
)DEFAULT CHARSET=utf8 comment "serve domains";


create table gw (
id int(64) not null auto_increment,
name varchar(255) not null,
addr varchar(255) not null,
typ varchar(32) not null default "normal",
enable bool not null default true,
unavailable int(16) not null default 0 comment "if other than zero, gw is unavailable, each bit indicate different reason",
primary key(id),
unique key(name)
)DEFAULT CHARSET=utf8 comment "network gateway, aka outlink";

create table netlinkset(
id int(64) not null auto_increment,
name varchar(255) not null, 
primary key(id),
unique key(name)
)DEFAULT CHARSET=utf8 comment "netlink set ";

create table netlink (
id int(64) not null auto_increment,
isp  varchar(255) not null,
region varchar(255) not null,
typ varchar(32) not null default "normal",
primary key(id)
)DEFAULT CHARSET=utf8 comment "netlink (isp + province or CP) of a target ip";

create table iptable (
id int(64) not null auto_increment,
ip_start varchar(40) not null,
ip_end varchar(40) not null,
ipnet varchar(40) not null, 
mask int(8) not null,
priority int not null default 0,
netlink_id int(64),
primary key(id),
unique key(ip_start, ip_end, priority),
foreign key(netlink_id) references netlink(id) on delete restrict
)DEFAULT CHARSET=utf8 comment "IP to netlink";

create table domainlink(
id int(64) not null auto_increment,
domain_pool_id int(64) not null,
netlink_id int(64) not null, 
netlinkset_id int(64) not null,
enable bool not null default true,
primary key(id),
unique key(domain_pool_id, netlink_id),
foreign key(domain_pool_id) references domain_pool(id) on delete restrict,
foreign key(netlink_id) references netlink(id) on delete restrict,
foreign key(netlinkset_id) references netlinkset(id) on delete cascade
)DEFAULT CHARSET=utf8 comment "bind domain_pool and netlink to a netlinkset";

create table gwlinkset (
id int(64) not null auto_increment,
name varchar(255) not null,
info text,
primary key(id),
unique key(name)
)DEFAULT CHARSET=utf8;

create table gwlink (
id int(64) not null auto_increment,
gw_id int(64) not null,
netlinkset_id int(64) not null,
gwlinkset_id int(64) not null,
enable bool not null default true,
priority int not null default 0,
score int comment "netlink performance index",
unavailable int(16) not null default 0 comment "if other than zero, gwlink is unavailable, each bit indicate different reason",
primary key(id),
unique key(netlinkset_id, gw_id, gwlinkset_id),
foreign key(gw_id) references gw(id) on delete restrict,
foreign key(netlinkset_id) references netlinkset(id) on delete restrict,
foreign key(gwlinkset_id) references gwlinkset(id) on delete cascade 
)DEFAULT CHARSET=utf8 comment "Performance of using gateway to serve paricular netlink";


create table policy (
id int(64) not null auto_increment,
name varchar(255) not null,
primary key(id),
unique key(name)
)DEFAULT CHARSET=utf8 comment "policy index of choose ldns upstream forwarder";

create table ldns (
id int(64) not null auto_increment,
name varchar(255) not null,
addr varchar(255) not null,
typ varchar(16) not null default "A",
enable bool not null default true,
unavailable int(16) not null default 0 comment "if other than zero, ldns is unavailable with each bit indicate different reason",
primary key(id),
unique key(name)
)DEFAULT CHARSET=utf8 comment "upstream ldns info";


create table rrset(
id int(64) not null auto_increment,
name varchar(64) not null default "",
ttl int(32) unsigned not null default 300,
enable bool not null default true,
primary key(id)
)DEFAULT CHARSET=utf8 comment "dns rrset(resource record)";

create table rr (
id int(64) not null auto_increment,
rrtype int(16) unsigned not null,
rrdata varchar(255) not null,
ttl int(32) unsigned not null default 300,
rrset_id int(64),
enable bool not null default true,
primary key(id),
foreign key(rrset_id) references rrset(id) on delete cascade
)DEFAULT CHARSET=utf8 comment "dns rr(resource record)";

create table policy_detail (
id int(64) not null auto_increment,
policy_id int(64) not null,
policy_sequence int not null default 0,
enable bool not null default true,
priority int not null default 0,
weight int not null default 100,
op varchar(255) not null default "and",
op_typ varchar(64) not null default "builtin",
ldns_id int(64) not null,
rrset_id int(64),
primary key(id),
foreign key(ldns_id) references ldns(id) on delete restrict,
foreign key(rrset_id) references rrset(id) on delete restrict,
foreign key(policy_id) references policy(id) on delete cascade
)DEFAULT CHARSET=utf8 comment "policy detail of choose ldns upstream forwarder";

create table viewer(
id int(64) not null auto_increment,
ipset_id int(64) not null,
domain_pool_id int(64) not null,
gwlinkset_id int(64) not null, 
policy_id int(64) not null,
enable bool not null default true,
primary key(id),
foreign key(ipset_id) references ipset(id) on delete restrict,
foreign key(domain_pool_id) references domain_pool(id) on delete restrict,
foreign key(gwlinkset_id) references gwlinkset(id) on delete restrict,
foreign key(policy_id) references policy(id) on delete restrict
)DEFAULT CHARSET=utf8 comment "map of <ipset , domain_pool> -> <policy, gwlinkset>";

create ALGORITHM = MERGE view ipset_view as
select ipnet.id as ipnet_id, ip_start, ip_end, ipnet, mask, priority, ipset_id, name as ipset_name from ipnet join ipset on ipset.id = ipnet.ipset_id;

create ALGORITHM = MERGE view domain_view as
select domain.id as domain_id, domain, domain_pool_id, domain_pool.name as pool_name from domain join domain_pool on domain.domain_pool_id= domain_pool.id;

create ALGORITHM = MERGE view netlink_view as
select iptable.id as iptable_id, ip_start, ip_end, ipnet, mask, priority, netlink_id, isp, region, typ
from netlink join iptable on iptable.netlink_id = netlink.id;

create ALGORITHM = MERGE view netlinkset_view as
select domainlink.domain_pool_id, domainlink.netlink_id, domainlink.netlinkset_id, domain_pool.name as domain_pool_name, isp, region
from  domain_pool, netlink, domainlink
where domain_pool.id = domainlink.domain_pool_id and netlink.id = domainlink.netlink_id
and domainlink.enable = true;

-- create ALGORITHM = MERGE view gwlink_view as
-- select gwlinkset_id, gwlinkset.name as gwlinkset_name, netlinkset_id,  netlinkset.name as netlinkset_name, gwlink.id as gwlink_id, gwlink.enable as gwlink_enbale, gwlink.priority as gwlink_priority, gwlink.score as gwlink_score, gwlink.unavailable as gwlink_unavailable, gw_id, gw.name as gw_name, gw.addr as gw_addr, gw.typ as gw_typ, gw.enable as gw_enalble, gw.unavailable as gw_unavailable
-- from netlinkset, gw, gwlink, gwlinkset
-- where netlinkset.id = gwlink.netlinkset_id and gw.id = gwlink.gw_id and gwlink.gwlinkset_id = gwlinkset.id;

create ALGORITHM = MERGE view gwlink_view as
select gwlinkset_id, gwlinkset.name as gwlinkset_name,  netlinkset_id, netlinkset.name as netlinkset_name, gwlink.id as gwlink_id, min(gwlink.priority) as gwlink_priority, max(gwlink.score) as gwlink_score, gw_id, gw.name as gw_name,  gw.addr as gw_addr, gw.typ as gw_typ
from netlinkset, gw, gwlink, gwlinkset
where netlinkset.id = gwlink.netlinkset_id and gw.id = gwlink.gw_id and gwlink.gwlinkset_id = gwlinkset.id
and gw.enable = true and gwlink.enable = true and gw.unavailable = 0 and gwlink.unavailable = 0
group by gwlinkset_id;

create ALGORITHM = MERGE view route_view as
select
viewer.ipset_id,
ipset.name as ipset_name,
viewer.domain_pool_id,
domain_pool.name as domain_pool_name,
viewer.gwlinkset_id, gwlinkset.name as gwlinkset_name,
netlinkset_id, netlinkset.name as netlinkset_name,
gwlink.id as gwlink_id, min(gwlink.priority) as gwlink_priority, max(gwlink.score) as gwlink_score,
gw_id, gw.name as gw_name,  gw.addr as gw_addr, gw.typ as gw_typ
from ipset, domain_pool, viewer, gwlinkset, netlinkset, gwlink, gw
where
ipset.id = viewer.ipset_id and domain_pool.id = viewer.domain_pool_id and viewer.gwlinkset_id = gwlinkset.id
and netlinkset.id = gwlink.netlinkset_id and gw.id = gwlink.gw_id and gwlink.gwlinkset_id = gwlinkset.id
and gw.enable = true and gwlink.enable = true and gw.unavailable = 0 and gwlink.unavailable = 0 and viewer.enable = true
group by viewer.gwlinkset_id;



create ALGORITHM = MERGE view policy_view as
select viewer.ipset_id, viewer.domain_pool_id, viewer.policy_id, policy.name as policy_name,
policy_detail.policy_sequence, policy_detail.priority, policy_detail.weight, policy_detail.op, policy_detail.op_typ,
policy_detail.ldns_id, ldns.name, ldns.addr, ldns.typ,
policy_detail.rrset_id
from viewer, policy, policy_detail, ldns
where policy.id = policy_detail.policy_id and ldns.id = policy_detail.ldns_id and viewer.policy_id = policy.id
and ldns.unavailable = 0 and ldns.enable = true and viewer.enable = true and policy_detail.enable = true;

create ALGORITHM = MERGE view the_view as
select ipset_id, ipset.name as ipset_name, domain_pool_id, domain_pool.name as domain_pool_name, gwlinkset_id, gwlinkset.name as gwlinkset_name
from ipset, domain_pool, viewer, gwlinkset
where ipset.id = viewer.ipset_id and domain_pool.id = viewer.domain_pool_id and gwlinkset.id = viewer.gwlinkset_id
and viewer.enable = true

-- create ALGORITHM = MERGE view the_view as
-- select ipnet.id as ipnet_id, domain.domain, policy_detail.policy_sequence, policy_detail.op, policy_detail.op_typ, policy_detail.priority as policy_prior, ldns.addr as ldns_addr, netlink.isp, netlink.region, gw.name as gw_name
-- from ipnet, ipset, domain_pool, domain, netlink, gw, gwlink, gwlinkset, policy, policy_detail, ldns, viewer
-- where ipset.id = ipnet.ipset_id and domain.domain_pool_id= domain_pool.id and netlink.id = gwlink.netlink_id and gw.id = gwlink.gw_id and gwlink.gwlinkset_id = gwlinkset.id and policy.id = policy_detail.policy_id and ldns.id = policy_detail.ldns_id
-- and viewer.enable = true and policy_detail.enable = true and gwlink.enable = true and ldns.enable = true and gwlink.unavailable = 0 and ldns.unavailable = 0
