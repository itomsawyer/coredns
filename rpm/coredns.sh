#/bin/bash

commit=$(git log| head -n 1 | awk '{printf $2}')
sed -e "s/__commit__/$commit/" coredns.spec  > coredns-$commit.spec

version=$(cat coredns.version)
sed -i -e  "s/__version__/$version/" coredns-$commit.spec

topdir=$(rpm --eval '${%_topdir}')
mkdir -p ${topdir}/{BUILD,BUILDROOT,SOURCES,RPMS,SRPMS,SPECS}
mv coredns-$commit.spec ${topdir}/SPECS

#rpmbuild -ba rpmbuild/SPECS/coredns-$commit.spec 
#
#if ! [[ $? == 0 ]] ; then
#    echo "error exists, abort"
#    exit -1
#fi

#find rpmbuild/RPMS  -print0  -iname *.rpm | xargs  -0  mv  .
