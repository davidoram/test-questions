package ring

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetAlwaysReturnsSameServer(t *testing.T) {
	r := NewRing("server-a", "server-b", "server-c", "server-d")

	keys := make([]string, 50)
	for i := range keys {
		keys[i] = fmt.Sprintf("key-%d", i)
	}

	first := make(map[string]string, len(keys))
	for _, k := range keys {
		first[k] = r.Get(k)
	}

	for range 10 {
		for _, k := range keys {
			require.Equal(t, first[k], r.Get(k), "key %q routed to different server", k)
		}
	}
}

func BenchmarkGet(b *testing.B) {
	const numKeys = 100_000
	r := NewRing("server-a", "server-b", "server-c", "server-d")

	keys := make([]string, numKeys)
	for i := range keys {
		keys[i] = fmt.Sprintf("key-%d", i)
	}

	b.ResetTimer()
	for i := range b.N {
		r.Get(keys[i%numKeys])
	}
}

func TestGetOnlyRoutesToKnownServers(t *testing.T) {
	servers := []string{"server-a", "server-b", "server-c", "server-d"}
	r := NewRing(servers...)

	valid := make(map[string]bool)
	for _, s := range servers {
		valid[s] = true
	}

	for i := range 200 {
		key := fmt.Sprintf("key-%d", i)
		got := r.Get(key)
		require.True(t, valid[got], "key %q routed to unknown server %q", key, got)
	}
}
