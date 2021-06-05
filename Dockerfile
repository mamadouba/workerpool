FROM alpine:latest 
RUN apk --no-cache add ca-certificates

WORKDIR /app 
COPY . .
RUN chmod +x ./app 
CMD ["./app"]