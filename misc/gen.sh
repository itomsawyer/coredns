#/bin/bash

echo "delimiter !" > trigger

while read tbl; do 

echo $tbl
echo "insert update delete" | xargs -n 1 echo | while read line; do 
cat >> trigger << EOF
create trigger oplog_${tbl}_${line} after $line on $tbl for each row
insert into oplog values(0, opr, "$tbl", "$line", id, null);
!
EOF

done 

done < t





