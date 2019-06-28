#!/bin/bash

find_netlink() {

cs=$(mysql -uroot iwg -e "select isp from netlink_wl_view where ip_start_int<= inet_aton('$1') and inet_aton('$1') <= ip_end_int;")

if ! [[ -z $cs ]]; then
   echo $cs | awk '{print $2}'
   exit 0
fi

cs=$(mysql -uroot iwg -e "select isp from netlink_view where ip_start_int<= inet_aton('$1') and inet_aton('$1') <= ip_end_int;")

if [[ -z $cs ]]; then
   echo "unknown"
else
   echo $cs | awk '{print $2}'
fi

}


awk '

{resp=$16; split(resp, ips, ","); for (i in ips) cnt[ips[i]]+=1}END{for (i in cnt) print i, cnt[i]}

' | sort -n -k2 -r | while read line; do
ip=$(echo "$line" | awk '{print $1}')
cnt=$(echo "$line" | awk '{print $2}')
netlink=$(find_netlink $ip)
echo $ip $cnt $netlink
done
