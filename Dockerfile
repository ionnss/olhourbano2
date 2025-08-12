# Golang image
FROM golang:1.24.4

# Install system dependencies for file processing
RUN apt-get update && apt-get install -y \
    ffmpeg \
    imagemagick \
    qpdf \
    exiftool \
    antiword \
    zip \
    unzip \
    pdftk \
    && rm -rf /var/lib/apt/lists/*

# Fix ImageMagick policy to allow PDF processing
RUN sed -i 's/<policy domain="coder" rights="none" pattern="PDF" \/>/<policy domain="coder" rights="read|write" pattern="PDF" \/>/' /etc/ImageMagick-6/policy.xml

# Work directory
WORKDIR /olhourbano2

# Copy go.mod abd go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download && go mod tidy

# Copy the wait-for-it script
COPY wait-for-it.sh /wait-for-it.sh

# Make the wait-for-it script executable
RUN chmod +x /wait-for-it.sh

# Copy the rest of the application code
COPY . .

# Ensure static directory exists and has the correct permissions
RUN mkdir -p /olhourbano2/static && chmod -R 755 /olhourbano2/static

# Compile the application with optimizations for production
RUN CGO_ENABLED=0 GOOS=linux go build -o /usr/local/bin/app_olhourbano2  -ldflags="-s -w" && \
    chmod +x /usr/local/bin/app_olhourbano2

# Expose the application port
EXPOSE 8080

# Command to run the application
CMD ["/wait-for-it.sh", "db:5432", "--", "/usr/local/bin/app_olhourbano2"]