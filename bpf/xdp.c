#include <linux/bpf.h>
#include <bpf/bpf_helpers.h>

struct stats_rec {
    __u64 rx_packets;
    __u64 rx_bytes;
};

struct {
    __uint(type, BPF_MAP_TYPE_ARRAY);
    __uint(max_entries, 1);
    __type(key, __u32);
    __type(value, struct stats_rec);
} packet_stats SEC(".maps");

SEC("xdp")
int xdp_stats(struct xdp_md *ctx)
{
    // Key is the index in the BPF array map (always 0 for this demo)
    // since we only track global statistics
    __u32 key = 0;

    // Look up the stats_rec structure in the packet_stats map
    // This structure stores our packet and byte counters
    struct stats_rec *rec = bpf_map_lookup_elem(&packet_stats, &key);

    // If map lookup fails, just pass the packet without updating statistics
    if (!rec) return XDP_PASS;

    // Get the length of the current packet
    __u64 length = bpf_xdp_get_buff_len(ctx);

    // Update packet statistics:
    // - Increment the total packet counter
    // - Add the current packet's length to the total bytes counter
    rec->rx_packets++;
    rec->rx_bytes += length;

    // Pass the packet to the normal network stack
    return XDP_PASS;
}

char LICENSE[] SEC("license") = "GPL";
