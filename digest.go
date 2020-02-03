/*
 * Copyright 2020 Aletheia Ware LLC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package conveygo

import (
	"encoding/base64"
	"github.com/AletheiaWareLLC/bcgo"
	"sort"
)

const (
	DIGEST_LIMIT = 4
)

type DigestEntry struct {
	Hash      string
	Topic     string
	Timestamp string
	Author    string
	Cost      uint64
	Reward    uint64
	Yield     int64
	Message   *Message
}

// Returns the 4 highest-yielding conversations from the given time period.
func GetDigestEntries(messages MessageStore, from, to uint64) ([]*DigestEntry, error) {
	conversations, err := messages.GetAllConversations(from, to)
	if err != nil {
		return nil, err
	}

	var entries []*DigestEntry
	for _, c := range conversations {
		entry := &DigestEntry{
			Hash:      base64.RawURLEncoding.EncodeToString(c.Hash),
			Topic:     c.Topic,
			Timestamp: bcgo.TimestampToString(c.Timestamp),
			Cost:      c.Cost,
			Author:    c.Author,
		}

		// Get first message
		if err := messages.GetMessage(c.Hash, nil, func(hash []byte, timestamp uint64, author string, cost uint64, message *Message) error {
			if message.Previous == nil || len(message.Previous) == 0 {
				entry.Message = message
			}
			return nil
		}); err != nil {
			return nil, err
		}

		// Get conversation yield
		cost, reward, err := messages.GetYield(c.Hash)
		if err != nil {
			return nil, err
		}

		entry.Cost += cost
		entry.Reward += reward
		entry.Yield = int64(entry.Reward) - int64(entry.Cost)
		entries = append(entries, entry)
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Yield > entries[j].Yield
	})

	if uint(len(entries)) > DIGEST_LIMIT {
		entries = entries[:DIGEST_LIMIT]
	}

	return entries, nil
}
