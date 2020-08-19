FROM alpine
COPY main /app/
WORKDIR /app
CMD ["./main"]