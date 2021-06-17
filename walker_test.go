package walker

import (
	"bytes"
	"encoding/hex"
	"io"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestGetMD5Response(t *testing.T) {
	var cases = map[string]struct {
		address string

		httpClientDoResponse *http.Response
		httpClientDoError    error

		wantHex   string
		wantError error
	}{
		"success": {
			address: "google.com",
			httpClientDoResponse: &http.Response{
				Body: io.NopCloser(bytes.NewBufferString(`hello!`)),
			},
			wantHex:   "5a8dd3ad0756a93ded72b823b19dd877",
			wantError: nil,
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			var (
				ctrl       = gomock.NewController(t)
				httpClient = NewMockHTTPClient(ctrl)
			)
			httpClient.EXPECT().
				Do(gomock.Any()).
				Return(c.httpClientDoResponse, c.httpClientDoError)

			hash, err := GetMD5Response(httpClient, c.address)
			hexhash, err := hex.DecodeString(c.wantHex)
			if err != nil {
				t.Error(err)
			}
			assert.ElementsMatch(t, hexhash, hash)
			if c.wantError == nil {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, c.wantError.Error())
			}
		})
	}
}
