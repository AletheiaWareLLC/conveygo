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
	"github.com/AletheiaWareLLC/aliasgo"
	"github.com/AletheiaWareLLC/bcgo"
	"github.com/AletheiaWareLLC/cryptogo"
	"github.com/AletheiaWareLLC/financego"
	"github.com/golang/protobuf/proto"
	"log"
)

type BCStore struct {
	Node     *bcgo.Node
	Listener bcgo.MiningListener
	KeyStore string
}

func (s *BCStore) AddKey(alias string, password []byte, key *rsa.PrivateKey) error {
	if s.HasKey(alias) {
		return errors.New(fmt.Sprintf(ERROR_KEY_ALREADY_EXISTS, alias))
	}
	// Write private key to keystore on filesystem
	return cryptogo.WriteRSAPrivateKey(key, s.KeyStore, alias, password)
}

func (s *BCStore) GetKey(alias string, password []byte) (*rsa.PrivateKey, error) {
	key, err := cryptogo.GetRSAPrivateKey(s.KeyStore, alias, password)
	if err != nil {
		log.Println(err)
		return nil, errors.New(ERROR_ACCESS_DENIED)
	}
	return key, nil
}

func (s *BCStore) HasKey(alias string) bool {
	return cryptogo.HasRSAPrivateKey(s.KeyStore, alias)
}

func (s *BCStore) RegisterAlias(alias string, password []byte, key *rsa.PrivateKey) error {
	// Create alias record
	record, err := aliasgo.CreateSignedAliasRecord(alias, key)
	if err != nil {
		return err
	}

	// Write alias record to cache
	if _, err := bcgo.WriteRecord(aliasgo.ALIAS, s.Node.Cache, record); err != nil {
		return err
	}

	// Get Alias Channel
	aliases, err := s.Node.GetChannel(aliasgo.ALIAS)
	if err != nil {
		return err
	}

	// Mine alias record
	if _, _, err := s.Node.Mine(aliases, aliasgo.ALIAS_THRESHOLD, s.Listener); err != nil {
		return err
	}

	// Push Update to Network
	if err := aliases.Push(s.Node.Cache, s.Node.Network); err != nil {
		return err
	}
	return nil
}

func (s *BCStore) RegisterCustomer(alias string, key *rsa.PrivateKey, customerId string) error {
	// Create Registration proto
	registration := &financego.Registration{
		MerchantAlias: s.Node.Alias,
		CustomerAlias: alias,
		Processor:     financego.PaymentProcessor_STRIPE,
		CustomerId:    customerId,
	}
	log.Println("Registration", registration)

	// Create Access Control List
	acl := map[string]*rsa.PublicKey{
		alias:        &key.PublicKey,
		s.Node.Alias: &s.Node.Key.PublicKey,
	}
	log.Println("Access", acl)

	// Marshal Registration proto
	registrationData, err := proto.Marshal(registration)
	if err != nil {
		return err
	}

	// Get Registration Channel
	registrations, err := s.Node.GetChannel(CONVEY_REGISTRATION)
	if err != nil {
		return err
	}

	if err := registrations.LoadCachedHead(s.Node.Cache); err != nil {
		log.Println(err)
	}

	if err := registrations.Pull(s.Node.Cache, s.Node.Network); err != nil {
		log.Println(err)
	}

	// Write Registration data to cache
	if _, err := s.Node.Write(bcgo.Timestamp(), registrations, acl, nil, registrationData); err != nil {
		return err
	}

	// Mine Registration Chain
	if _, _, err := s.Node.Mine(registrations, bcgo.THRESHOLD_G, s.Listener); err != nil {
		return err
	}

	// Push Update to Network
	if err := registrations.Push(s.Node.Cache, s.Node.Network); err != nil {
		log.Println(err)
	}

	return nil
}

func (s *BCStore) GetRegistration(alias string) (*financego.Registration, error) {
	// Get Registration Channel
	registrations, err := s.Node.GetChannel(CONVEY_REGISTRATION)
	if err != nil {
		return nil, err
	}
	var registration *financego.Registration
	// Read Registration Channel
	if err := bcgo.Read(registrations.Name, registrations.Head, nil, s.Node.Cache, s.Node.Network, s.Node.Alias, s.Node.Key, nil, func(entry *bcgo.BlockEntry, key, data []byte) error {
		// Unmarshal as Registration
		r := &financego.Registration{}
		err := proto.Unmarshal(data, r)
		if err != nil {
			return err
		}
		if r.CustomerAlias == alias {
			registration = r
			return bcgo.StopIterationError{}
		}
		return nil
	}); err != nil {
		switch err.(type) {
		case bcgo.StopIterationError:
			// Do nothing
			break
		default:
			return nil, err
		}
	}
	return registration, nil
}

func (s *BCStore) NewConversation(conversationHash []byte, conversationRecord *bcgo.Record, messageHash []byte, messageRecord *bcgo.Record) error {
	conversations, err := s.Node.GetChannel(CONVEY_CONVERSATION)
	if err != nil {
		return err
	}

	if err := s.MineBlockEntry(conversations, &bcgo.BlockEntry{
		RecordHash: conversationHash,
		Record:     conversationRecord,
	}); err != nil {
		return err
	}

	channel := CONVEY_PREFIX_MESSAGE + base64.RawURLEncoding.EncodeToString(conversationHash)
	messages := bcgo.OpenPoWChannel(channel, bcgo.THRESHOLD_G)
	s.Node.AddChannel(messages)

	if err := s.MineBlockEntry(messages, &bcgo.BlockEntry{
		RecordHash: messageHash,
		Record:     messageRecord,
	}); err != nil {
		return err
	}

	return nil
}

func (s *BCStore) GetConversation(conversationHash []byte) (*Listing, error) {
	conversations, err := s.Node.GetChannel(CONVEY_CONVERSATION)
	if err != nil {
		return nil, err
	}
	var listing *Listing
	if err := bcgo.Iterate(conversations.Name, conversations.Head, nil, s.Node.Cache, s.Node.Network, func(h []byte, b *bcgo.Block) error {
		for _, entry := range b.Entry {
			if bytes.Equal(conversationHash, entry.RecordHash) {
				listing, err = ConversationEntryToListing(entry)
				if err != nil {
					return err
				}
				return bcgo.StopIterationError{}
			}
		}
		return nil
	}); err != nil {
		switch err.(type) {
		case bcgo.StopIterationError:
			// Do nothing
			break
		default:
			return nil, err
		}
	}
	if listing == nil {
		return nil, errors.New(fmt.Sprintf(ERROR_NO_SUCH_CONVERSATION, base64.RawURLEncoding.EncodeToString(conversationHash)))
	}
	return listing, nil
}

func (s *BCStore) GetAllConversations(since uint64) ([]*Listing, error) {
	conversations, err := s.Node.GetChannel(CONVEY_CONVERSATION)
	if err != nil {
		return nil, err
	}
	var listings []*Listing
	if err := bcgo.Iterate(conversations.Name, conversations.Head, nil, s.Node.Cache, s.Node.Network, func(h []byte, b *bcgo.Block) error {
		if b.Timestamp < since {
			return bcgo.StopIterationError{}
		}
		for _, entry := range b.Entry {
			if entry.Record.Timestamp >= since {
				listing, err := ConversationEntryToListing(entry)
				if err != nil {
					return err
				}
				listings = append(listings, listing)
			}
		}
		return nil
	}); err != nil {
		switch err.(type) {
		case bcgo.StopIterationError:
			// Do nothing
			break
		default:
			return nil, err
		}
	}
	return listings, nil
}

func (s *BCStore) GetRecentConversations(limit uint) ([]*Listing, error) {
	conversations, err := s.Node.GetChannel(CONVEY_CONVERSATION)
	if err != nil {
		return nil, err
	}
	var listings []*Listing
	if err := bcgo.Iterate(conversations.Name, conversations.Head, nil, s.Node.Cache, s.Node.Network, func(h []byte, b *bcgo.Block) error {
		for _, entry := range b.Entry {
			listing, err := ConversationEntryToListing(entry)
			if err != nil {
				return err
			}
			listings = append(listings, listing)
			limit--
			if limit <= 0 {
				return bcgo.StopIterationError{}
			}
		}
		return nil
	}); err != nil {
		switch err.(type) {
		case bcgo.StopIterationError:
			// Do nothing
			break
		default:
			return nil, err
		}
	}
	return listings, nil
}

func (s *BCStore) AddMessage(conversationHash, messageHash []byte, messageRecord *bcgo.Record) error {
	conversationHashString := base64.RawURLEncoding.EncodeToString(conversationHash)
	channel := CONVEY_PREFIX_MESSAGE + conversationHashString
	messages, err := s.Node.GetChannel(channel)
	if err != nil {
		return errors.New(fmt.Sprintf(ERROR_NO_SUCH_CONVERSATION, conversationHashString))
	}

	if err := s.MineBlockEntry(messages, &bcgo.BlockEntry{
		RecordHash: messageHash,
		Record:     messageRecord,
	}); err != nil {
		return err
	}

	return nil
}

func (s *BCStore) MineBlockEntry(channel *bcgo.Channel, entry *bcgo.BlockEntry) error {
	// Write Record to Cache
	if err := s.Node.Cache.PutBlockEntry(channel.Name, entry); err != nil {
		return err
	}

	// Mine Channel
	if _, _, err := s.Node.Mine(channel, bcgo.THRESHOLD_G, s.Listener); err != nil {
		return err
	}

	if s.Node.Network != nil {
		// Push Update to Peers
		if err := channel.Push(s.Node.Cache, s.Node.Network); err != nil {
			return err
		}
	}

	return nil
}

func (s *BCStore) GetMessage(conversationHash, messageHash []byte, callback func([]byte, uint64, string, uint64, *Message) error) error {
	conversationHashString := base64.RawURLEncoding.EncodeToString(conversationHash)
	channel := CONVEY_PREFIX_MESSAGE + conversationHashString
	messages, err := s.Node.GetChannel(channel)
	if err != nil {
		return errors.New(fmt.Sprintf(ERROR_NO_SUCH_CONVERSATION, conversationHashString))
	}

	return bcgo.Iterate(channel, messages.Head, nil, s.Node.Cache, s.Node.Network, func(h []byte, b *bcgo.Block) error {
		for _, entry := range b.Entry {
			if messageHash == nil || bytes.Equal(messageHash, entry.RecordHash) {
				// Unmarshal as Message
				m := &Message{}
				err := proto.Unmarshal(entry.Record.Payload, m)
				if err != nil {
					return err
				}
				if err := callback(entry.RecordHash, entry.Record.GetTimestamp(), entry.Record.GetCreator(), Cost(entry.Record), m); err != nil {
					return err
				}
			}
		}
		return nil
	})
}

type ReplyNode struct {
	Previous string
	Cost     uint64
	Reward   uint64
}

func (s *BCStore) GetYield(conversationHash []byte) (uint64, uint64, error) {
	var messageHash string
	var messageCost uint64
	var messageReward uint64
	replies := make(map[string]*ReplyNode)
	if err := s.GetMessage(conversationHash, nil, func(hash []byte, timestamp uint64, author string, cost uint64, message *Message) error {
		key := base64.RawURLEncoding.EncodeToString(hash)
		if message.Previous == nil || len(message.Previous) == 0 {
			messageHash = key
			messageCost = cost
		} else {
			replies[key] = &ReplyNode{
				Previous: base64.RawURLEncoding.EncodeToString(message.Previous),
				Cost:     cost,
			}
		}
		return nil
	}); err != nil {
		return 0, 0, err
	}

	for _, reply := range replies {
		// Calculate rewards
		reward := reply.Cost
		for reply != nil && reward > 0 {
			half := reward / 2 // Integer division so half of 3 is 1
			if reply.Previous == messageHash {
				messageReward += half
				break
			} else {
				reply = replies[reply.Previous]
				reward -= half
			}
		}
	}
	return messageCost, messageReward, nil
}
