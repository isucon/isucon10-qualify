FROM rust:1.46

EXPOSE 1323

WORKDIR /usr/src/isuumo

COPY Cargo.toml Cargo.lock ./
RUN mkdir src && echo 'fn main() {}' > src/main.rs && cargo build --locked --release
COPY src ./src
RUN touch src/main.rs && cargo build --locked --frozen --release
CMD ["/usr/src/isuumo/target/release/isuumo"]
