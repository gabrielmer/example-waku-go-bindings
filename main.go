package main

import (
	"context"
	"fmt"
	"time"

	"github.com/waku-org/waku-go-bindings/waku"
	"go.uber.org/zap"
)

func main() {
	// Create logger
	logger, err := zap.NewDevelopment()
	if err != nil {
		fmt.Printf("Failed to create logger: %v\n", err)
		return
	}
	defer logger.Sync()

	const requestTimeout = 30 * time.Second
	// Configure dialer node
	dialerNodeWakuConfig := waku.WakuConfig{
		Relay:           true,
		LogLevel:        "DEBUG",
		Discv5Discovery: false,
		ClusterID:       16,
		Shards:          []uint16{64},
		Discv5UdpPort:   9020,
		TcpPort:         60020,
	}

	// Create and start dialer node
	dialerNode, err := waku.NewWakuNode(&dialerNodeWakuConfig, logger.Named("dialerNode"))
	if err != nil {
		fmt.Printf("Failed to create dialer node: %v\n", err)
		return
	}
	if err := dialerNode.Start(); err != nil {
		fmt.Printf("Failed to start dialer node: %v\n", err)
		return
	}
	defer dialerNode.Stop()
	time.Sleep(1 * time.Second)

	// Configure receiver node
	receiverNodeWakuConfig := waku.WakuConfig{
		Relay:           true,
		LogLevel:        "DEBUG",
		Discv5Discovery: false,
		ClusterID:       16,
		Shards:          []uint16{64},
		Discv5UdpPort:   9021,
		TcpPort:         60021,
	}

	// Create and start receiver node
	receiverNode, err := waku.NewWakuNode(&receiverNodeWakuConfig, logger.Named("receiverNode"))
	if err != nil {
		fmt.Printf("Failed to create receiver node: %v\n", err)
		return
	}
	if err := receiverNode.Start(); err != nil {
		fmt.Printf("Failed to start receiver node: %v\n", err)
		return
	}
	defer receiverNode.Stop()
	time.Sleep(1 * time.Second)

	// Get receiver node's multiaddress
	receiverMultiaddr, err := receiverNode.ListenAddresses()
	if err != nil {
		fmt.Printf("Failed to get receiver node addresses: %v\n", err)
		return
	}

	// Check initial peer counts
	dialerPeerCount, err := dialerNode.GetNumConnectedPeers()
	if err != nil {
		fmt.Printf("Failed to get dialer peer count: %v\n", err)
		return
	}
	fmt.Printf("Dialer initial peer count: %d\n", dialerPeerCount)

	receiverPeerCount, err := receiverNode.GetNumConnectedPeers()
	if err != nil {
		fmt.Printf("Failed to get receiver peer count: %v\n", err)
		return
	}
	fmt.Printf("Receiver initial peer count: %d\n", receiverPeerCount)

	// Dial peer
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()
	if err := dialerNode.Connect(ctx, receiverMultiaddr[0]); err != nil {
		fmt.Printf("Failed to dial peer: %v\n", err)
		return
	}
	time.Sleep(1 * time.Second)

	// Check final peer counts
	dialerPeerCount, err = dialerNode.GetNumConnectedPeers()
	if err != nil {
		fmt.Printf("Failed to get dialer peer count: %v\n", err)
		return
	}
	fmt.Printf("Dialer final peer count: %d\n", dialerPeerCount)

	receiverPeerCount, err = receiverNode.GetNumConnectedPeers()
	if err != nil {
		fmt.Printf("Failed to get receiver peer count: %v\n", err)
		return
	}
	fmt.Printf("Receiver final peer count: %d\n", receiverPeerCount)
}
