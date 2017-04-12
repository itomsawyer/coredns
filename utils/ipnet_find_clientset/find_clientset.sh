echo clientset_wl 
mysql -uroot iwg -e "select * from clientset_wl_view where ip_start_int<= inet_aton('$1') and inet_aton('$1') <= ip_end_int;"
echo clientset
mysql -uroot iwg -e "select * from clientset_view where ip_start_int<= inet_aton('$1') and inet_aton('$1') <= ip_end_int;"
