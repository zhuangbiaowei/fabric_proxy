# fabric_proxy

`go build`

`./fabric_proxy invoke initLedger symbol name owner 1000`

返回值：

{"ID":"b2f06dbe1a26535fe5643cd886172cebcb5c44c772266d9fafc7324a6246b712","Nonce":"T3w491sMsOvwrclQxiaAS+xW5dc7McgX"}


`./fabric_proxy query balance symbol owner`

返回值：

"1000"
