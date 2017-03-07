#/bin/bash

commit=$(git log| head -n 1 | awk '{printf $2}')
sed -e "s/__commit__/$commit/" coredns.spec  > coredns-$commit.spec

version=$(cat coredns.version)
sed -e  "s/__version__/$version/" coredns-$commit.spec > coredns-$commit-${version}.spec
rm coredns-$commit.spec 

which rpmbuild
if [[ $? != 0 ]]; then
    echo rpmbuild command not found
    exit -1
fi

topdir=$(rpmbuild --eval '${%_topdir}')
mkdir -p ${topdir}/{BUILD,BUILDROOT,SOURCES,RPMS,SRPMS,SPECS}
mv coredns-$commit-${version}.spec ${topdir}/SPECS

#rpmbuild -ba rpmbuild/SPECS/coredns-$commit.spec 
#
#if ! [[ $? == 0 ]] ; then
#    echo "error exists, abort"
#    exit -1
#fi

version=$(cat coredns.version | awk 'BEGIN{FS="."}{printf $1"."$2"."($3+1)}')

#find rpmbuild/RPMS  -print0  -iname *.rpm | xargs  -0  mv  .
