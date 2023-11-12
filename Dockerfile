# Downloads font: 'Noto Serif CJK JP'
FROM alpine:latest as downloader
WORKDIR /noto-cjk
RUN apk add --no-cache unzip curl
RUN curl -O https://noto-website.storage.googleapis.com/pkgs/NotoSerifCJKjp-hinted.zip && \
    unzip NotoSerifCJKjp-hinted.zip

# Build
FROM rust:slim as builder
WORKDIR /usr/src/app
COPY . .
RUN --mount=type=cache,target=/usr/local/cargo,from=rust:slim,source=/usr/local/cargo \
    --mount=type=cache,target=target \
    cargo build --release --features docker && mv ./target/release/amabot ./amabot

# Final minimum image
FROM gcr.io/distroless/cc-debian12:latest
WORKDIR /app
COPY --from=downloader /noto-cjk/NotoSerifCJKjp-Regular.otf .
COPY --from=builder /usr/src/app/amabot .
CMD ["/app/amabot"]
