FROM rust:slim as builder
WORKDIR /usr/src/app
COPY . .
RUN --mount=type=cache,target=/usr/local/cargo,from=rust:slim,source=/usr/local/cargo \
    --mount=type=cache,target=target \
    cargo build --release && mv ./target/release/amabot ./amabot

FROM debian:stable-slim
RUN useradd -m app
USER app
COPY --from=builder /usr/src/app/amabot /app/amabot
WORKDIR /app
CMD ["/app/amabot"]
