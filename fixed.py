#!/usr/bin/python

from json import dumps
from getpass import getpass
from os import system, remove

def main():
    code = '''package main

import (
    "C"
    "math/big"
    "strings"
)

//export SignTxWithFixedPrivKey
func SignTxWithFixedPrivKey(txJson *C.char) *C.char {
    privKey := strings.Join([]string{%s}, "")
    maxValue, _ := new(big.Int).SetString("%d", 10)
    rawTx, err := signTxWithPrivKey(C.GoString(txJson), privKey, maxValue)
    if err == nil {
        return C.CString(rawTx)
    } else {
        return C.CString("ERROR: " + err.Error())
    }
}''' % (dumps(list(getpass('Private Key (not echoed): ')))[1:-1], int(float(raw_input('Max value in ether (0 means no limit): ') or '0') * 1e18))

    path = 'fixed.go'
    with open(path, 'w') as file: file.write(code)
    system('go build -buildmode=c-shared -v -ldflags="-s -w" -o SignTx.so')
    remove(path)

main()
