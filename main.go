package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/cilium/ebpf/link"
	"github.com/cilium/ebpf/rlimit"
)

//go:generate go run github.com/cilium/ebpf/cmd/bpf2go xdp bpf/xdp.c

func main() {
	// Remove the memory locking limit to allow eBPF program loading.
	// This is required because eBPF maps and programs need locked memory
	// to ensure they can be accessed quickly from kernel space.
	if err := rlimit.RemoveMemlock(); err != nil {
		log.Fatalf("removing memlock rlimit: %v", err)
	}

	// Load pre-compiled programs into the kernel
	objs := xdpObjects{}
	if err := loadXdpObjects(&objs, nil); err != nil {
		log.Fatalf("loading objects: %v", err)
	}
	defer objs.Close()

	// Get interface name
	name := "eth0"
	if len(os.Args) > 1 {
		name = os.Args[1]
	}

	// Get interface by name
	iface, err := net.InterfaceByName(name)
	if err != nil {
		log.Fatalf("getting interface: %v", err)
	}

	// Attach the compiled eBPF program to the specified network interface.
	// link.AttachXDP sets up the XDP hook so that incoming packets
	// will be processed by our eBPF program in the kernel before
	// they reach the network stack, enabling high-performance packet processing.
	l, err := link.AttachXDP(link.XDPOptions{
		Program:   objs.XdpStats,
		Interface: iface.Index,
	})
	if err != nil {
		log.Fatalf("attaching xdp: %v", err)
	}
	defer l.Close()

	log.Printf("XDP program attached to %s\n", iface.Name)

	// Print stats every second
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	// Catch SIGINT and SIGTERM
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	// Initialize previous stats record for delta calculations
	var prev xdpStatsRec
	
	// Main loop: periodically retrieve packet counters from the eBPF map (packet_stats)
	// and compute packets per second (pps) and throughput in bits per second (bps).
	// The statistics are calculated by comparing current values with previous measurements.
	for {
		select {
		case <-ticker.C:
			var key uint32 = 0
			var rec xdpStatsRec
			err := objs.PacketStats.Lookup(&key, &rec)
			if err != nil {
				log.Printf("map lookup: %v", err)
				continue
			}

			if prev.RxPackets != 0 {
				pps := rec.RxPackets - prev.RxPackets
				bps := (rec.RxBytes - prev.RxBytes) * 8
				log.Printf("PPS: %d, Throughput: %d bps\n", pps, bps)
			}

			prev = rec
		case <-sig:
			log.Println("Received signal, exiting...")
			return
		}
	}
}
