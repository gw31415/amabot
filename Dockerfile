FROM rust:slim as builder
WORKDIR /usr/src/app
COPY . .
RUN --mount=type=cache,target=/usr/local/cargo,from=rust:slim,source=/usr/local/cargo \
    --mount=type=cache,target=target \
    cargo build --release --features docker && mv ./target/release/amabot ./amabot

FROM gcr.io/distroless/cc-debian12
COPY --from=builder /usr/src/app/amabot /app/amabot
WORKDIR /app
ADD https://github.com/notofonts/noto-cjk/raw/727f898acdf7b100d308af8edf63c3953b626a1b/Serif/SubsetOTF/JP/NotoSerifJP-Regular.otf .
CMD ["/app/amabot"]
