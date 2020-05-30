FROM photoprism/development as build

# Set up project directory
WORKDIR "/go/src/github.com/photoprism/photoprism"
COPY . .

# Build PhotoPrism
RUN make dep build-js install

# Same base image as photoprism/development
FROM debian:10-slim

# Set environment variables
ENV DEBIAN_FRONTEND noninteractive

# Configure apt-get
RUN echo 'Acquire::Retries "10";' > /etc/apt/apt.conf.d/80retry && \
    echo 'APT::Install-Recommends "false";' > /etc/apt/apt.conf.d/80recommends && \
    echo 'APT::Install-Suggests "false";' > /etc/apt/apt.conf.d/80suggests && \
    echo 'APT::Get::Assume-Yes "true";' > /etc/apt/apt.conf.d/80forceyes && \
    echo 'APT::Get::Fix-Missing "true";' > /etc/apt/apt.conf.d/80fixmissin

# Install additional distribution packages
RUN apt-get update && \
    apt-get install --no-install-recommends \
    ca-certificates \
    curl \
    exiftool \
    ffmpeg \
    gnupg2 \
    libheif-examples \
    tzdata

# Copy dependencies
COPY --from=build /usr/lib/libtensorflow.so /usr/lib/libtensorflow.so
COPY --from=build /usr/lib/libtensorflow_framework.so /usr/lib/libtensorflow_framework.so
RUN ldconfig

# Install RAW to JPEG converter
# Suppress an apt-key warning about standard out not being a terminal
ENV APT_KEY_DONT_WARN_ON_DANGEROUS_USAGE=true
RUN sh -c "echo 'deb http://download.opensuse.org/repositories/graphics:/darktable/Debian_10/ /' > /etc/apt/sources.list.d/graphics:darktable.list" && \
    curl -Ls https://download.opensuse.org/repositories/graphics:/darktable/Debian_10/Release.key | (OUT=$(apt-key add - 2>&1) || echo $OUT) && \
    apt-get update && \
    apt-get install darktable && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /photoprism

ENV TF_CPP_MIN_LOG_LEVEL 2
ENV PATH /photoprism/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin

ENV PHOTOPRISM_ORIGINALS_PATH /photoprism/originals
ENV PHOTOPRISM_IMPORT_PATH /photoprism/import
ENV PHOTOPRISM_TIDB_PATH /photoprism/database
ENV PHOTOPRISM_TEMP_PATH /photoprism/temp
ENV PHOTOPRISM_CACHE_PATH /photoprism/cache
ENV PHOTOPRISM_CONFIG_PATH /photoprism/config
ENV PHOTOPRISM_CONFIG_FILE /photoprism/config/photoprism.yml
ENV PHOTOPRISM_ASSETS_PATH /photoprism/assets
ENV PHOTOPRISM_RESOURCES_PATH /photoprism/assets/resources
ENV PHOTOPRISM_LOG_FILENAME /photoprism/photoprism.log
ENV PHOTOPRISM_PID_FILENAME /photoprism/photoprism.pid

# Copy files to /photoprism
COPY --from=build /root/.local/bin/photoprism /photoprism/bin/photoprism
COPY --from=build /root/.local/share/photoprism /photoprism/assets
COPY --from=build /root/.config/photoprism /photoprism/config

# Create directories & show photoprism version
RUN mkdir -p \
    /photoprism/originals \
    /photoprism/import \
    /photoprism/database \
    /photoprism/temp \
    /photoprism/cache && \
    chmod -R 777 /photoprism && \
    photoprism -v

# Expose http and database ports
EXPOSE 2342 2343 4000

# Run server
CMD photoprism start
