package storage

import (
	"tasks-service-demo/internal/storage/channel"
	"tasks-service-demo/internal/storage/naive"
	"tasks-service-demo/internal/storage/shard"
	"testing"
)

func Test_InitMemoryStore(t *testing.T) {
	ResetStore()
	InitStore(naive.NewMemoryStore())
	store := GetStore()

	if _, ok := store.(*naive.MemoryStore); !ok {
		t.Error("Unexpected MemoryStore store init")
	}
}

func Test_InitShardStore(t *testing.T) {
	ResetStore()
	InitStore(shard.NewShardStore(4))
	store := GetStore()

	if _, ok := store.(*shard.ShardStore); !ok {
		t.Error("Unexpected ShardStore store init")
	}
}

func Test_InitChannelStore(t *testing.T) {
	ResetStore()
	InitStore(channel.NewChannelStore(4))
	store := GetStore()

	if _, ok := store.(*channel.ChannelStore); !ok {
		t.Error("Unexpected ChannelStore store init")
	}
}

func Test_InitShardPoolStore(t *testing.T) {
	ResetStore()
	InitStore(shard.NewShardStoreGopool(4))
	store := GetStore()

	if _, ok := store.(*shard.ShardStoreGopool); !ok {
		t.Error("Unexpected ShardStoreGopool store init")
	}
}

func Test_GetStore(t *testing.T) {
	ResetStore()
	InitStore(naive.NewMemoryStore())
	store := GetStore()

	if _, ok := store.(*naive.MemoryStore); !ok {
		t.Error("Unexpected store init")
	}
}
