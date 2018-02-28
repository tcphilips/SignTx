#!/usr/bin/python

from json import dumps
from getpass import getpass
from os import system, remove

def main():
    code = '''package main

import (
    "C"
    "strings"
)

//export SignTxWithFixedPrivKey
func SignTxWithFixedPrivKey(txJson *C.char) *C.char {
    privKey := strings.Join([]string{%s}, "")
    rawTx, err := signTxWithPrivKey(C.GoString(txJson), privKey)
    if err == nil {
        return C.CString(rawTx)
    } else {
        return C.CString("ERROR: " + err.Error())
    }
}''' % dumps(list(getpass('Private Key (not echoed): ')))[1:-1]

    path = 'fixed.go'
    with open(path, 'w') as file: file.write(code)
    system('go build -buildmode=c-shared -v -ldflags="-s -w" -o SignTx.so')
    remove(path)

main()
