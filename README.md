SignTx
======

Sign transaction with private key.

```
go build -buildmode=c-shared -v -ldflags="-s -w" -o SignTx.so
```

## Funcs

### `SignTxWithPrivKey(txJson, privKey string) (rawTx string)`

#### txJson

```json
{
    "nonce": 16,
    "gasPrice": 20000000000,
    "gas": 4000000,
    "to": "0x0000000000000000000000000000000000000000",
    "value": 1000000000000000000,
    "input": "SGVsbG8gV29ybGQh",
    "chainId": 123456
}
```

#### privKey

```
b0e565021eb427c1f71b810de4c8916ba478b5f6c9d339824c3e8b98f61762cb
```

#### rawTx (can be used in eth.sendRawTransaction)

```
0xf87c108504a817c800833d0900940000000000000000000000000000000000000000880de0b6b3
a76400008c48656c6c6f20576f726c64218303c4a4a0d956a764f81866d9b16d744cd0d513daeac0
a5831215a2d34ed9a0163b292462a05d78bfeaed709951333506f5043fdcb3b75ebecb4d35d64841
6149da37273f84
```

### `SignTxWithFixedPrivKey(txJson string) (rawTx string)`

Run `./fixed.py` to build with a fixed private key with max value per transaction.
