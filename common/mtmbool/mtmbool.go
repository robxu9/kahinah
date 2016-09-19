// Package mtmbool presents a threaded map[uint]bool instance,
// allowing for a system of atomic reads and writes where reads
// will always occur after all pending writes have finished
package mtmbool

import "sync"

// Map represents a thread-safe map[uint]bool
type Map struct {
	container     map[uint]bool
	containerLock *sync.Mutex
	writes        chan write

	reading *sync.Mutex
}

type write struct {
	key    uint
	value  bool
	delete bool
}

// NewMap creates a new thread-safe map.
func NewMap() *Map {
	m := &Map{
		container:     map[uint]bool{},
		containerLock: &sync.Mutex{},
		writes:        make(chan write, 100),
		reading:       &sync.Mutex{},
	}
	go m.flush()

	return m
}

func (m *Map) flush() {
	for {
		m.containerLock.Lock()
	InnerL:
		for {
			select {
			case w, ok := <-m.writes:
				if !ok {
					m.containerLock.Unlock()
					return
				}

				if w.delete {
					delete(m.container, w.key)
				} else {
					m.container[w.key] = w.value
				}
			default:
				break InnerL
			}
		}
		m.containerLock.Unlock()

		// grab the read lock at the end in case a read is waiting so that we don't flush
		m.reading.Lock()
		m.reading.Unlock()
	}
}

// Read retrieves the value corresponding to the key specified.
func (m *Map) Read(key uint) bool {
	m.reading.Lock()
	defer m.reading.Unlock()

	// wait for pending writes to flush
	m.containerLock.Lock() // we're waiting..... (when this returns, we've got control)
	defer m.containerLock.Unlock()

	return m.container[key]
}

// Write retrieves the value corresponding to the key specified.
func (m *Map) Write(key uint, value bool) {
	m.writes <- write{key: key, value: value}
}

// CompareAndSet writes the value you're trying to set ONLY if the read
// condition is set; if not, returns false.
func (m *Map) CompareAndSet(key uint, expected bool, value bool) bool {
	m.reading.Lock()
	defer m.reading.Unlock()

	// wait for pending writes to flush
	m.containerLock.Lock()
	defer m.containerLock.Unlock()

	if m.container[key] != expected {
		return false
	}

	m.container[key] = value
	return true
}

// Delete deletes the value corresponding to the key specified.
func (m *Map) Delete(key uint) {
	m.writes <- write{key: key, delete: true}
}

// Destroy stops flushing the map and allows it to be GC'ed.
func (m *Map) Destroy() {
	close(m.writes) // no more values can be written
}
