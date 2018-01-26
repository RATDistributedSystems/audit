FROM scratch

COPY audit /app/
WORKDIR "/app"
EXPOSE 44444
CMD ["./audit"]
