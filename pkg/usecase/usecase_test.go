package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/coocood/freecache"
	"github.com/PacoXXD/p2p-monitor/pkg/models"
	"github.com/stretchr/testify/assert"
)

func TestXxx(t *testing.T) {
	cacheSize := 100 * 1024 * 1024
	cache := freecache.NewCache(cacheSize)

	var share_key = "5DFA3FFA401A058A9C6E4518E92AB7FD"
	var peer_ip = "139.64.165.246"

	uc := NewUsecase(cache)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	err := uc.ReportPeer(ctx, peer_ip, "32676", "127.0.0.1", "127.0.0.1", share_key, models.Online)
	assert.NoError(t, err)

	peer, err := uc.GetPeer(ctx, share_key)
	assert.NoError(t, err)
	assert.NotEmpty(t, peer)
	// peer.ShareKey
	assert.Equal(t, share_key, peer.ShareKey)
	assert.Equal(t, peer_ip, peer.IP)

	err2 := uc.ReportPeer(ctx, peer_ip, "32676", "127.0.0.1", "127.0.0.1", share_key, models.Offline)
	assert.NoError(t, err2)

	peer3, _ := uc.GetPeer(ctx, share_key)
	assert.Empty(t, peer3)
}
