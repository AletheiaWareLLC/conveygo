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
	"crypto/rsa"
	"github.com/AletheiaWareLLC/bcgo"
	"github.com/AletheiaWareLLC/financego"
)

const (
	ERROR_ACCESS_DENIED        = "Access denied"
	ERROR_NO_SUCH_CONVERSATION = "No such conversation: %s"
	ERROR_KEY_ALREADY_EXISTS   = "Key already exists: %s"
)

type ConversationStore interface {
	NewConversation(conversationHash []byte, conversationRecord *bcgo.Record, messageHash []byte, messageRecord *bcgo.Record) error
	GetConversation(conversationHash []byte) (*Listing, error)
	GetAllConversations(from, to uint64) ([]*Listing, error)
	GetRecentConversations(limit uint) ([]*Listing, error)
}

type MessageStore interface {
	ConversationStore
	AddMessage(conversationHash, messageHash []byte, messageRecord *bcgo.Record) error
	GetMessage(conversationHash, messageHash []byte, callback func([]byte, uint64, string, uint64, *Message) error) error
	GetYield(conversationHash []byte) (uint64, uint64, error)
}

type UserStore interface {
	AddKey(alias string, password []byte, key *rsa.PrivateKey) error
	GetKey(alias string, password []byte) (*rsa.PrivateKey, error)
	HasKey(alias string) bool
	RegisterAlias(alias string, password []byte, key *rsa.PrivateKey) error
	RegisterCustomer(alias string, key *rsa.PrivateKey, customerId string) error
	GetRegistration(alias string) (*financego.Registration, error)
	// TODO(v3) SubscribeCustomer(alias string, key *rsa.PrivateKey, customer, product, plan string) error
	// TODO(v3) GetSubscription(alias string) (*financego.Subscription, error)
}
