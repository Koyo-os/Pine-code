package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/peerstore"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/p2p/discovery/mdns"

	pubsub "github.com/libp2p/go-libp2p-pubsub"
)

type (
	Request struct {
		From    string          `json:"from"`
		Payload json.RawMessage `json:"payload"`
	}

	NodeType string

	Node struct {
		host   host.Host
		PubSub *pubsub.PubSub
		Topic  *pubsub.Topic
		Sub    *pubsub.Subscription
		ctx    context.Context
		Peers  map[peer.ID]*peer.AddrInfo
		mux    *sync.RWMutex
	}

	discoveryNotifee struct {
		Host host.Host
	}
)

func (n *discoveryNotifee) HandlePeerFound(pi peer.AddrInfo) {
	n.Host.Peerstore().AddAddrs(pi.ID, pi.Addrs, peerstore.PermanentAddrTTL)
}

func newNode(ctx context.Context, port int, topicID string) (*Node, error) {
	h, err := libp2p.New(
		libp2p.ListenAddrStrings(fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", port)),
	)
	if err != nil {
		return nil, err
	}

	ps, err := pubsub.NewGossipSub(ctx, h)
	if err != nil {
		return nil, err
	}

	topic, err := ps.Join(topicID)
	if err != nil {
		return nil, err
	}
	sub, err := topic.Subscribe()
	if err != nil {
		return nil, err
	}

	node := &Node{
		host:   h,
		PubSub: ps,
		Topic:  topic,
		Sub:    sub,
		ctx:    ctx,
		Peers:  make(map[peer.ID]*peer.AddrInfo),
	}

	mdnsService := mdns.NewMdnsService(h, topicID, &discoveryNotifee{
		Host: h,
	})

	if err = mdnsService.Start(); err != nil {
		return nil, err
	}

	return node, nil
}
