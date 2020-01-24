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
	"errors"
	"fmt"
	"github.com/AletheiaWareLLC/bcgo"
	"github.com/golang/protobuf/proto"
)

const ERROR_CREATOR_SENDER_DONT_MATCH = "Record Creator and Transaction Sender don't match: %s vs %s"

type TransactionValidator struct {
}

func (t *TransactionValidator) Validate(channel *bcgo.Channel, cache bcgo.Cache, network bcgo.Network, hash []byte, block *bcgo.Block) error {
	return bcgo.Iterate(channel.Name, hash, block, cache, network, func(h []byte, b *bcgo.Block) error {
		for _, entry := range b.Entry {
			// Unmarshal as Transaction
			t := &Transaction{}
			err := proto.Unmarshal(entry.Record.Payload, t)
			if err != nil {
				return err
			}
			// Check Record Creator matches Sender of Tokens
			if entry.Record.Creator != t.Sender {
				return errors.New(fmt.Sprintf(ERROR_CREATOR_SENDER_DONT_MATCH, entry.Record.Creator, t.Sender))
			}
		}
		return nil
	})
}
