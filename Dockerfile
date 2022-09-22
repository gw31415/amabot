FROM golang:latest as builder
RUN apt update
RUN apt install -y libgs-dev
WORKDIR /builder
COPY . .
RUN go build -o /builder/amabot

FROM alpine:latest
RUN mkdir /lib64
RUN ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2
RUN apk add texlive-dvi ghostscript-dev
WORKDIR /ama
RUN adduser -h /ama -HDS ama
USER ama
COPY --from=builder /builder/amabot /ama/amabot
CMD [ "/ama/amabot" ]
