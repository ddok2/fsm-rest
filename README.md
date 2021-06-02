# fsm rest for test 
[![Go Report Card](https://goreportcard.com/badge/github.com/ddok2/fsm-rest)](https://goreportcard.com/report/github.com/ddok2/fsm-rest)

using state machines for test.

Just amuse myself.

## Configure
Update `common/config/config.yml` file
```yaml
version: '2'

tester:
  port: 8089
  operationMode: blockchain #exchange or blockchain
  txInterval: 1000 #interval to call
  remittanceAmount: 100
  remittanceFee: 0
  chargeAmount: 10000
  adminChargeAmount: 10000000000
  transactionSendCount: 2000

booster:
  addr: txbooster.nuriflex.com
  port: 8080

exchange:
  addr: dex
  port: 8080

members:
  count: 10
  prefix: caaa
  adminId: NURI-GENERAL
  trashId: NURI-TRASH

```

## Run
```shell
$ go mod vendor
...

$ go run main.go
INFO: 2021/06/01 19:46:41.540778 [Bot: enter_state_new  - id:  caaaa0]
INFO: 2021/06/01 19:46:44.069260 [Bot: enter_state_new  - id:  caaaa1]
INFO: 2021/06/01 19:46:47.036188 [Bot: enter_state_new  - id:  caaaa2]
INFO: 2021/06/01 19:46:49.650210 [Bot: enter_state_new  - id:  caaaa3]
INFO: 2021/06/01 19:46:51.655037 [Bot: enter_state_new  - id:  caaaa4]
...
```
