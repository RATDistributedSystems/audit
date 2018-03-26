FROM scratch

COPY audit /app/
WORKDIR "/app"
EXPOSE 44443
CMD ["./audit"]
