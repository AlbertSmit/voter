# see https://github.com/deanbaker/react-go-heroku/blob/master/Dockerfile

# Build the Go API
FROM golang:latest AS builder

## We copy everything in the root directory
## into our /app directory
ADD . /app

# Move to working directory /build
WORKDIR /app/server

# add go modules lockfiles
RUN go mod download

# Copy the code into the container
COPY . .

# build the binary with all dependencies
RUN CGO_ENABLED=0 go build -o /main .

# Build the Svelte application
FROM node:14.15-alpine3.12 AS node_builder
COPY --from=builder /app/client ./
RUN yarn install && \
    yarn build

# Final stage build, this will be the container
# that we will deploy to production
FROM alpine:latest
COPY --from=builder /main ./
COPY --from=node_builder /build ./web
RUN chmod +x ./main

# The EXPOSE instruction does not actually publish the port. 
# It functions as a type of documentation between the person 
# who builds the image and the person who runs the container, 
# about which ports are intended to be published.
# — Docker — Dockerfile Reference
EXPOSE 8080

CMD ./main
