FROM alpine:3.21

RUN apk add --no-cache ffmpeg ca-certificates tzdata

ENV FFMPEG_BIN=/usr/bin/ffmpeg

WORKDIR /app
COPY comic2video /usr/local/bin/comic2video

VOLUME ["/app/app_data"]

EXPOSE 8080
ENTRYPOINT ["/usr/local/bin/comic2video"]
CMD ["serve"]
