FROM alpine:latest 
COPY app .
RUN apk --no-cache --update add ca-certificates
CMD ["./app"]
EXPOSE 8080