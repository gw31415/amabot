FROM alpine:latest as downloader
WORKDIR /noto-cjk
RUN apk add --no-cache unzip curl
RUN curl -O https://noto-website.storage.googleapis.com/pkgs/NotoSerifCJKjp-hinted.zip && \
    unzip NotoSerifCJKjp-hinted.zip

FROM rust:slim as builder
WORKDIR /usr/src/app
COPY . .
RUN --mount=type=cache,target=/usr/local/cargo,from=rust:slim,source=/usr/local/cargo \
    --mount=type=cache,target=target \
    cargo build --release --features docker && mv ./target/release/amabot ./amabot

FROM gcr.io/distroless/cc-debian12:latest
WORKDIR /app
COPY --from=downloader /noto-cjk/NotoSerifCJKjp-Regular.otf /app/amabot
COPY --from=builder /usr/src/app/amabot /app/amabot
CMD ["/app/amabot"]
