package peers

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMarshap(t *testing.T) {
	tests := map[string]struct {
		input  string
		output []Peer
		fails  bool
	}{
		"correctly parses peer": {
			input: string([]byte{127, 0, 0, 1, 0x00, 0x50, 1, 1, 1, 1, 0x01, 0xbb}),
			output: []Peer{
				{IP: net.IP{127, 0, 0, 1}, Port: 80},
				{IP: net.IP{1, 1, 1, 1}, Port: 443},
			},
		},
		"not enough bytes in peer": {
			input: string([]byte{127, 0, 0, 1, 0x00}),
			fails: true,
		},
	}
	for _, test := range tests {
		peers, err := Unmarshal([]byte(test.input))
		if test.fails {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err)
		}
		t.Logf("Peers: %+v", peers)
		assert.Equal(t, test.output, peers)
	}
}
