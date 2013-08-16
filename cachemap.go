package cachemap

import (
    "time"
    "log"
    "errors"
    "sync"
    "runtime"
)

/* NOTE: to use CacheMap, another object must implement CacheMapCallbacks */

type Callbacks interface {
    // persist new key, value somewhere else
    Create (key interface{}, item *Item) error

    // get value from key somewhere else
    Read   (key interface{}) (item *Item, err error)

    // persist updated key, value somewhere else
    Update (key interface{}, item *Item) error

    // delete key somewhere else
    Delete (key interface{}) error
}

// dont create this struct or modify these fields yourself
type Item struct {
    Value interface{}

    // when this item was originally Add()ed
    Created *time.Time

    // last time this value was Get()
    LastUsed *time.Time

    // size of the key, in bytes
    KeyBytes uint64

    // size of the value, in bytes
    ValueBytes uint64
}

type CacheMap struct {
    cache                    map[interface{}]*Item
    callbacks                Callbacks
    callbackMutex            sync.Mutex
    mutex                    sync.Mutex
    CallbacksAreThreadSafe   bool
    CallbackTimeout          time.Duration
    MaxItems                 int
    _bytes                   uint64
    MaxBytes                 uint64
    evictorRunning           bool
    lazyEvictor              bool
    EvictorDelay             time.Duration
}


const NoTimeout = 0 * time.Second
const DefaultEvictorDelay = 3 * time.Second
const DefaultMaxBytes = 128*1024*1024

func NewCacheMap(c Callbacks) CacheMap {
    return CacheMap {
        cache: make(map[interface{}]*Item),
        MaxItems: 10000,
        _bytes: 0,
        EvictorDelay: DefaultEvictorDelay,
        MaxBytes: DefaultMaxBytes,
        callbacks: c,
    }
}

func NewLazyCacheMap(c Callbacks) CacheMap {
    return CacheMap {
        cache: make(map[interface{}]*Item),
        MaxItems: 10000,
        _bytes: 0,
        EvictorDelay: DefaultEvictorDelay,
        MaxBytes: DefaultMaxBytes,
        lazyEvictor: true,
        callbacks: c,
    }
}

func (m CacheMap) CachedItems() int {
    return len(m.cache)
}

func (m CacheMap) Bytes() uint64 {
    return m._bytes
}

func (m CacheMap) wouldNeedEviction(newbytes uint64, items int) bool {
    return newbytes > m.MaxBytes || items > m.MaxItems
}

func (m *CacheMap) callCreate(key interface{}, item *Item) error {
    if ! m.CallbacksAreThreadSafe {
        defer m.callbackMutex.Unlock()
        m.callbackMutex.Lock()
    }
    if m.CallbackTimeout != NoTimeout {
        return m.callbacks.Create(key, item)
    }
    c := make(chan error)
    go func() {
        c <- m.callbacks.Create(key, item)
    }()
    select {
    case err := <-c:
        return err
    case <-time.After(m.CallbackTimeout):
        return errors.New("cachemap callCreate timeout")
    }
}

func (m *CacheMap) callRead(key interface{}) (item *Item, err error) {
    if ! m.CallbacksAreThreadSafe {
        defer m.callbackMutex.Unlock()
        m.callbackMutex.Lock()
    }
    if m.CallbackTimeout != NoTimeout {
        return m.callbacks.Read(key)
    }
    type loadResult struct {
        *Item
        error
    }
    c := make(chan loadResult)
    go func() {
        i, e := m.callbacks.Read(key)
        c <- loadResult{i, e}
    }()
    select {
    case result := <-c:
        err = result.error
        item = result.Item
    case <-time.After(m.CallbackTimeout):
        err = errors.New("cachemap callLoad timeout")
    }
    return
}

func (m *CacheMap) callUpdate(key interface{}, item *Item) error {
    if ! m.CallbacksAreThreadSafe {
        defer m.callbackMutex.Unlock()
        m.callbackMutex.Lock()
    }
    if m.CallbackTimeout != NoTimeout {
        return m.callbacks.Update(key, item)
    }
    c := make(chan error)
    go func() {
        c <- m.callbacks.Update(key, item)
    }()
    select {
    case err := <-c:
        return err
    case <-time.After(m.CallbackTimeout):
        return errors.New("cachemap callUpdate timeout")
    }
}

func (m *CacheMap) callDelete(key interface{}) error {
    if ! m.CallbacksAreThreadSafe {
        defer m.callbackMutex.Unlock()
        m.callbackMutex.Lock()
    }
    if m.CallbackTimeout != NoTimeout {
        return m.callbacks.Delete(key)
    }
    c := make(chan error)
    go func() {
        c <- m.callbacks.Delete(key)
    }()
    select {
    case err := <-c:
        return err
    case <-time.After(m.CallbackTimeout):
        return errors.New("cachemap callDelete timeout")
    }
}


// replace existing item in cache
func (m *CacheMap) updateItem(key interface{}, prev, _new *Item) (err error) {
    err = m.callUpdate(key, _new)
    if err != nil {
        return
    }
    // should an immediate evict happen?
    //
    tempbytes := m._bytes + (_new.KeyBytes + _new.ValueBytes) - (prev.KeyBytes + prev.ValueBytes)
    if !m.lazyEvictor && m.wouldNeedEviction(tempbytes, m.CachedItems()+1) {
        err = m.chooseEvictionVictim()
        if err != nil {
            return
        }
    }
    // update the cache item
    prev.Value = _new.Value
    prev.KeyBytes = _new.KeyBytes
    prev.ValueBytes = _new.ValueBytes
    prev.Created = _new.Created
    m.removeBytes(prev.KeyBytes + prev.ValueBytes)
    m.addBytes(_new.KeyBytes + _new.ValueBytes)
    return
}

var noTime = time.Time{}

// add a new item to cache
func (m *CacheMap) createItem(key interface{}, _new *Item) (err error) {
    err = m.callCreate(key, _new)
    if err != nil {
        return
    }
    // should an immediate evict happen?
    tempbytes := m._bytes + (_new.KeyBytes + _new.ValueBytes)
    if !m.lazyEvictor && m.wouldNeedEviction(tempbytes, m.CachedItems()+1) {
        err = m.chooseEvictionVictim()
        if err != nil {
            return
        }
    }
    // persist the item
    m.cache[key] = _new
    m.addBytes(_new.KeyBytes + _new.ValueBytes)
    return
}

// add key, value to the map
func (m *CacheMap) Add(key, value interface{}, KeyBytes, ValueBytes uint64) error {
    now := time.Now()
    item := &Item{
        Value: value,
        KeyBytes: KeyBytes,
        ValueBytes: ValueBytes,
        Created: &now,
    }
    if prev := m.cache[key]; prev != nil {
        // replace a previous item
        return m.updateItem(key, prev, item)
    }
    // add a new item
    return m.createItem(key, item)
}

// Empty the cache, dont delete anything
func (m *CacheMap) DropCaches() {
    m.cache = make(map[interface{}]*Item)
    m._bytes = 0
}

// get an item from cache.  if not cached, fetch it
func (m *CacheMap) Get(key interface{}) (value interface{}, err error) {
    defer m.mutex.Unlock()
    m.mutex.Lock()
    cached := m.cache[key]
    // is cached?
    if cached != nil {
        now := time.Now()
        cached.LastUsed = &now
        value = cached.Value
        return
    }
    // is persisted?
    item, err := m.callRead(key)
    if err == nil {
        m.cache[key] = item
        value = item.Value
    }
    return
}

func (m *CacheMap) addBytes(bytes uint64) {
    m._bytes += bytes
}

func (m *CacheMap) removeBytes(bytes uint64) {
    m._bytes -= bytes
}

func (m *CacheMap) deleteFromCache(key interface{}) {
    item := m.cache[key]
    if item != nil {
        delete(m.cache, key)
        m.removeBytes(item.KeyBytes + item.ValueBytes)
    }
}

// remove an item from the in-memory cache
// the item was already saved
func (m *CacheMap) Evict(key interface{}) (err error) {
    cached := m.cache[key]
    if cached == nil {
        err = errors.New("cannot persist item not in cache")
        return
    }
    m.deleteFromCache(key)
    return
}

// delete an item from the cache
func (m *CacheMap) Delete(key interface{}) error {
    m.deleteFromCache(key)
    return m.callDelete(key)
}

func (m *CacheMap) chooseEvictionVictim() error {
    // 1. find OLDEST UNUSED item
    var oldestKey *interface{}
    var oldestItem *Item
    for key, item := range m.cache {
        if item.LastUsed == nil {
            continue
        }
        if oldestItem == nil || item.Created.Before(*oldestItem.Created) {
            oldestKey = &key
            oldestItem = item
        }
    }
    // was there an unused item found?
    if oldestKey != nil && oldestItem != nil {
        // dont evict items that are < 10 seconds old
        // if now - 10s < oldestUnused.created 
        tenSecondsAgo := time.Now().Add(-10 * time.Second)
        if oldestItem.Created.Before(tenSecondsAgo) {
            return m.Evict(oldestKey)
        }
    }

    // 2. if no unused item(s), find the OLDEST item used
    for key, item := range m.cache {
        if oldestItem == nil || item.LastUsed.Before(*oldestItem.LastUsed) {
            oldestItem = item
            oldestKey = &key
        }
    }
    // ws there an oldest item?
    if oldestKey != nil {
        return m.Evict(oldestKey)
    }
    panic("no item to evict from cachemap")
    return nil
}

func (m CacheMap) IsLazyEvictor() bool {
    return m.lazyEvictor
}

// stop evictor
func (m *CacheMap) StopEvictor() {
    if ! m.lazyEvictor {
        panic("Not a lazy evictor")
    }
    if ! m.evictorRunning {
        panic("evictor not running")
    }
    m.evictorRunning = false
}

// start evictor
func (m *CacheMap) StartEvictor() {
    if ! m.lazyEvictor {
        panic("Not a lazy evictor")
    }
    if m.evictorRunning {
        panic("evictor already started")
    }
    m.evictorRunning = true
    go func() {
        log.Println("evictor started")
        for m.evictorRunning {
            for len(m.cache) > m.MaxItems && m._bytes > m.MaxBytes {
                if err := m.chooseEvictionVictim(); err != nil {
                    log.Println("eviction error", err)
                }
            }
            log.Println("running GC")
            runtime.GC()
            time.Sleep(m.EvictorDelay)
        }
        log.Println("evictor stopped")
    }()
}
