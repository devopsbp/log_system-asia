FROM golang:latest as builder
ENV GO111MODULE=on

# build stage
WORKDIR /go/src/app/
COPY app ./
RUN go mod download
LABEL MULTISTAGE_CACHE=rds_log_transfer
RUN CGO_ENABLED=0 GOOS=linux go build -o rds_log_transfer .
RUN chmod 0744 rds_log_transfer

# running app stage
FROM gcr.io/distroless/static:debug as runner
WORKDIR /app/
COPY --from=builder --chown=nonroot /go/src/app/rds_log_transfer ./

ENTRYPOINT [ "/app/rds_log_transfer" ]
