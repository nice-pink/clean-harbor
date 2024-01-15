ARG TARGET

FROM cgr.dev/chainguard/go:latest-dev AS builder
ARG TARGET=cleaner

# Info
LABEL org.opencontainers.image.authors="r@nice.pink"
LABEL org.opencontainers.image.source="https://github.com/nice-pink/clean-harbor/blob/main/Dockerfile"

WORKDIR /app

# get go module ready
COPY go.mod go.sum ./
RUN go mod download

# copy module code
COPY . .

# RUN CGO_ENABLED=0 GOOS=linux go build -o /repo
RUN chmod u+x ./build && ./build ${TARGET}

FROM cgr.dev/chainguard/git:latest-root-dev AS runner
ARG TARGET=cleaner

# add glibc compatibility
RUN apk add --update gcompat jq

# Info
LABEL org.opencontainers.image.authors="r@nice.pink"
LABEL org.opencontainers.image.source="https://github.com/nice-pink/clean-harbor/blob/main/Dockerfile"

WORKDIR /app

# copy executable
COPY --from=builder /app/${TARGET} /app/${TARGET}
ENTRYPOINT [ "/app/${TARGET}" ]
