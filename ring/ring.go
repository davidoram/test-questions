package ring

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"sort"
)

// Map data to servers.
// When servers are added or removed, minimize the remapping needed.
// Consists of a ring 0 -> 2^32 -1
// Map servers and keys to teh same ring
// Hash key -> find position on ring. Walk clockwise until find the server
// Add server: Keys remapped between previous & new server
// Remove server: Move that servers keys to the next server
// Might map servers to multiple places on teh ring to distribute load more evenly
//
// Minimal remapping on topology changes
// Horizontally scaleable
// No central co-ordinator
// Even distrinution

// Hash algorithim needs to be fast, uniform distribution across the ring, low collision rate.
// SHA-256
type Ring struct {
	keys    []uint32          // sorted hash values Maybe faster with Balanced Binary tree - but for small numbers (< 20k?) not relevant
	hashMap map[uint32]string // hash → server name
}

func NewRing(servers ...string) *Ring {
	r := &Ring{hashMap: make(map[uint32]string)}
	for _, s := range servers {
		r.Add(s)
	}
	return r
}

const virtualNodes = 100

func hashKey(key string) uint32 {
	h := sha256.Sum256([]byte(key))
	return binary.BigEndian.Uint32(h[:4])
}

func (r *Ring) Add(server string) {
	for i := range virtualNodes {
		hash := hashKey(fmt.Sprintf("%s#%d", server, i))
		r.keys = append(r.keys, hash)
		r.hashMap[hash] = server
	}
	sort.Slice(r.keys, func(i, j int) bool { return r.keys[i] < r.keys[j] })
}

// Return the server that key is hashed to, or error if no servers
func (r *Ring) Get(key string) string {
	val := hashKey(key)
	i := sort.Search(len(r.keys), func(i int) bool {
		return r.keys[i] >= val
	})
	if i == len(r.keys) {
		i = 0 // wrap around
	}
	return r.hashMap[r.keys[i]]
}
