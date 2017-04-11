echo iptable_wl
mysql -uroot iwg -e "select * from netlink_wl_view where ip_start_int<= inet_aton('$1') and inet_aton('$1') <= ip_end_int;"
echo iptable
mysql -uroot iwg -e "select * from netlink_view where ip_start_int<= inet_aton('$1') and inet_aton('$1') <= ip_end_int;"
