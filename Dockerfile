# Multi-stage build for opentrans

# UI builder stage
FROM node:20-alpine AS ui-builder

# Set working directory
WORKDIR /app

# Copy web directory for building
COPY web/ ./web/

# Create backend/webui/dist directory for build output
RUN mkdir -p backend/webui/dist

# Install dependencies and build
# Vite config outputs to ../backend/webui/dist
RUN cd web && npm install && npm run build

# Builder stage
FROM golang:1.24-alpine AS builder

RUN apk update && apk add --no-cache make gcc musl-dev git

# Set working directory to backend where go.mod is located
WORKDIR /app

COPY backend/ .

# Install Go dependencies
RUN go env -w GOPROXY='https://goproxy.cn,direct' && go mod tidy

# Install yzma CLI tool to download llama.cpp libraries
RUN go install github.com/hybridgroup/yzma/cmd/yzma@latest

# Download llama.cpp precompiled libraries for Linux (will be placed in ./lib by default)
RUN yzma lib get || true

# Ensure lib directory exists (even if yzma lib get fails)
RUN mkdir -p /app/lib

# Copy built UI assets from ui-builder stage
COPY --from=ui-builder /app/backend/webui/dist ./webui/dist

# Build the application using Makefile
RUN make

# Final stage - minimal Alpine image
FROM alpine:latest

# Install runtime dependencies including libyzma support
RUN apk --no-cache add \
    ca-certificates \
    tzdata \
    libstdc++ \
    libgcc \
    libffi \
    && update-ca-certificates

# Create non-root user
RUN addgroup -g 65532 --system nonroot && \
    adduser -D -u 65532 -G nonroot --system nonroot

# Copy the binary from builder stage
COPY --from=builder /app/opentrans /usr/local/bin/opentrans

# Copy yzma/llama.cpp libraries from builder stage
COPY --from=builder --chown=nonroot:nonroot /app/lib /usr/local/lib/yzma

# Create app directory, models directory for local Llama models, library directory, and set permissions
RUN mkdir -p /app /app/models /usr/local/lib/yzma && \
    chown -R nonroot:nonroot /app && \
    chmod -R 755 /usr/local/lib/yzma

# Set yzma library path
ENV LD_LIBRARY_PATH=/usr/local/lib/yzma
ENV YZMA_LIB=/usr/local/lib/yzma

WORKDIR /app

# Make binary executable
RUN chmod +x /usr/local/bin/opentrans

# Switch to non-root user
USER nonroot:nonroot

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD nc -z localhost 8080 || exit 1

# Default command
ENTRYPOINT ["opentrans"]
CMD ["serve"]

# Expose default port
EXPOSE 8080
