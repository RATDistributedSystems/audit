FROM scratch

COPY audit /app/
COPY config.json /app/
WORKDIR "/app"
EXPOSE 44444
CMD ["./audit"]
