package usecase

import (
	"context"

	"github.com/PacoXXD/p2p-monitor/pkg/models"
)

type MoitorUsecase interface {
	ReportPeer(ctx context.Context, ip, port, tracker_url, chat_url, share_key string, status models.PeerStatus) error
	ListPeer(ctx context.Context, share_key string) ([]models.Peer, error)
	GetPeer(ctx context.Context, share_key string) (*models.Peer, error)
}
