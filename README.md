migrate -source file://migrations/ -database=$DB_DSN up 1

migrate -source file://migrations/ -database=$DB_DSN down 1

migrate create -ext sql -seq -dir ./migrations name

export DB_DSN=postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable

## TODO

1. add logging to handlers (encoding to json returns error)
2. unified logging in storage package
3. set issuer for JWT via env variables
4. panic-recover middleware
5. add env variable (dev | prod)

##

comments

// the refresh token rotation mechanism is used.
//
// during the refresh, the user exchanges the old refresh token (R1) for new refresh (R2) and access (A) tokens.
// if at some point in time after the issuance of the token (Rx), a refresh request arrives with one of the
// ancestors of the token (Rx), we believe that there has been a leak and invalidate the entire token branch.