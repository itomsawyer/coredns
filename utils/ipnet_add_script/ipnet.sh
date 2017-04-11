if [[ -z $1 ]]; then
    echo "usage: $0 <clientset_id>"
    exit -1
fi

while read line; do
    mask=${line#*/}
    ipnet=${line%/*}
    echo "insert into ipnet values(0, \"$ipnet\", inet_ntoa(inet_aton(\"$ipnet\") + pow(2, 32 - $mask) -1), \"$ipnet\", $mask, $1);"
done
