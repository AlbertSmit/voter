# see https://github.com/deanbaker/react-go-heroku/blob/master/Dockerfile
# and https://github.com/prisma/prisma-client-go/blob/master/docs/deploy.md#set-up-go-generate

# Build the Go API
FROM golang:latest AS builder
ADD . /app
# WORKDIR /app/server

# Move to working directory /build
WORKDIR /build

# add go modules lockfiles
# COPY go.mod go.sum ./
RUN go mod download

# prefetch the binaries, so that they will be cached and not downloaded on each change
RUN go run github.com/prisma/prisma-client-go prefetch

COPY . .

# generate the Prisma Client Go client
RUN go generate ./...

# build the binary with all dependencies
RUN go build -o /main .

# Build the Svelte application
# FROM node:14.15-alpine3.12 AS node_builder
# COPY --from=builder /app/client ./
# RUN yarn install && \
#     yarn build

# Final stage build, this will be the container
# that we will deploy to production
# FROM alpine:latest
# RUN apk --no-cache add ca-certificates
# COPY --from=builder /main ./
# COPY --from=node_builder /build ./web
# RUN chmod +x ./main
# EXPOSE 8080

# CMD ./main

# Move to /dist directory as the place for resulting binary folder
WORKDIR /dist

# Copy binary from build to main folder
RUN cp /build/main .

# Export necessary port
EXPOSE 8080

# Command to run when starting the container
CMD ["/dist/main"]