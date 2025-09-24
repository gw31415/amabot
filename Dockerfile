# Downloads font: 'Noto Serif CJK JP'
FROM alpine:latest as downloader
WORKDIR /noto-cjk
RUN apk add --no-cache unzip curl
RUN curl -O https://noto-website-2.storage.googleapis.com/pkgs/NotoSerifCJKjp-hinted.zip && \
    unzip NotoSerifCJKjp-hinted.zip
WORKDIR /librusty
RUN curl -O https://github.com/denoland/rusty_v8/releases/download/v139.0.0/librusty_v8_release_x86_64-unknown-linux-gnu.a.gz

# Build
FROM rust:slim as builder
WORKDIR /usr/src/app
COPY . .
COPY --from=downloader /librusty/librusty_v8_release_x86_64-unknown-linux-gnu.a.gz /usr/local/cargo/.rusty_v8/https___github_com_denoland_rusty_v8_releases_download_v139_0_0_librusty_v8_release_x86_64_unknown_linux_gnu_a_gz
RUN --mount=type=cache,target=/usr/local/cargo,from=rust:slim,source=/usr/local/cargo \
    --mount=type=cache,target=target \
    cargo build --release --features docker && mv ./target/release/amabot ./amabot

# Final minimum image
FROM gcr.io/distroless/cc-debian12:latest
WORKDIR /app
COPY --from=downloader /noto-cjk/NotoSerifCJKjp-Regular.otf .
COPY --from=builder /usr/src/app/amabot .
CMD ["/app/amabot"]
