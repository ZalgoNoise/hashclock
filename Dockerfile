# syntax=docker/dockerfile:1

FROM golang:alpine as builder

WORKDIR /app

COPY app/ ./

RUN go mod download && go build -o /hashclock

FROM scratch

COPY --from=builder /hashclock .

CMD [ "/hashclock" ]
ENTRYPOINT [ "/hashclock" ]