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

package conveygo_test

import (
	"crypto/rand"
	"crypto/rsa"
	"github.com/AletheiaWareLLC/bcgo"
	"github.com/AletheiaWareLLC/conveygo"
	"github.com/AletheiaWareLLC/testinggo"
	"github.com/golang/protobuf/proto"
	"io/ioutil"
	"log"
	"os"
	"testing"
	"time"
)

func makeNode(t *testing.T, alias string, key *rsa.PrivateKey) *bcgo.Node {
	t.Helper()
	node := &bcgo.Node{
		Alias:    alias,
		Key:      key,
		Cache:    bcgo.NewMemoryCache(1),
		Network:  &bcgo.TcpNetwork{},
		Channels: make(map[string]*bcgo.Channel),
	}
	node.AddChannel(bcgo.OpenPoWChannel(conveygo.CONVEY_YEAR, bcgo.THRESHOLD_Z))
	node.AddChannel(bcgo.OpenPoWChannel(conveygo.CONVEY_CONVERSATION, bcgo.THRESHOLD_Z))
	node.AddChannel(bcgo.OpenPoWChannel(conveygo.CONVEY_TRANSACTION, bcgo.THRESHOLD_Z))
	return node
}

func checkLedger(t *testing.T, ledger *conveygo.Ledger) {
	t.Helper()
	log.Println("Aliases:", ledger.Aliases)
	log.Println("Minted:", ledger.Minted)
	log.Println("Burned:", ledger.Burned)
	log.Println("Bought:", ledger.Bought)
	log.Println("Sold:", ledger.Sold)
	log.Println("Earned:", ledger.Earned)
	log.Println("Spent:", ledger.Spent)
	// Every Credit must have a matching Debit
	{
		var total int64
		// Credits
		for _, v := range ledger.Bought {
			total += int64(v)
		}
		// Debits
		for _, v := range ledger.Sold {
			total -= int64(v)
		}
		if total != 0 {
			t.Errorf("Tokens bought should equal tokens sold, instead got %d", total)
		}
	}
	{
		var total int64
		// Credits
		for _, v := range ledger.Earned {
			total += int64(v)
		}
		for _, v := range ledger.Spent {
			total -= int64(v)
		}
		if total != 0 {
			t.Errorf("Tokens earned should equal tokens spent, instead got %d", total)
		}
	}
	for alias := range ledger.Aliases {
		balance := ledger.GetBalance(alias)
		if balance < 0 {
			t.Errorf("Token balance for %s cannot be negative, instead got %d", alias, balance)
		}
	}
}

func makePeriodicValidationBlock(t *testing.T, node *bcgo.Node, listener bcgo.MiningListener, channel *bcgo.Channel, hash []byte, block *bcgo.Block) ([]byte, *bcgo.Block) {
	unix := uint64(time.Now().UnixNano())
	entries, err := bcgo.CreateValidationEntries(unix, node)
	testinggo.AssertNoError(t, err)
	b := bcgo.CreateValidationBlock(unix, channel.Name, node.Alias, hash, block, entries)
	h, _, err := node.MineBlock(channel, bcgo.THRESHOLD_Z, listener, b)
	testinggo.AssertNoError(t, err)
	return h, b
}

func makeTransactionBlock(t *testing.T, node *bcgo.Node, listener bcgo.MiningListener, channel *bcgo.Channel, alias string, key *rsa.PrivateKey, amount uint64) {
	transaction := &conveygo.Transaction{
		Sender:   node.Alias,
		Receiver: alias,
		Amount:   amount,
	}
	data, err := proto.Marshal(transaction)
	testinggo.AssertNoError(t, err)
	_, record, err := bcgo.CreateRecord(bcgo.Timestamp(), node.Alias, node.Key, nil, nil, data)
	testinggo.AssertNoError(t, err)
	_, err = bcgo.WriteRecord(channel.Name, node.Cache, record)
	testinggo.AssertNoError(t, err)
	_, _, err = node.Mine(channel, bcgo.THRESHOLD_Z, listener)
	testinggo.AssertNoError(t, err)
}

func TestLedger(t *testing.T) {
	listener := &bcgo.PrintingMiningListener{
		Output: os.Stdout,
	}
	aliasNode := "Node"
	keyNode, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		t.Error("Could not generate key:", err)
	}
	aliasAlice := "Alice"
	keyAlice, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		t.Error("Could not generate key:", err)
	}
	aliasBob := "Bob"
	keyBob, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		t.Error("Could not generate key:", err)
	}
	aliasCharlie := "Charlie"
	keyCharlie, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		t.Error("Could not generate key:", err)
	}
	t.Run("Conversation", func(t *testing.T) {
		node := makeNode(t, aliasNode, keyNode)
		years, err := node.GetChannel(conveygo.CONVEY_YEAR)
		testinggo.AssertNoError(t, err)
		makePeriodicValidationBlock(t, node, listener, years, nil, nil)
		transactions, err := node.GetChannel(conveygo.CONVEY_TRANSACTION)
		testinggo.AssertNoError(t, err)
		makeTransactionBlock(t, node, listener, transactions, aliasAlice, keyAlice, conveygo.YEARLY_PVC_REWARD)

		keystore, err := ioutil.TempDir("", "keystore")
		testinggo.AssertNoError(t, err)
		defer os.RemoveAll(keystore)
		store := &conveygo.BCStore{
			Node:     node,
			Listener: nil,
			KeyStore: keystore,
		}

		timestamp := bcgo.Timestamp()
		conversationHash, conversationRecord, err := conveygo.ProtoToRecord(aliasAlice, keyAlice, timestamp, &conveygo.Conversation{
			Topic: "Test123",
		})
		testinggo.AssertNoError(t, err)
		messageHash, messageRecord, err := conveygo.ProtoToRecord(aliasAlice, keyAlice, timestamp, &conveygo.Message{
			Content: []byte("Foo"),
			Type:    conveygo.MediaType_TEXT_PLAIN,
		})
		testinggo.AssertNoError(t, err)
		testinggo.AssertNoError(t, store.NewConversation(conversationHash, conversationRecord, messageHash, messageRecord))

		ledger := conveygo.NewLedger(node)
		ledger.TriggerUpdate()
		log.Println("Conversation Record:", conveygo.Cost(conversationRecord))
		log.Println("Message Record:", conveygo.Cost(messageRecord))

		checkLedger(t, ledger)

		if !ledger.Aliases[aliasNode] {
			t.Errorf("Missing alias: %s", aliasNode)
		}
		if !ledger.Aliases[aliasAlice] {
			t.Errorf("Missing alias: %s", aliasAlice)
		}
		if ledger.Aliases[aliasBob] {
			t.Errorf("Unexpected alias: %s", aliasBob)
		}
		if ledger.Aliases[aliasCharlie] {
			t.Errorf("Unexpected alias: %s", aliasCharlie)
		}
	})
	t.Run("Conversation_Reply", func(t *testing.T) {
		node := makeNode(t, aliasNode, keyNode)
		years, err := node.GetChannel(conveygo.CONVEY_YEAR)
		testinggo.AssertNoError(t, err)
		pvh, pvb := makePeriodicValidationBlock(t, node, listener, years, nil, nil)
		makePeriodicValidationBlock(t, node, listener, years, pvh, pvb)
		transactions, err := node.GetChannel(conveygo.CONVEY_TRANSACTION)
		testinggo.AssertNoError(t, err)
		makeTransactionBlock(t, node, listener, transactions, aliasAlice, keyAlice, conveygo.YEARLY_PVC_REWARD)
		makeTransactionBlock(t, node, listener, transactions, aliasBob, keyBob, conveygo.YEARLY_PVC_REWARD)

		keystore, err := ioutil.TempDir("", "keystore")
		testinggo.AssertNoError(t, err)
		defer os.RemoveAll(keystore)
		store := &conveygo.BCStore{
			Node:     node,
			Listener: nil,
			KeyStore: keystore,
		}

		timestamp := bcgo.Timestamp()
		conversationHash, conversationRecord, err := conveygo.ProtoToRecord(aliasAlice, keyAlice, timestamp, &conveygo.Conversation{
			Topic: "Test123",
		})
		testinggo.AssertNoError(t, err)
		messageHash, messageRecord, err := conveygo.ProtoToRecord(aliasAlice, keyAlice, timestamp, &conveygo.Message{
			Content: []byte("Foo"),
			Type:    conveygo.MediaType_TEXT_PLAIN,
		})
		testinggo.AssertNoError(t, err)
		testinggo.AssertNoError(t, store.NewConversation(conversationHash, conversationRecord, messageHash, messageRecord))

		replyHash, replyRecord, err := conveygo.ProtoToRecord(aliasBob, keyBob, bcgo.Timestamp(), &conveygo.Message{
			Previous: messageHash,
			Content:  []byte("Bar"),
			Type:     conveygo.MediaType_TEXT_PLAIN,
		})
		testinggo.AssertNoError(t, err)
		testinggo.AssertNoError(t, store.AddMessage(conversationHash, replyHash, replyRecord))

		ledger := conveygo.NewLedger(node)
		ledger.TriggerUpdate()
		log.Println("Conversation Record:", conveygo.Cost(conversationRecord))
		log.Println("Message Record:", conveygo.Cost(messageRecord))
		log.Println("Reply Record:", conveygo.Cost(replyRecord))

		checkLedger(t, ledger)

		if !ledger.Aliases[aliasNode] {
			t.Errorf("Missing alias: %s", aliasNode)
		}
		if !ledger.Aliases[aliasAlice] {
			t.Errorf("Missing alias: %s", aliasAlice)
		}
		if !ledger.Aliases[aliasBob] {
			t.Errorf("Missing alias: %s", aliasBob)
		}
		if ledger.Aliases[aliasCharlie] {
			t.Errorf("Unexpected alias: %s", aliasCharlie)
		}
	})
	t.Run("Conversation_Replies", func(t *testing.T) {
		node := makeNode(t, aliasNode, keyNode)
		years, err := node.GetChannel(conveygo.CONVEY_YEAR)
		testinggo.AssertNoError(t, err)
		pvh1, pvb1 := makePeriodicValidationBlock(t, node, listener, years, nil, nil)
		pvh2, pvb2 := makePeriodicValidationBlock(t, node, listener, years, pvh1, pvb1)
		makePeriodicValidationBlock(t, node, listener, years, pvh2, pvb2)
		transactions, err := node.GetChannel(conveygo.CONVEY_TRANSACTION)
		testinggo.AssertNoError(t, err)
		makeTransactionBlock(t, node, listener, transactions, aliasAlice, keyAlice, conveygo.YEARLY_PVC_REWARD)
		makeTransactionBlock(t, node, listener, transactions, aliasBob, keyBob, conveygo.YEARLY_PVC_REWARD)
		makeTransactionBlock(t, node, listener, transactions, aliasCharlie, keyCharlie, conveygo.YEARLY_PVC_REWARD)

		keystore, err := ioutil.TempDir("", "keystore")
		testinggo.AssertNoError(t, err)
		defer os.RemoveAll(keystore)
		store := &conveygo.BCStore{
			Node:     node,
			Listener: nil,
			KeyStore: keystore,
		}

		timestamp := bcgo.Timestamp()
		conversationHash, conversationRecord, err := conveygo.ProtoToRecord(aliasAlice, keyAlice, timestamp, &conveygo.Conversation{
			Topic: "Test123",
		})
		testinggo.AssertNoError(t, err)
		messageHash, messageRecord, err := conveygo.ProtoToRecord(aliasAlice, keyAlice, timestamp, &conveygo.Message{
			Content: []byte("Foo"),
			Type:    conveygo.MediaType_TEXT_PLAIN,
		})
		testinggo.AssertNoError(t, err)
		testinggo.AssertNoError(t, store.NewConversation(conversationHash, conversationRecord, messageHash, messageRecord))

		reply1Hash, reply1Record, err := conveygo.ProtoToRecord(aliasBob, keyBob, bcgo.Timestamp(), &conveygo.Message{
			Previous: messageHash,
			Content:  []byte("Bar"),
			Type:     conveygo.MediaType_TEXT_PLAIN,
		})
		testinggo.AssertNoError(t, err)
		testinggo.AssertNoError(t, store.AddMessage(conversationHash, reply1Hash, reply1Record))

		reply2Hash, reply2Record, err := conveygo.ProtoToRecord(aliasCharlie, keyCharlie, bcgo.Timestamp(), &conveygo.Message{
			Previous: reply1Hash,
			Content:  []byte("Baz"),
			Type:     conveygo.MediaType_TEXT_PLAIN,
		})
		testinggo.AssertNoError(t, err)
		testinggo.AssertNoError(t, store.AddMessage(conversationHash, reply2Hash, reply2Record))

		reply3Hash, reply3Record, err := conveygo.ProtoToRecord(aliasAlice, keyAlice, bcgo.Timestamp(), &conveygo.Message{
			Previous: reply2Hash,
			Content:  []byte("FooBarBaz"),
			Type:     conveygo.MediaType_TEXT_PLAIN,
		})
		testinggo.AssertNoError(t, err)
		testinggo.AssertNoError(t, store.AddMessage(conversationHash, reply3Hash, reply3Record))

		ledger := conveygo.NewLedger(node)
		ledger.TriggerUpdate()
		log.Println("Conversation Record:", conveygo.Cost(conversationRecord))
		log.Println("Message Record:", conveygo.Cost(messageRecord))
		log.Println("Reply 1 Record:", conveygo.Cost(reply1Record))
		log.Println("Reply 2 Record:", conveygo.Cost(reply2Record))
		log.Println("Reply 3 Record:", conveygo.Cost(reply3Record))

		checkLedger(t, ledger)

		if !ledger.Aliases[aliasNode] {
			t.Errorf("Missing alias: %s", aliasNode)
		}
		if !ledger.Aliases[aliasAlice] {
			t.Errorf("Missing alias: %s", aliasAlice)
		}
		if !ledger.Aliases[aliasBob] {
			t.Errorf("Missing alias: %s", aliasBob)
		}
		if !ledger.Aliases[aliasCharlie] {
			t.Errorf("Missing alias: %s", aliasCharlie)
		}
	})
}
