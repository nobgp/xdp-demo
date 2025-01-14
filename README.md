# XDP Demo

This project demonstrates the usage of XDP (eXpress Data Path) with Go, showcasing how to implement high-performance packet processing using eBPF (extended Berkeley Packet Filter) technology. The demo attaches an XDP program to a network interface and provides real-time statistics about packet throughput.

## Overview

XDP allows you to run custom packet processing code directly in the network driver, before the kernel networking stack. This enables high-performance operations like packet filtering, counting, and manipulation at the earliest possible point in the software stack.

This demo:
- Attaches an XDP program to a network interface (default: eth0)
- Counts incoming packets and bytes
- Calculates and displays packets per second (PPS) and throughput statistics
- Demonstrates the integration between Go and eBPF/XDP

## Prerequisites

To build and run this demo, you need:

- Docker (recommended method)
- Or, if building locally:
  - Go 1.23 or later
  - Clang and LLVM
  - Linux headers
  - libbpf development files

The Docker-based build process automatically handles all dependencies.

## Building

### Using Make (Recommended)

```bash
# Build the Docker image
make build
```

### Using Docker Directly

```bash
# Build the Docker image
docker build -t xdp-demo .
```

## Running the Demo

### Using Make (Recommended)

```bash
# Run the demo (attaches to eth0 by default)
make demo

# For development/debugging
make shell
```

### Using Docker Directly

```bash
# Run the demo with required capabilities
docker run --rm -it --name xdp-demo \
  --cap-add SYS_RESOURCE --cap-add NET_ADMIN --cap-add BPF \
  xdp-demo

# Optionally specify a different interface
docker run --rm -it --name xdp-demo \
  --cap-add SYS_RESOURCE --cap-add NET_ADMIN --cap-add BPF \
  xdp-demo eth1
```

The program will display packet statistics every second, including:
- Packets per second (PPS)
- Network throughput (bits per second)

Press Ctrl+C to exit gracefully.

## References

1. [Official eBPF Documentation](https://docs.cilium.io/en/stable/bpf/)
2. [Linux Kernel Documentation - BPF and XDP Reference Guide](https://www.kernel.org/doc/html/latest/bpf/index.html)
3. [Cilium's eBPF Go Library](https://github.com/cilium/ebpf)
4. [XDP Project Documentation](https://github.com/xdp-project/xdp-project)
