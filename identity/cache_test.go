package identity

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

var file = "/tmp/cache.json"

func TestIdentityCache_StoreIdentity(t *testing.T) {
	identity := FromAddress("0x000000000000000000000000000000000000000A")
	cache := IdentityCache{
		File: file,
	}

	err := cache.StoreIdentity(identity)
	assert.Nil(t, err)
}

func TestIdentityCache_GetIdentity(t *testing.T) {
	identity := FromAddress("0x000000000000000000000000000000000000000A")
	cache := IdentityCache{
		File: file,
	}

	err := cache.StoreIdentity(identity)
	assert.Nil(t, err)
	id, err := cache.GetIdentity()

	assert.Equal(t, id, identity)
	assert.Nil(t, err)
}

func TestIdentityCache_cacheExists(t *testing.T) {
	identity := FromAddress("0x000000000000000000000000000000000000000A")
	cache := IdentityCache{
		File: file,
	}

	err := cache.StoreIdentity(identity)
	assert.Nil(t, err)

	assert.True(t, cache.cacheExists())

	_, err = os.Stat(file)
	assert.True(t, err == nil && !os.IsNotExist(err))
}
