FROM golang:1.25.3

ARG VERSION=v0.1.0
ARG BINARY=endoflife-${VERSION}

WORKDIR /app

# Copy the prebuilt binary
COPY ${BINARY} .

# Set permissions and create non-root user
RUN chmod +x ${BINARY} && \
    groupadd -r app && \
    useradd -r -g app app && \
    chown app:app ${BINARY}

USER app

# Use shell form to allow variable expansion
CMD ./endoflife-${VERSION}