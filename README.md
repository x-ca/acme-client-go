# acme-client-go
golang acme client for x-ca/pebble

## Use

```
# 1. install root x-ca

# 2. run ssl req
go mod vendor
go run main

# 3. test cert
docker run -it -d \
  -p 8443:443 \
  -v $(pwd)/examples/default.conf:/etc/nginx/conf.d/default.conf \
  -v $(pwd)/certs/xiexianbin.cn/xiexianbin.cn.crt:/etc/pki/nginx/server.crt \
  -v $(pwd)/certs/xiexianbin.cn/xiexianbin.cn.key:/etc/pki/nginx/private/server.key \
  nginx
```
