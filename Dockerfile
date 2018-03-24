FROM scratch

COPY audit config.json /app/
WORKDIR "/app"
EXPOSE 44443
CMD ["./audit"]
