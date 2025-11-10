# syntax=docker/dockerfile:1

FROM alpine:3.22.2

# Metadata
LABEL maintainer="techprimate GmbH <opensource@techprimate.com>"
LABEL description="Container for github-actions-utils-cli"

# ARG for platform detection
ARG TARGETARCH

# Copy the appropriate binary based on target architecture
COPY dist/github-actions-utils-cli-linux-${TARGETARCH} /tmp/github-actions-utils-cli

# Install binary to PATH
RUN install \
    -o root \
    -g root \
    -m 0755 \
    /tmp/github-actions-utils-cli /usr/local/bin/github-actions-utils-cli && \
    rm -f /tmp/github-actions-utils-cli

# Smoke test
RUN set -x && \
    github-actions-utils-cli --version

# Set environment variables
ENV TZ=UTC

# Entrypoint
ENTRYPOINT ["/usr/local/bin/github-actions-utils-cli"]

