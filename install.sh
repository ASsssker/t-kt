#!/bin/bash

go_url=https://go.dev/dl/go1.23.0.linux-amd64.tar.gz
user_profile="/home/${SUDO_USER}/.profile"

trap "echo Ошибка при установке; exit" ERR

function install_go() {
    trap "echo Ошибка при установке; rm -f ${filename} exit" ERR
    filename="$(mktemp go.XXXXX.tar.gz)"
    
    echo $(wget -O "${filename}" "${go_url}")
    tar -C "/usr/local" -xzf "${filename}"
    rm -f "${filename}"
    echo "export PATH=\$PATH:/usr/local/go/bin" >> "${user_profile}"
}

function install_select() {
    if [[ ! "Y y" == *"$1"* ]]; then
        echo "Установка прервана"
        exit 0
    else
       install_go
    fi
}


if [ -d /usr/local/go ]; then
    user_go_version="$(/usr/local/go/bin/go version | awk '{print $3}' | cut -d 'o' -f 2)"

    # Если версия Go ниже необходимой
    if [[  $( echo "${user_go_version}" | awk -F '.' '{print $2}') -lt 24 ]]; then
        echo "В данный момент установлен Go версии ${user_go_version}, а для сборки проекта необходима версия не ниже 1.23.0"
        read -p "Обновить Go до версии 1.23.0?(Y/y-обновить):" is_install
        install_select $is_install
    fi
else
    echo "Необходимо устаноть Go 1.23.0"
    read -p "Установить Go 1.23.0?(Y/y-установить):" is_install
    install_select $is_install
fi

/usr/local/go/bin/go build -o t-kt ./cmd/cli/main.go
cp t-kt /bin/t-kt

echo "Успешно установлено"

