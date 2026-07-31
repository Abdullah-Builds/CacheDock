package cache

import (
	"container/list"
	"encoding/json"
	"os"
	"time"
)

func (c *Cache) Save(filename string) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	// convert internal map[string]*list.Element to map[string]Item for snapshot
	snapshotData := make(map[string]Item, len(c.data))
	for k, elem := range c.data {
		if elem == nil {
			continue
		}
		switch v := elem.Value.(type) {
		case *Entry:
			snapshotData[k] = v.Item
		case Item:
			snapshotData[k] = v
		case *Item:
			snapshotData[k] = *v
		}
	}

	snapshot := Snapshot{
		Version: 1,
		SavedAt: time.Now(),
		Data:    snapshotData,
	}

	data, err := json.MarshalIndent(snapshot, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filename, data, 0644)
}

func (c *Cache) Load(filename string) error {

	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	var snapshot Snapshot

	if err := json.Unmarshal(data, &snapshot); err != nil {
		return err
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	// Rebuild both the map and LRU list from the snapshot.
	newData := make(map[string]*list.Element, len(snapshot.Data))
	newLRU := list.New()
	for k, it := range snapshot.Data {
		if !it.ExpiresAt.IsZero() && !time.Now().Before(it.ExpiresAt) {
			continue
		}
		newData[k] = newLRU.PushFront(&Entry{Key: k, Item: it})
	}

	c.data = newData
	c.lru = newLRU
	c.usedMemory = 0
	for key, elem := range c.data {
		c.usedMemory += int64(len(key) + len(elem.Value.(*Entry).Item.Value))
	}
	return nil
}

func (c *Cache) StartAutoSave(filename string, interval time.Duration) {
	go func() {

		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for range ticker.C {
			c.Save(filename)
		}

	}()
}
