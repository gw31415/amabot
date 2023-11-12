FROM alpine:latest as downloader
WORKDIR /noto-cjk
RUN apk add --no-cache git
RUN --mount=type=cache,target=/noto-cjk/.git \
    git init || \
    git remote add origin https://github.com/notofonts/noto-cjk || \
    git fetch --depth 1 origin 727f898acdf7b100d308af8edf63c3953b626a1b || \
    git reset --hard FETCH_HEAD

FROM rust:slim as builder
WORKDIR /usr/src/app
COPY . .
RUN --mount=type=cache,target=/usr/local/cargo,from=rust:slim,source=/usr/local/cargo \
    --mount=type=cache,target=target \
    cargo build --release --features docker && mv ./target/release/amabot ./amabot

FROM gcr.io/distroless/cc-debian12:latest
WORKDIR /app
COPY --from=downloader /noto-cjk/Serif/OTF/Japanese/NotoSerifCJKjp-Regular.otf /app/amabot
COPY --from=builder /usr/src/app/amabot /app/amabot
CMD ["/app/amabot"]
