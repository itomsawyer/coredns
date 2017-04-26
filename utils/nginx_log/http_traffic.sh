#httpfile=/home/log/nginx/access.log-`date +%Y%m%d`.gz
httpfile=/home/log/nginx/access.log
#tcpfile=/home/log/nginx/stream.log-`date +%Y%m%d`.gz
#httpstat=/home/tmp/`date +%Y%m%d`.http
#peakhttpstat=/home/tmp/`date +%Y%m%d`-19_22.http
#tcpstat=/home/tmp/`date +%Y%m%d`.tcp
peakstart=`date  +%d/%b/%Y:21:0`
peakend=`date  +%d/%b/%Y:21:0`
outlink=/home/log/nginx/outlink.conf.cszx-tproxy2

#zcat  $httpfile |  awk -F'"' '{print $1,$3,$NF} ' |awk '{print $7,$11,$(NF-2)} ' |sed 's/\// /' |awk '  {split($1,b,/:/) ;a[b[1]" "$(NF-2)]+=$NF}  END { for (i in a) print i,a[i]} ' |awk 'NR==FNR { if (/\[/) name=$1 ; if (/src/) a[$2]=name;} NR!=FNR { print a[$1],$0} '  $outlink  - |sort -k4nr  | head -n 1000 > $httpstat

#zcat  $httpfile |  sed -n '\#'$peakstart'#,\#'$peakend'#p' |  awk -F'"' '{print $1,$3,$NF} ' |awk '{print $7,$11,$(NF-2)} ' |sed 's/\// /' |awk '  {split($1,b,/:/) ;a[b[1]" "$(NF-2)]+=$NF}  END { for (i in a) print i,a[i]} ' |awk 'NR==FNR { if (/\[/) name=$1 ; if (/src/) a[$2]=name;} NR!=FNR { print a[$1],$0} '  $outlink - |sort -k4nr  | head -n 1000 > $peakhttpstat

#zcat $tcpfile |awk '{print $3,$2,$(NF-2) } ' |sed 's/\// /' | awk ' {  split($1,b,/:/);  a[b[1]" "$2]+=$3 } END { for (i in a) print i,a[i]} ' | awk 'NR==FNR { if (/\[/) name=$1 ; if (/src/) a[$2]=name;} NR!=FNR { print a[$1]":",$0} ' $outlink  -  | sort -k4nr |head -n 1000 > $tcpstat


cat  $httpfile |  sed -n '\#'$peakstart'#,\#'$peakend'#p' |  awk -F'"' '{print $1,$3,$NF} ' |awk '{print $7,$11,$(NF-2)} ' |sed 's/\// /' |awk '  {split($1,b,/:/) ;a[b[1]" "$(NF-2)]+=$NF}  END { for (i in a) print i,a[i]} ' |awk 'NR==FNR { if (/\[/) name=$1 ; if (/src/) a[$2]=name;} NR!=FNR { print a[$1],$0} '  $outlink - |sort -k4nr  > out
