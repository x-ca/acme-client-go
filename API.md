# API

ACME-based CA server for [x-ca](https://github.com/x-ca/x-ca). support by [letsencrypt/boulder](https://github.com/letsencrypt/boulder)

## use

### Start Server

```
bash start.sh
```

- `PEBBLE_VA_ALWAYS_VALID=1` Skipping Validation

### python client

ref `test/chisel2.py`

```
python chisel2.py foo.com bar.com
```

### new order

```
POST /acme/new-order HTTP/1.1
Host: example.com
Content-Type: application/jose+json

{
    "protected": base64url({
        "alg": "ES256",
        "kid": "https://example.com/acme/acct/1",
        "nonce": "5XJ1L3lEkMG7tR6pA00clA",
        "url": "https://example.com/acme/new-order"
    }),
    "payload": base64url({
        "identifiers": [{"type:"dns","value":"example.com"}],
        "notBefore": "2016-01-01T00:00:00Z",
        "notAfter": "2016-01-08T00:00:00Z"
    }),
    "signature": "H6ZXtGjTZyUnPeKn...wEA4TklBdh3e454g"
}

HTTP/1.1 201 Created
Replay-Nonce: MYAuvOpaoIiywTezizk5vw
Location: https://example.com/acme/order/asdf

{
    "status": "pending",
    "expires": "2016-01-01T00:00:00Z",

    "notBefore": "2016-01-01T00:00:00Z",
    "notAfter": "2016-01-08T00:00:00Z",

    "identifiers": [
        { "type:"dns", "value":"example.com" },
        { "type:"dns", "value":"www.example.com" }
    ],

    "authorizations": [
        "https://example.com/acme/authz/1234",
        "https://example.com/acme/authz/2345"
    ],

    "finalize": "https://example.com/acme/order/asdf/finalize"
}
```

### order finalize

```
POST /acme/order/asdf/finalize HTTP/1.1
Host: example.com
Content-Type: application/jose+json

{
    "protected": base64url({
        "alg": "ES256",
        "kid": "https://example.com/acme/acct/1",
        "nonce": "MSF2j2nawWHPxxkE3ZJtKQ",
        "url": "https://example.com/acme/order/asdf/finalize"
    }),
    "payload": base64url({
        "csr": "5jNudRx6Ye4HzKEqT5...FS6aKdZeGsysoCo4H9P",
    }),
    "signature": "uOrUfIIk5RyQ...nw62Ay1cl6AB"
}

HTTP/1.1 200 Ok
Replay-Nonce: CGf81JWBsq8QyIgPCi9Q9X
Location: https://example.com/acme/order/asdf

{
    "status": "valid",
    "expires": "2016-01-01T00:00:00Z",

    "notBefore": "2016-01-01T00:00:00Z",
    "notAfter": "2016-01-08T00:00:00Z",
    
    "identifiers": [
        { "type":"dns", "value":"example.com" },
        { "type":"dns", "value":"www.example.com" }
    ],
    
    "authorizations": [
        "https://example.com/acme/authz/1234",
        "https://example.com/acme/authz/2345"
    ],
    
    "finalize": "https://example.com/acme/order/asdf/finalize",
    
    "certificate": "https://example.com/acme/cert/asdf"
}
```

status:
- "invalid": The certificate will not be issued.  Consider this order process abandoned.
- "pending": The server does not believe that the client has fulfilled the requirements.  Check the "authorizations" array for entries that are still pending.
- "processing": The server agrees that the requirements have been fulfilled, and is in the process of generating the certificate. Retry after the time given in the "Retry-After" header field of the response, if any.
- "valid": The server has issued the certificate and provisioned its URL to the "certificate" field of the order.

### new authz

```
POST /acme/new-authz HTTP/1.1
Host: example.com
Content-Type: application/jose+json

{
    "protected": base64url({
        "alg": "ES256",
        "kid": "https://example.com/acme/acct/1",
        "nonce": "uQpSjlRb4vQVCjVYAyyUWg",
        "url": "https://example.com/acme/new-authz"
    }),
    "payload": base64url({
        "identifier": {
            "type": "dns",
            "value": "example.net"
        }
    }),
    "signature": "nuSDISbWG8mMgE7H...QyVUL68yzf3Zawps"
}
```

### down cert

```
GET /acme/cert/asdf HTTP/1.1
Host: example.com
Accept: application/pkix-cert

HTTP/1.1 200 OK
Content-Type: application/pem-certificate-chain
Link: <https://example.com/acme/some-directory>;rel="index"

-----BEGIN CERTIFICATE-----
[End-entity certificate contents]
-----END CERTIFICATE-----
-----BEGIN CERTIFICATE-----
[Issuer certificate contents]
-----END CERTIFICATE-----
-----BEGIN CERTIFICATE-----
[Other certificate contents]
-----END CERTIFICATE-----
```
