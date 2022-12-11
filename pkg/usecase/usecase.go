package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/coocood/freecache"
	"github.com/PacoXXD/p2p-monitor/pkg/models"
	log "github.com/sirupsen/logrus"
)

type monitorUsecase struct {
	cache freecache.Cache
}

type PeerList struct {
	Peers []models.Peer `json:"peers"`
}

func NewUsecase(cache *freecache.Cache) MoitorUsecase {
	return &monitorUsecase{
		cache: *cache,
	}
}

func (m *monitorUsecase) ReportPeer(ctx context.Context, ip, port, tracker_url, chat_url, share_key string, status models.PeerStatus) error {
	if ip == "" || port == "" || share_key == "" {
		return fmt.Errorf("invalid param")
	}
	var newPeer = models.Peer{
		TrackerURL: tracker_url,
		ShareKey:  share_key,
		ChatUrl:    chat_url,
		IP:         ip,
		Port:       port,
		States:     status,
	}

	got, err := m.cache.Get([]byte(share_key))
	if err != nil {
		// not found

		//offline peer not save record
		if status == models.Offline {
			return nil
		}
		res := PeerList{}
		res.Peers = append(res.Peers, newPeer)
		re, err := json.Marshal(res)
		if err != nil {
			return fmt.Errorf("failed marshal,err:%s", err)
		}
		log.WithFields(log.Fields{"peer": string(re)}).Info("[ADD] peer")

		return m.cache.Set([]byte(share_key), re, -1)
	} else {
		res := PeerList{}
		if err := json.Unmarshal(got, &res); err != nil {
			return fmt.Errorf("failed unmarshal,err:%s", err)
		}

		var newpeers []models.Peer
		for _, v := range res.Peers {
			if v.IP == newPeer.IP {
				if newPeer.States == models.Offline {
					//remove offline peer from list
					continue
				}
			}
			// add online peer to list
			newpeers = append(newpeers, v)
		}
		res.Peers = newpeers

		if len(res.Peers) == 0 {
			m.cache.Del([]byte(share_key))
			return nil
		}

		re, err := json.Marshal(res)
		if err != nil {
			return fmt.Errorf("failed marshal,err:%s", err)
		}

		log.WithFields(log.Fields{"len": len(res.Peers)}).Info("[COUNT]current peers")

		return m.cache.Set([]byte(share_key), re, -1)
	}

}

func (m *monitorUsecase) ListPeer(ctx context.Context, share_key string) ([]models.Peer, error) {
	got, err := m.cache.Get([]byte(share_key))
	if err != nil {
		// not found
		return []models.Peer{}, nil
	}
	res := PeerList{}
	if err := json.Unmarshal(got, &res); err != nil {
		log.Errorf("failed unmarshal,err:%s", err)
	}
	return res.Peers, nil
}

func (m *monitorUsecase) GetPeer(ctx context.Context, share_key string) (*models.Peer, error) {
	if share_key == "" {
		return nil, fmt.Errorf("invalid share_key:%s", share_key)
	}
	got, err := m.cache.Get([]byte(share_key))
	if err != nil {
		// not found
		return nil, err
	}
	res := PeerList{}
	if err := json.Unmarshal(got, &res); err != nil {
		return nil, fmt.Errorf("failed unmarshal,err:%s", err)
	}
	peers := res.Peers
	log.WithFields(log.Fields{"peers_count": len(peers)}).Info("found peers")

	if len(peers) == 0 {
		return nil, nil
	}
	// random chose peer
	rand.Seed(time.Now().Unix())
	index := rand.Intn(len(peers))

	log.WithFields(log.Fields{"peer_index": index, "peer_ip": &peers[index].IP}).Info("chose peers")

	return &peers[index], nil
}
