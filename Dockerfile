# Step 1: Use official Golang image as the base image (adjusted to the actual Go version)
FROM golang:1.23-alpine

# Step 2: Set the working directory inside the container
WORKDIR /app

# Step 3: Copy go.mod and go.sum files to download dependencies
COPY go.mod . 
COPY go.sum .

# Step 4: Download the dependencies
RUN go mod download

RUN apk add --no-cache tzdata

# Step 6: ENV VARS
ENV TZ=UTC
ENV SERVER_HOST=0.0.0.0
ENV SERVER_PORT=8090
ENV GIN_RELEASE_MODE=false          
ENV DB_HOST=127.0.0.1
ENV DB_PORT=5433
ENV DB_USER=postgres
ENV DB_PASSWORD=postgres
ENV DB_NAME=fyc
ENV SSLMode=disable
ENV Prefix=fyc
ENV JWT_Secret=fyc4711
ENV ExpireTokenTime=12

#VALKEY CONFIG
ENV VALKEY_HOST=192.168.1.107
ENV VALKEY_PORT=6379
ENV VALKEY_CHANNEL=fyc_valkey

# this is the prefix passed in token like bearer token (if true)
ENV TokenPrefBackoffice=true
ENV TokenPref3rdParty=false

# this checker used to check token data in db 
ENV TokenCheck=false
ENV ExtraLog=false
ENV SaveXml=false
ENV SwaggerBasePath=/

# Admin Backoffice 
ENV USERNAME=admin
ENV PASSWORD=admin


COPY . .

# Step 8: Build the Go application
RUN go build -o bin/app .

# Step 9: Set the entry point for the container to run the application
ENTRYPOINT [ "/app/bin/app" ]
