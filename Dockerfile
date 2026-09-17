FROM golang:1.19-alpine AS builder

WORKDIR /app

COPY . .

RUN go build -mod=mod -tags server -o /out/server ./cmd/server

FROM alpine:3.18

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

COPY --from=builder /out/server ./server
COPY iptvhost ./iptvhost
COPY iptvhost6 ./iptvhost6

ENV PORT=8080
EXPOSE 8080

CMD ["./server"]