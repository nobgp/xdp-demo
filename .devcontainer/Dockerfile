FROM golang:bookworm

# Install dependencies and tools
RUN apt-get update && apt-get install -y --no-install-recommends \
    build-essential \
    clang \
    linux-libc-dev \
    libbpf-dev \
    llvm \
    wget \
    xz-utils \
    rsync \
    iproute2 \
    net-tools \
    vim \
    tcpdump \
    bpftool && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

# Download, extract, and install kernel headers in a single layer
RUN KERNEL_VERSION=$(uname -r | grep -oP '^\d+\.\d+(\.[1-9]\d*)?') && \
    MAJOR_VERSION=${KERNEL_VERSION%%.*} && \
    ARCH=$(uname -m) && \
    case "$ARCH" in \
        x86_64) ARCH="x86" ;; \
        aarch64) ARCH="arm64" ;; \
    esac && \
    wget https://cdn.kernel.org/pub/linux/kernel/v${MAJOR_VERSION}.x/linux-${KERNEL_VERSION}.tar.xz -O /tmp/linux-${KERNEL_VERSION}.tar.xz && \
    mkdir -p /usr/src/linux-${KERNEL_VERSION} && \
    tar -xf /tmp/linux-${KERNEL_VERSION}.tar.xz -C /usr/src/linux-${KERNEL_VERSION} --strip-components=1 && \
    rm /tmp/linux-${KERNEL_VERSION}.tar.xz && \
    cd /usr/src/linux-${KERNEL_VERSION} && \
    make ARCH=${ARCH} headers_install INSTALL_HDR_PATH=/usr/local && \
    rm -rf /usr/src/linux-${KERNEL_VERSION}
