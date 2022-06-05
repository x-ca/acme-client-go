#!/bin/bash
# x-ca/ca/certs/370bfb24-4fb4-11e8-a6a8-bb6671df87db/xiexianbin.cn.key -o test/certs/localhost/key.pem
# x-ca/ca/certs/370bfb24-4fb4-11e8-a6a8-bb6671df87db/xiexianbin.cn.crt -o test/certs/localhost/cert.pem

docker run \
  -e "PEBBLE_VA_NOSLEEP=1" \
  -e "PEBBLE_VA_ALWAYS_VALID=1" \
  -p 14000:14000 \
  -p 15000:15000 \
  -v $(pwd)/test:/test:rw \
  letsencrypt/pebble
