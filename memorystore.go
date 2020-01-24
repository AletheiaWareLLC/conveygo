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
	"bytes"
	"crypto/rsa"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/AletheiaWareLLC/bcgo"
	"github.com/AletheiaWareLLC/financego"
	"github.com/golang/protobuf/proto"
	"sort"
)

type MemoryStore struct {
	Passwords     map[string][]byte
	Keys          map[string]*rsa.PrivateKey
	Timestamps    map[string]uint64
	Conversations map[string]*bcgo.Record
	Mappings      map[string][]string
	Messages      map[string]*bcgo.Record
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		Passwords:     make(map[string][]byte),
		Keys:          make(map[string]*rsa.PrivateKey),
		Timestamps:    make(map[string]uint64),
		Conversations: make(map[string]*bcgo.Record),
		Mappings:      make(map[string][]string),
		Messages:      make(map[string]*bcgo.Record),
	}
}

func (s *MemoryStore) AddKey(alias string, password []byte, key *rsa.PrivateKey) error {
	if _, ok := s.Passwords[alias]; ok {
		return errors.New(fmt.Sprintf(ERROR_KEY_ALREADY_EXISTS, alias))
	}
	s.Passwords[alias] = password
	s.Keys[alias] = key
	return nil
}

func (s *MemoryStore) GetKey(alias string, password []byte) (*rsa.PrivateKey, error) {
	pwd, ok := s.Passwords[alias]
	if ok && bytes.Equal(pwd, password) {
		return s.Keys[alias], nil
	}
	return nil, errors.New(ERROR_ACCESS_DENIED)
}

func (s *MemoryStore) HasKey(alias string) bool {
	_, ok := s.Passwords[alias]
	return ok
}

func (s *MemoryStore) RegisterAlias(alias string, password []byte, key *rsa.PrivateKey) error {
	// TODO
	return nil
}

func (s *MemoryStore) RegisterCustomer(alias string, key *rsa.PrivateKey, customerId string) error {
	// TODO
	return nil
}

func (s *MemoryStore) GetRegistration(alias string) (*financego.Registration, error) {
	// TODO
	return nil, nil
}

func (s *MemoryStore) SubscribeCustomer(alias string, key *rsa.PrivateKey, customer, payment, product, plan string) error {
	// TODO
	return nil
}

func (s *MemoryStore) GetSubscription(alias string) (*financego.Subscription, error) {
	// TODO
	return nil, nil
}

func (s *MemoryStore) NewConversation(conversationHash []byte, conversationRecord *bcgo.Record, messageHash []byte, messageRecord *bcgo.Record) error {
	conversationHashString := base64.RawURLEncoding.EncodeToString(conversationHash)
	s.Conversations[conversationHashString] = conversationRecord
	s.Timestamps[conversationHashString] = conversationRecord.Timestamp

	if err := s.AddMessage(conversationHash, messageHash, messageRecord); err != nil {
		return err
	}

	return nil
}

func (s *MemoryStore) GetConversation(conversationHash []byte) (*Listing, error) {
	hash := base64.RawURLEncoding.EncodeToString(conversationHash)
	record, ok := s.Conversations[hash]
	if !ok {
		return nil, errors.New(fmt.Sprintf(ERROR_NO_SUCH_CONVERSATION, hash))
	}
	c := &Conversation{}
	if err := proto.Unmarshal(record.Payload, c); err != nil {
		return nil, err
	}
	return &Listing{
		Hash:      conversationHash,
		Timestamp: record.Timestamp,
		Author:    record.Creator,
		Topic:     c.Topic,
		Cost:      Cost(record),
	}, nil
}

func (s *MemoryStore) GetAllConversations(since uint64) ([]*Listing, error) {
	var listings []*Listing
	for conversationHashString, value := range s.Conversations {
		if s.Timestamps[conversationHashString] >= since {
			conversationHash, err := base64.RawURLEncoding.DecodeString(conversationHashString)
			if err != nil {
				return nil, err
			}
			c := &Conversation{}
			if err := proto.Unmarshal(value.Payload, c); err != nil {
				return nil, err
			}
			listings = append(listings, &Listing{
				Hash:      conversationHash,
				Timestamp: s.Timestamps[conversationHashString],
				Author:    "",
				Topic:     c.Topic,
				Cost:      Cost(value),
			})
		}
	}
	return listings, nil
}

func (s *MemoryStore) GetRecentConversations(limit uint) ([]*Listing, error) {
	var listings []*Listing
	for conversationHashString, value := range s.Conversations {
		conversationHash, err := base64.RawURLEncoding.DecodeString(conversationHashString)
		if err != nil {
			return nil, err
		}
		c := &Conversation{}
		if err := proto.Unmarshal(value.Payload, c); err != nil {
			return nil, err
		}
		listings = append(listings, &Listing{
			Hash:      conversationHash,
			Timestamp: s.Timestamps[conversationHashString],
			Author:    "",
			Topic:     c.Topic,
			Cost:      Cost(value),
		})
	}
	sort.Slice(listings, func(i, j int) bool {
		return listings[i].Timestamp > listings[j].Timestamp
	})
	if uint(len(listings)) > limit {
		listings = listings[:limit]
	}
	return listings, nil
}

func (s *MemoryStore) AddMessage(conversationHash, messageHash []byte, messageRecord *bcgo.Record) error {
	conversationKey := base64.RawURLEncoding.EncodeToString(conversationHash)
	_, ok := s.Conversations[conversationKey]
	if !ok {
		return errors.New(fmt.Sprintf(ERROR_NO_SUCH_CONVERSATION, conversationKey))
	}
	messageKey := base64.RawURLEncoding.EncodeToString(messageHash)
	s.Mappings[conversationKey] = append(s.Mappings[conversationKey], messageKey)
	s.Messages[messageKey] = messageRecord
	return nil
}

func (s *MemoryStore) GetMessage(conversationHash, messageHash []byte, callback func([]byte, uint64, string, uint64, *Message) error) error {
	conversationHashString := base64.RawURLEncoding.EncodeToString(conversationHash)
	mappings, ok := s.Mappings[conversationHashString]
	if !ok {
		return errors.New(fmt.Sprintf(ERROR_NO_SUCH_CONVERSATION, conversationHashString))
	}
	for _, m := range mappings {
		hash, err := base64.RawURLEncoding.DecodeString(m)
		if err != nil {
			return err
		}
		if messageHash == nil || bytes.Equal(messageHash, hash) {
			record := s.Messages[m]
			message := &Message{}
			if err := proto.Unmarshal(record.Payload, message); err != nil {
				return err
			}
			if err := callback(hash, 0, "", Cost(record), message); err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *MemoryStore) GetYield(conversationHash []byte) (uint64, uint64, error) {
	// TODO
	return 0, 0, nil
}
