# syntax=docker/dockerfile:1

# # # # # # # # # # # # # # #
# ====== BUILD STAGE ====== #
# # # # # # # # # # # # # # #

FROM golang:alpine as builder
WORKDIR /go/src/
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

# # # # # # # # # # # # # # #
# ====== FINAL STAGE ====== #
# # # # # # # # # # # # # # #

FROM busybox
WORKDIR /root/
COPY --from=builder /go/src/templates ./templates
COPY --from=builder /go/src/app .
EXPOSE 8080
CMD ["./app"]

# docker build --tag shortlink -t shortlink:multistage .
# docker run --env-file .env -p 8080:8080 shortlink:multistage