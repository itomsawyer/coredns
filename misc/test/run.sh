#!/bin/bash
user=root
pass="-proot"
db=iwg

mysql -u$user $pass $db < clientset.sql
mysql -u$user $pass $db < ipnet.sql
mysql -u$user $pass $db < netlink.sql
mysql -u$user $pass $db < iptable.sql
mysql -u$user $pass $db < domain_pool.sql
mysql -u$user $pass $db < domain.sql
mysql -u$user $pass $db < netlinkset.sql
mysql -u$user $pass $db < domainlink.sql
mysql -u$user $pass $db < outlink.sql
mysql -u$user $pass $db < routeset.sql
mysql -u$user $pass $db < route.sql
mysql -u$user $pass $db < ldns.sql
mysql -u$user $pass $db < policy.sql
mysql -u$user $pass $db < policy_detail.sql
mysql -u$user $pass $db < viewer.sql 
mysql -u$user $pass $db < filter.sql
mysql -u$user $pass $db < rrset.sql
mysql -u$user $pass $db < rr.sql

