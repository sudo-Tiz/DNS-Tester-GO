# Dockerfile for Docusaurus documentation (dev only)
# Production build handled by GitHub Actions

FROM node:24-alpine

WORKDIR /app

# Copy package files first (better cache)
COPY package.json package-lock.json* ./

# Install dependencies
# If package-lock.json exists → use npm ci (secure, CI/CD)
# If not → use npm install (dev, generates lock)
RUN if [ -f package-lock.json ]; then \
      npm ci --prefer-offline --no-audit; \
    else \
      npm install --prefer-offline --no-audit; \
    fi

# Copy source files
COPY . .

# Expose port
EXPOSE 3000

# Health check
HEALTHCHECK --interval=30s --timeout=5s --start-period=30s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:3000 || exit 1

# Dev server with hot reload
CMD ["npm", "start", "--", "--host", "0.0.0.0", "--poll", "1000"]
