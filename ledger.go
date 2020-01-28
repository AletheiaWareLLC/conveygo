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
	"github.com/golang/protobuf/proto"
	"log"
	"math"
	"strings"
)

/*
Ledger Economics

Minted - an Alias mints Tokens by mining Periodic Validation Chains (PVC) which strengthen the Network;
    - Convey-Hour: The Hourly Periodic Validation Chain awards a miner 3600 Tokens per Block
    - Convey-Day: The Daily Periodic Validation Chain awards a miner 86400 Tokens per Block
    - Convey-Week: The Weekly Periodic Validation Chain awards a miner 604800 Tokens per Block
    - Convey-Year: The Yearly Periodic Validation Chain awards a miner 31557600 Tokens per Block
    - Convey-Decade: The Decennially Periodic Validation Chain awards a miner 315576000 Tokens per Block
    - Convey-Century: The Centennially Periodic Validation Chain awards a miner 3155760000 Tokens per Block

Burned - an Alias burns Tokens by starting Conversations;
    - Convey-Conversation: The Conversation Chain burns 1 Token per 100 Bytes for each Record in each Blocks
    - Convey-Message-<HASH>: The Message Chain burns 1 Token per 100 Bytes for the first Message in each Conversation (subsequent Messages are Replies in which tokens are Spent by the Reply Author and Earned by the Author of the Message being replied to)s

Bought - an Alias gains Tokens by buying them from another.

Sold - an Alias loses Tokens by selling them to another.

Earned - an Alias earns 1/2 Token per 100 Bytes for each Record in any Message Chain that replies to a Message they Authored.
	- It gets slightly more complicated with replies to a reply;
		1. Alice starts a conversation.
		2. Bob replies to Alice, half the tokens Bob spent go to Alice, and the rest are burned.
		3. Charlie replies to Bob, half the tokens Charlie spent go to Bob, half the remaining tokens go to Alice, and the rest are burned.
		4. Daniel replies to Charlie, half the tokens Daniel spent go to Charlie, half the remaining tokens go to Bob, half the remaining tokens go to Alice, and the rest are burned.
		5. Emma replies to Daniel, half the tokens Emma spent go to Daniel, half the remaining tokens go to Charlie, half the remaining tokens go to Bob, half the remaining tokens go to Alice, and the rest are burned.
	- Half is an integer division so if the number of tokens is odd (ie 3) the smaller half (ie 1) is awarded to the Message being replied to and the larger half (ie 2) is awarded up the Message hierarchy. This means 1 token will always get burned.

Spent - an Alias spends 1 Token per 100 Bytes for each Record they Author in any Conversation or Message Chain.

Balance - an Alias's balance is;
    - credited for each Token
        - minted by Mining PVC
        - bought through a Transaction
        - earned from Replies
    - debited for each Token
        - sold through a Transaction
        - burned or spent in Conversations and Messages
*/

const (
	// PVC Reward is 1 Token per Second
	// 6 PVCs means 6 Tokens per Second average
	HOURLY_PVC_REWARD       = 3600       // (60 * 60)
	DAILY_PVC_REWARD        = 86400      // (60 * 60 * 24)
	WEEKLY_PVC_REWARD       = 604800     // (60 * 60 * 24 * 7)
	YEARLY_PVC_REWARD       = 31557600   // (60 * 60 * 24 * 365.25)
	DECENNIALLY_PVC_REWARD  = 315576000  // (60 * 60 * 24 * 365.25 * 10)
	CENTENNIALLY_PVC_REWARD = 3155760000 // (60 * 60 * 24 * 365.25 * 100)
)

type Ledger struct {
	Node      *bcgo.Node
	Processed map[string]map[string]bool // Channel Name -> Block Hash -> Processed Flag
	Aliases   map[string]bool            // Alias -> Seen Flag
	Minted    map[string]uint64
	Burned    map[string]uint64
	Bought    map[string]uint64
	Sold      map[string]uint64
	Earned    map[string]uint64
	Spent     map[string]uint64
	Trigger   chan bool
}

func NewLedger(node *bcgo.Node) *Ledger {
	ledger := &Ledger{
		Node:      node,
		Processed: make(map[string]map[string]bool),
		Aliases:   make(map[string]bool),
		Minted:    make(map[string]uint64),
		Burned:    make(map[string]uint64),
		Bought:    make(map[string]uint64),
		Sold:      make(map[string]uint64),
		Earned:    make(map[string]uint64),
		Spent:     make(map[string]uint64),
	}
	return ledger
}

func Record(m map[string]uint64, key string, amount uint64) {
	a, ok := m[key]
	if !ok {
		a = 0
	}
	a += amount
	m[key] = a
}

func (l *Ledger) RecordMinted(alias string, amount uint64) {
	// log.Println(alias, "minted", amount)
	l.Aliases[alias] = true
	Record(l.Minted, alias, amount)
}

func (l *Ledger) RecordBurned(alias string, amount uint64) {
	// log.Println(alias, "burned", amount)
	l.Aliases[alias] = true
	Record(l.Burned, alias, amount)
}

func (l *Ledger) RecordBought(alias string, amount uint64) {
	// log.Println(alias, "bought", amount)
	l.Aliases[alias] = true
	Record(l.Bought, alias, amount)
}

func (l *Ledger) RecordSold(alias string, amount uint64) {
	// log.Println(alias, "sold", amount)
	l.Aliases[alias] = true
	Record(l.Sold, alias, amount)
}

func (l *Ledger) RecordEarned(alias string, amount uint64) {
	// log.Println(alias, "earned", amount)
	l.Aliases[alias] = true
	Record(l.Earned, alias, amount)
}

func (l *Ledger) RecordSpent(alias string, amount uint64) {
	// log.Println(alias, "spent", amount)
	l.Aliases[alias] = true
	Record(l.Spent, alias, amount)
}

func (l *Ledger) GetBalance(alias string) int64 {
	var balance int64
	balance += int64(l.Minted[alias])
	balance -= int64(l.Burned[alias])
	balance += int64(l.Bought[alias])
	balance -= int64(l.Sold[alias])
	balance += int64(l.Earned[alias])
	balance -= int64(l.Spent[alias])
	return balance
}

func Cost(record *bcgo.Record) uint64 {
	return uint64(math.Ceil(float64(proto.Size(record)) / 100))
}

// Iterates through unprocessed blocks in the given channel
func (l *Ledger) iterate(channel string, hash []byte, callback func([]byte, *bcgo.Block) error) error {
	processed, ok := l.Processed[channel]
	if !ok {
		processed = make(map[string]bool)
		l.Processed[channel] = processed
	}
	if err := bcgo.Iterate(channel, hash, nil, l.Node.Cache, l.Node.Network, func(h []byte, b *bcgo.Block) error {
		key := base64.RawURLEncoding.EncodeToString(h)
		if processed[key] {
			return bcgo.StopIterationError{}
		}
		processed[key] = true
		return callback(h, b)
	}); err != nil {
		switch err.(type) {
		case bcgo.StopIterationError:
			// Do nothing
			break
		default:
			return err
		}
	}
	return nil
}

type MessageNode struct {
	Author   string
	Cost     uint64
	Previous string
}

func (l *Ledger) Update(name string, hash []byte) error {
	if hash == nil {
		return nil
	}
	// log.Println("Ledger Update", name, base64.RawURLEncoding.EncodeToString(hash))
	switch name {
	case CONVEY_HOUR:
		// Block Miner mints 3600 Tokens per Block
		if err := l.iterate(name, hash, func(h []byte, b *bcgo.Block) error {
			l.RecordMinted(b.Miner, HOURLY_PVC_REWARD)
			return nil
		}); err != nil {
			return err
		}
	case CONVEY_DAY:
		// Block Miner mints 86400 Tokens per Block
		if err := l.iterate(name, hash, func(h []byte, b *bcgo.Block) error {
			l.RecordMinted(b.Miner, DAILY_PVC_REWARD)
			return nil
		}); err != nil {
			return err
		}
	case CONVEY_WEEK:
		// Block Miner mints 604800 Tokens per Block
		if err := l.iterate(name, hash, func(h []byte, b *bcgo.Block) error {
			l.RecordMinted(b.Miner, WEEKLY_PVC_REWARD)
			return nil
		}); err != nil {
			return err
		}
	case CONVEY_YEAR:
		// Block Miner mints 31557600 Tokens per Block
		if err := l.iterate(name, hash, func(h []byte, b *bcgo.Block) error {
			l.RecordMinted(b.Miner, YEARLY_PVC_REWARD)
			return nil
		}); err != nil {
			return err
		}
	case CONVEY_DECADE:
		// Block Miner mints 315576000 Tokens per Block
		if err := l.iterate(name, hash, func(h []byte, b *bcgo.Block) error {
			l.RecordMinted(b.Miner, DECENNIALLY_PVC_REWARD)
			return nil
		}); err != nil {
			return err
		}
	case CONVEY_CENTURY:
		// Block Miner mints 3155760000 Tokens per Block
		if err := l.iterate(name, hash, func(h []byte, b *bcgo.Block) error {
			l.RecordMinted(b.Miner, CENTENNIALLY_PVC_REWARD)
			return nil
		}); err != nil {
			return err
		}
	case CONVEY_TRANSACTION:
		// Holds transactions where Sender sells Tokens, Recipient buys Tokens
		if err := l.iterate(name, hash, func(h []byte, b *bcgo.Block) error {
			for _, entry := range b.Entry {
				record := entry.Record
				// Unmarshal as Transaction
				t := &Transaction{}
				err := proto.Unmarshal(record.Payload, t)
				if err != nil {
					return err
				}
				amount := t.Amount
				l.RecordSold(t.Sender, amount)
				l.RecordBought(t.Receiver, amount)
			}
			return nil
		}); err != nil {
			return err
		}
	case CONVEY_CONVERSATION:
		// Record Author burns 1 Token per 100 Bytes
		if err := l.iterate(name, hash, func(h []byte, b *bcgo.Block) error {
			for _, entry := range b.Entry {
				l.RecordBurned(entry.Record.Creator, Cost(entry.Record))
			}
			return nil
		}); err != nil {
			return err
		}
	default:
		if strings.HasPrefix(name, CONVEY_PREFIX_MESSAGE) {
			// If First Message: Author burns 1 Token per 100 Bytes
			// Else: Author spends 1 Token per 100 Bytes
			blocks := make(map[string]string)
			nodes := make(map[string]*MessageNode)
			if err := bcgo.Iterate(name, hash, nil, l.Node.Cache, l.Node.Network, func(h []byte, b *bcgo.Block) error {
				blockKey := base64.RawURLEncoding.EncodeToString(h)
				for _, entry := range b.Entry {
					recordKey := base64.RawURLEncoding.EncodeToString(entry.RecordHash)
					blocks[recordKey] = blockKey
					record := entry.Record
					node := &MessageNode{
						Author: record.Creator,
						Cost:   Cost(record),
					}
					// Unmarshal as Message
					m := &Message{}
					err := proto.Unmarshal(record.Payload, m)
					if err != nil {
						return err
					}
					if m.Previous != nil && len(m.Previous) > 0 {
						node.Previous = base64.RawURLEncoding.EncodeToString(m.Previous)
					}
					nodes[recordKey] = node
				}
				return nil
			}); err != nil {
				return err
			}

			processed, ok := l.Processed[name]
			if !ok {
				processed = make(map[string]bool)
				l.Processed[name] = processed
			}

			for key, node := range nodes {
				if processed[blocks[key]] {
					continue // Skip blocks that have already been processed
				}
				author := node.Author
				cost := node.Cost
				prev := node.Previous
				if prev == "" {
					l.RecordBurned(author, cost)
				} else {
					for {
						half := cost / 2
						previous := nodes[prev]
						if previous.Previous == "" {
							// half awarded to author of previous, remaining tokens are burned
							l.RecordSpent(author, half)
							l.RecordEarned(previous.Author, half)
							l.RecordBurned(author, cost-half)
							break
						} else {
							// half awarded to author of previous, remaining tokens go up the hierarchy
							l.RecordSpent(author, half)
							l.RecordEarned(previous.Author, half)
							cost -= half
							prev = previous.Previous
						}
					}
				}
			}

			for _, block := range blocks {
				processed[block] = true
			}
		}
	}
	return nil
}

func (l *Ledger) UpdateAll() error {
	for _, channel := range l.Node.GetChannels() {
		if err := l.Update(channel.Name, channel.Head); err != nil {
			return err
		}
	}
	return nil
}

func (l *Ledger) Start() {
	l.Trigger = make(chan bool)
	for ok := true; ok; ok = <-l.Trigger {
		if err := l.UpdateAll(); err != nil {
			log.Println(err)
			return
		}
	}
}

func (l *Ledger) Stop() {
	close(l.Trigger)
	l.Trigger = nil
}

func (l *Ledger) TriggerUpdate() {
	if l.Trigger != nil {
		l.Trigger <- true
	}
}
