FROM cgr.dev/chainguard/go:latest AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/root/go/pkg/mod,sharing=locked \
  go mod download

# Copy source code

COPY . .

RUN --mount=type=cache,target=/root/go/pkg/mod,sharing=locked \
  --mount=type=cache,target=/root/.cache/go-build,sharing=locked \
  go build -x -v \
  -o /build/quantm \
  ./cmd/quantm

# Runtime Stage (with git)
FROM cgr.dev/chainguard/git:latest-glibc AS quantm

COPY --from=builder /build/quantm /bin/quantm

ENTRYPOINT ["/bin/quantm"]
CMD []
