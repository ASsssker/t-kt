#!/bin/bash

go_url=https://go.dev/dl/go1.23.0.linux-amd64.tar.gz
user_profile="/home/${SUDO_USER}/.profile"

trap "echo Ошибка при установке; exit" ERR

if [ ! -d /usr/local/go ]; then
    filename="$(mktemp go.XXXXX.tar.gz)"
    
    echo $(wget -O "${filename}" "${go_url}")
    tar -C "/usr/local" -xzf "${filename}"
    rm -f "${filename}"
    echo "export PATH=\$PATH:/usr/local/go/bin" >> "${user_profile}"
    source "${user_profile}"
 
fi

/usr/local/go/bin/go build -o t-kt ./cmd/cli/main.go
cp t-kt /bin/t-kt

echo "Успешно установлено"

