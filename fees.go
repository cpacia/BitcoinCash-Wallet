package bitcoincash

import (
	"golang.org/x/net/proxy"
	"net"
	"net/http"
	"time"
	"github.com/OpenBazaar/wallet-interface"
)

type httpClient interface {
	Get(string) (*http.Response, error)
}

type feeCache struct {
	fees        *Fees
	lastUpdated time.Time
}

type Fees struct {
	FastestFee  uint64
	HalfHourFee uint64
	HourFee     uint64
}

type FeeProvider struct {
	maxFee      uint64
	priorityFee uint64
	normalFee   uint64
	economicFee uint64

	httpClient httpClient

	cache *feeCache
}

func NewFeeProvider(maxFee, priorityFee, normalFee, economicFee uint64, proxy proxy.Dialer) *FeeProvider {
	fp := FeeProvider{
		maxFee:      maxFee,
		priorityFee: priorityFee,
		normalFee:   normalFee,
		economicFee: economicFee,
		cache:       new(feeCache),
	}
	dial := net.Dial
	if proxy != nil {
		dial = proxy.Dial
	}
	tbTransport := &http.Transport{Dial: dial}
	httpClient := &http.Client{Transport: tbTransport, Timeout: time.Second * 10}
	fp.httpClient = httpClient
	return &fp
}

func (fp *FeeProvider) GetFeePerByte(feeLevel wallet.FeeLevel) uint64 {
	defaultFee := func() uint64 {
		switch feeLevel {
		case wallet.PRIOIRTY:
			return fp.priorityFee
		case wallet.NORMAL:
			return fp.normalFee
		case wallet.ECONOMIC:
			return fp.economicFee
		case wallet.FEE_BUMP:
			return fp.priorityFee * 2
		default:
			return fp.normalFee
		}
	}
	return defaultFee()
}
