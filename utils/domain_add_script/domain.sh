if [[ -z $1 ]]; then
    echo "usage: $0 <domain_pool_id>"
    exit -1
fi

while read line; do
    echo "insert into domain values(0, \"$line\", $1);"
done
