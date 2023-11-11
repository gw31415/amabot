FROM rust:latest as builder

WORKDIR /usr/src/app
COPY . .
RUN --mount=type=cache,target=/usr/local/cargo,from=rust:latest,source=/usr/local/cargo \
    --mount=type=cache,target=target \
    cargo build --release && mv ./target/release/amabot ./amabot

FROM alpine:latest
RUN apk --no-cache add font-noto-cjk-extra
RUN adduser -S app
USER app
WORKDIR /app
COPY --from=builder /usr/src/app/amabot /app/amabot
CMD /app/amabot
