FROM golang:latest as builder
WORKDIR /builder
COPY . .
RUN go build -o /builder/amabot

FROM alpine:latest
RUN mkdir /lib64
RUN ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2
WORKDIR /ama/
COPY --from=builder /builder/amabot /ama/amabot
CMD [ "/ama/amabot" ]
