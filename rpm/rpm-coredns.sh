#/bin/bash

echo "Generate spec file based on version"
commit=$(git log| head -n 1 | awk '{printf $2}')
sed -e "s/__commit__/$commit/" coredns.spec  > coredns-$commit.spec

version=$(cat coredns.version)
sed -e  "s/__version__/$version/" coredns-$commit.spec > coredns-$commit-${version}.spec
rm coredns-$commit.spec 

if ! [[ -f coredns-$commit-${version}.spec ]]; then
    echo Error: cannot generate spec file
    exit -1
fi


echo "Check rpmbuild command"
which rpmbuild > /dev/null
if [[ $? != 0 ]]; then
    echo rpmbuild command not found
    exit -1
fi

echo "Copy spec file"
topdir=$(rpm --eval '%_topdir')
mkdir -p ${topdir}/{BUILD,BUILDROOT,SOURCES,RPMS,SRPMS,SPECS}
echo "mv coredns-$commit-${version}.spec $topdir/SPECS"
mv coredns-$commit-${version}.spec ${topdir}/SPECS

echo "Packing repo to a tar file"
tar cf coredns.tar  ../../coredns --exclude=../../coredns/rpm/coredns.tar
gzip coredns.tar
mv coredns.tar.gz ${topdir}/SOURCES

echo "Building rpm"
rpmbuild -ba ${topdir}/SPECS/coredns-$commit-${version}.spec 
if ! [[ $? == 0 ]] ; then
    echo "error exists, abort"
    exit -1
fi

echo "Fetch all rpm to local dir"
find /root/rpmbuild/ -name *${version}*.rpm -print0 | xargs -0 -n1 -i%  mv % .

echo "Done"

version=$(cat coredns.version | awk 'BEGIN{FS="."}{printf $1"."$2"."($3+1)}')
