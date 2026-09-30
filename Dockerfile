FROM golang:1.26-bookworm AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o music-server .

FROM debian:bookworm-slim


RUN apt-get update \
    && apt-get install -y --no-install-recommends ffmpeg ca-certificates \
    && rm -rf /var/lib/apt/lists/*

COPY --from=denoland/deno:bin /deno /usr/local/bin/deno

WORKDIR /app

RUN mkdir -p /app/data

COPY --from=build /app/music-server ./music-server
COPY yt-dlp/yt-dlp_linux ./yt-dlp/yt-dlp_linux
RUN chmod +x ./yt-dlp/yt-dlp_linux

CMD ["./music-server"]
