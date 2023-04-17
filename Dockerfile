# This is a multi build Dockerfile, it builds the server and places into into a minimal image

# Step 1: Build the app
FROM golang:1.20-alpine AS build

WORKDIR /app

# Copy all files
COPY . .

# Download the dependencies and build
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o chassit-on-repeat ./main.go


# Step 2: Copy to a minimal image
FROM alpine:latest

WORKDIR /app/
VOLUME [ "/files", "/config" ]
# Default port
EXPOSE 8080


# Copy from build images
COPY --from=build /app/chassit-on-repeat ./chassit-on-repeat

ENTRYPOINT ["./chassit-on-repeat"]