#!/usr/bin/env bash

error()
{
    echo "$@"
    return 1
}

[ -n "${GOPATH}" -a "${GOPATH}" != "" ] || error "GOPATH not exist" || exit 1

export PATH=$PATH:$GOPATH/bin
which govendor 2>/dev/null
if [ $? -ne 0 ];then
    echo "install govendor tool"
    go get -u github.com/kardianos/govendor
else
    echo "govendor tool installed yet"
fi

if ! [ -f vendor/vendor.json ];then
    echo "init vendor repertory"
    govendor init
else
    echo "vendor init yet"
fi

govendor sync