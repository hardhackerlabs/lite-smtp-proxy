# lite-smtp-proxy
A lite smtp proxy that can be used to proxy SES.

## Build
`
go build .
`

## Setup

`SMTP_PROXY_UPSTREAM` is an environment variable that must be configured, and its value is the address of the destination mail server that needs to be forwarded.

```
SMTP_PROXY_UPSTREAM={ses}:587
```

If you need to use TLS to encrypt the network connection between the client and the proxy, then you need to provide a valid certificate.

```
SMTP_PROXY_CERT={path}/cert.pem
SMTP_PROXY_KEY={path}/key.pem
```

Use the `SMTP_PROXY_PORT` environment variable to set the listening port of the proxy server.

```
SMTP_PROXY_PORT=587
```

## Run

```
nohup lite-smtp-proxy > output.log 2>&1 &
```
