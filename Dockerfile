FROM alpine:latest

# Copy caddy configuration
COPY ci-build/verbose /usr/bin/verbose
COPY credentials.json /usr/lib/verbose/credentials.json

ENTRYPOINT [ "/usr/bin/verbose", "--port=8000", "--credentials-file=/usr/lib/verbose/credentials.json", "--vocabulary-file=/usr/lib/verbose/vocabulary.json" ]
