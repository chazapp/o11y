# Auth

The authentication microservice of o11y.

Accessible via API Gateway HTTP Route.

## Generate JWT keys

```bash
$ go install github.com/go-jose/go-jose/v4/jose-util@latest
$ jose-util generate-key -use sig -alg RS256
# Rename keys to more usable: privateKey.json ; publicKey.json
```

## HTTP Routes

- /register

- /auth

- /jwks