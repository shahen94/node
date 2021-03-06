/*
 * Copyright (C) 2018 The "MysteriumNetwork/node" Authors.
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package session

import (
	"sync"
)

// NewStorageMemory initiates new session storage
func NewStorageMemory() *StorageMemory {
	return &StorageMemory{
		sessions:  make([]Session, 0),
		idToIndex: make(map[ID]int),
		lock:      sync.Mutex{},
	}
}

// StorageMemory maintains all currents sessions in memory
type StorageMemory struct {
	sessions  []Session
	idToIndex map[ID]int
	lock      sync.Mutex
}

// Add puts given session to storage. Multiple sessions per peerID is possible in case different services are used
func (storage *StorageMemory) Add(sessionInstance Session) {
	storage.lock.Lock()
	defer storage.lock.Unlock()

	sessionIndex := len(storage.sessions)
	storage.sessions = append(storage.sessions, sessionInstance)

	storage.idToIndex[sessionInstance.ID] = sessionIndex
}

// GetAll returns all sessions in storage
func (storage *StorageMemory) GetAll() ([]Session, error) {
	storage.lock.Lock()
	defer storage.lock.Unlock()

	return storage.sessions, nil
}

// Find returns underlying session instance
func (storage *StorageMemory) Find(id ID) (Session, bool) {
	storage.lock.Lock()
	defer storage.lock.Unlock()

	if sessionIndex, found := storage.idToIndex[id]; found {
		instance := storage.sessions[sessionIndex]
		if len(storage.sessions) == 1 {
			instance.Last = true
		}
		return instance, true
	}

	return Session{}, false
}

// Remove removes given session from underlying storage
func (storage *StorageMemory) Remove(id ID) {
	storage.lock.Lock()
	defer storage.lock.Unlock()

	if sessionIndex, found := storage.idToIndex[id]; found {
		delete(storage.idToIndex, id)
		storage.sessions = append(storage.sessions[:sessionIndex], storage.sessions[sessionIndex+1:]...)
	}
}
