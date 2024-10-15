#!/bin/bash

# Nama file Go dan binary
GO_FILE="main.go"
BINARY_NAME="proxy-downloader"

# Fungsi untuk mengecek apakah Go sudah terinstal
check_go() {
    if ! command -v go &> /dev/null; then
        echo "Go tidak terinstal. Menginstal Go..."
        install_go
    else
        echo "Go sudah terinstal."
    fi
}

# Fungsi untuk menginstal Go
install_go() {
    OS=$(uname -s | tr '[:upper:]' '[:lower:]')
    ARCH=$(uname -m)

    if [ "$ARCH" == "x86_64" ]; then
        ARCH="amd64"
    elif [[ "$ARCH" == "aarch64" || "$ARCH" == "arm64" ]]; then
        ARCH="arm64"
    else
        echo "Arsitektur $ARCH tidak didukung."
        exit 1
    fi

    GO_VERSION="1.21.1"
    FILE="go${GO_VERSION}.${OS}-${ARCH}.tar.gz"
    URL="https://golang.org/dl/${FILE}"

    echo "Mengunduh Go dari $URL..."
    curl -OL $URL

    echo "Ekstrak file..."
    sudo tar -C /usr/local -xzf $FILE
    rm -f $FILE

    export PATH=$PATH:/usr/local/go/bin
    echo "export PATH=\$PATH:/usr/local/go/bin" >> ~/.bashrc
    source ~/.bashrc

    echo "Go telah diinstal."
}

# Fungsi untuk mengompilasi kode Go
compile_go() {
    if [ ! -f "$GO_FILE" ]; then
        echo "File $GO_FILE tidak ditemukan!"
        exit 1
    fi

    echo "Mengompilasi $GO_FILE menjadi binary..."
    go build -o $BINARY_NAME $GO_FILE

    if [ $? -eq 0 ]; then
        echo "Kompilasi berhasil. Binary: $BINARY_NAME"
        echo "Menghapus $GO_FILE..."
        rm -f $GO_FILE
    else
        echo "Kompilasi gagal."
        exit 1
    fi
}

# Eksekusi proses instalasi
check_go
compile_go

echo "Instalasi selesai. Jalankan ./$BINARY_NAME untuk menggunakan program."
