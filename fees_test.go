package bitcoincash

import (
	"bytes"
	"net/http"
	"testing"
	"github.com/OpenBazaar/wallet-interface"
)

type ClosingBuffer struct {
	*bytes.Buffer
}

func (cb *ClosingBuffer) Close() (err error) {
	return
}

type mockHttpClient struct{}

func (m *mockHttpClient) Get(url string) (*http.Response, error) {
	data := `{"fastestFee":450,"halfHourFee":420,"hourFee":390}`
	cb := &ClosingBuffer{bytes.NewBufferString(data)}
	resp := &http.Response{
		Body: cb,
	}
	return resp, nil
}

func TestFeeProvider_GetFeePerByte(t *testing.T) {
	fp := NewFeeProvider(2000, 360, 320, 280, nil)

	// Test no API provided
	if fp.GetFeePerByte(wallet.PRIOIRTY) != 360 {
		t.Error("Returned incorrect fee per byte")
	}
	if fp.GetFeePerByte(wallet.NORMAL) != 320 {
		t.Error("Returned incorrect fee per byte")
	}
	if fp.GetFeePerByte(wallet.ECONOMIC) != 280 {
		t.Error("Returned incorrect fee per byte")
	}
	if fp.GetFeePerByte(wallet.FEE_BUMP) != 720 {
		t.Error("Returned incorrect fee per byte")
	}
}
