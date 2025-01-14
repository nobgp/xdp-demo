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
    __u32 key = 0;

    struct stats_rec *rec = bpf_map_lookup_elem(&packet_stats, &key);

    if (!rec) return XDP_PASS;

    __u64 length = bpf_xdp_get_buff_len(ctx);

    rec->rx_packets++;
    rec->rx_bytes += length;

    return XDP_PASS;
}

char LICENSE[] SEC("license") = "GPL";
