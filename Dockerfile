FROM rust:slim as builder
WORKDIR /usr/src/app
COPY . .
RUN --mount=type=cache,target=/usr/local/cargo,from=rust:slim,source=/usr/local/cargo \
    --mount=type=cache,target=target \
    cargo build --release && mv ./target/release/amabot ./amabot

FROM debian:stable-slim
RUN apt-get update \
 && apt-get install -y --no-install-recommends \
    fontconfig \
    fonts-noto-cjk-extra \
 && apt-get -y clean \
 && rm -rf /var/lib/apt/lists/* \
 && useradd app
USER app
COPY --from=builder /usr/src/app/amabot /home/app/amabot
WORKDIR /home/app
CMD ["/home/app/amabot"]
