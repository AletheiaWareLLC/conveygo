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
	"crypto/rsa"
	"fmt"
	"github.com/AletheiaWareLLC/bcgo"
	"github.com/AletheiaWareLLC/conveygo"
	"github.com/AletheiaWareLLC/testinggo"
	"testing"
)

func testUserStore_AddKey_Exists(t *testing.T, s conveygo.UserStore, alias string, password []byte, key *rsa.PrivateKey) {
	t.Helper()
	testinggo.AssertNoError(t, s.AddKey(alias, password, key))
	testinggo.AssertError(t, fmt.Sprintf(conveygo.ERROR_KEY_ALREADY_EXISTS, alias), s.AddKey(alias, password, key))
}

func testUserStore_AddKey_NotExists(t *testing.T, s conveygo.UserStore, alias string, password []byte, key *rsa.PrivateKey) {
	t.Helper()
	testinggo.AssertNoError(t, s.AddKey(alias, password, key))
}

func testUserStore_GetKey_Exists(t *testing.T, s conveygo.UserStore, alias string, password []byte, key *rsa.PrivateKey) {
	t.Helper()
	testinggo.AssertNoError(t, s.AddKey(alias, password, key))
	_, err := s.GetKey(alias, append(password, password...))
	testinggo.AssertError(t, conveygo.ERROR_ACCESS_DENIED, err)
	actual, err := s.GetKey(alias, password)
	testinggo.AssertNoError(t, err)
	testinggo.AssertPrivateKeyEqual(t, key, actual)
}

func testUserStore_GetKey_NotExists(t *testing.T, s conveygo.UserStore, alias string, password []byte, key *rsa.PrivateKey) {
	t.Helper()
	_, err := s.GetKey(alias, password)
	testinggo.AssertError(t, conveygo.ERROR_ACCESS_DENIED, err)
}

func testUserStore_HasKey_Exists(t *testing.T, s conveygo.UserStore, alias string, password []byte, key *rsa.PrivateKey) {
	t.Helper()
	testinggo.AssertNoError(t, s.AddKey(alias, password, key))
	if !s.HasKey(alias) {
		t.Error("HasKey should return true")
	}
}

func testUserStore_HasKey_NotExists(t *testing.T, s conveygo.UserStore, alias string, password []byte, key *rsa.PrivateKey) {
	t.Helper()
	if s.HasKey(alias) {
		t.Error("HasKey should return false")
	}
}

func testUserStore_RegisterAlias_Exists(t *testing.T, s conveygo.UserStore, alias, email, payment string, key *rsa.PrivateKey) {
	t.Helper()
	// TODO
}

func testUserStore_RegisterAlias_NotExists(t *testing.T, s conveygo.UserStore, alias, email, payment string, key *rsa.PrivateKey) {
	t.Helper()
	// TODO
}

func testUserStore_RegisterCustomer_Exists(t *testing.T, s conveygo.UserStore, alias, email, payment string, key *rsa.PrivateKey) {
	t.Helper()
	// TODO
}

func testUserStore_RegisterCustomer_NotExists(t *testing.T, s conveygo.UserStore, alias, email, payment string, key *rsa.PrivateKey) {
	t.Helper()
	// TODO
}

func testUserStore_GetRegistration_Exists(t *testing.T, s conveygo.UserStore, alias, email, payment string, key *rsa.PrivateKey) {
	t.Helper()
	// TODO
}

func testUserStore_GetRegistration_NotExists(t *testing.T, s conveygo.UserStore, alias, email, payment string, key *rsa.PrivateKey) {
	t.Helper()
	// TODO
}

func testUserStore_SubscribeCustomer_Exists(t *testing.T, s conveygo.UserStore, alias, email, payment string, key *rsa.PrivateKey) {
	t.Helper()
	// TODO
}

func testUserStore_SubscribeCustomer_NotExists(t *testing.T, s conveygo.UserStore, alias, email, payment string, key *rsa.PrivateKey) {
	t.Helper()
	// TODO
}

func testUserStore_GetSubscription_Exists(t *testing.T, s conveygo.UserStore, alias, email, payment string, key *rsa.PrivateKey) {
	t.Helper()
	// TODO
}

func testUserStore_GetSubscription_NotExists(t *testing.T, s conveygo.UserStore, alias, email, payment string, key *rsa.PrivateKey) {
	t.Helper()
	// TODO
}

func testConversationStore_NewConversation(t *testing.T, s conveygo.ConversationStore, alias string, key *rsa.PrivateKey) {
	t.Helper()
	timestamp := bcgo.Timestamp()
	conversationHash, conversationRecord, err := conveygo.ProtoToRecord(alias, key, timestamp, &conveygo.Conversation{
		Topic: "Test123",
	})
	testinggo.AssertNoError(t, err)
	messageHash, messageRecord, err := conveygo.ProtoToRecord(alias, key, timestamp, &conveygo.Message{
		Content: []byte("FooBar"),
		Type:    conveygo.MediaType_TEXT_PLAIN,
	})
	testinggo.AssertNoError(t, err)
	testinggo.AssertNoError(t, s.NewConversation(conversationHash, conversationRecord, messageHash, messageRecord))
}

func testConversationStore_GetConversation_Exists(t *testing.T, s conveygo.ConversationStore, alias string, key *rsa.PrivateKey) {
	t.Helper()
	expected := "Test123"
	timestamp := bcgo.Timestamp()
	conversationHash, conversationRecord, err := conveygo.ProtoToRecord(alias, key, timestamp, &conveygo.Conversation{
		Topic: expected,
	})
	testinggo.AssertNoError(t, err)
	messageHash, messageRecord, err := conveygo.ProtoToRecord(alias, key, timestamp, &conveygo.Message{
		Content: []byte("FooBar"),
		Type:    conveygo.MediaType_TEXT_PLAIN,
	})
	testinggo.AssertNoError(t, err)
	testinggo.AssertNoError(t, s.NewConversation(conversationHash, conversationRecord, messageHash, messageRecord))
	listing, err := s.GetConversation(conversationHash)
	testinggo.AssertNoError(t, err)
	if listing.Topic != expected {
		t.Errorf("Wrong topic; expected '%s', got '%s'", expected, listing.Topic)
	}
}

func testConversationStore_GetConversation_NotExists(t *testing.T, s conveygo.ConversationStore) {
	t.Helper()
	listing, err := s.GetConversation([]byte("DoesNotExist"))
	testinggo.AssertError(t, "No such conversation: RG9lc05vdEV4aXN0", err)
	if listing != nil {
		t.Error("Expected Listing to be nil")
	}
}

func testConversationStore_GetAllConversations_Empty(t *testing.T, s conveygo.ConversationStore) {
	t.Helper()
	listings, err := s.GetAllConversations(uint64(0))
	testinggo.AssertNoError(t, err)
	if len(listings) != 0 {
		t.Errorf("Wrong listings; expected '%d', got '%d'", 0, len(listings))
	}
}

func testConversationStore_GetAllConversations_NotEmpty(t *testing.T, s conveygo.ConversationStore, alias string, key *rsa.PrivateKey) {
	t.Helper()
	expected := "Foo"
	timestamp := bcgo.Timestamp()
	conversationHash, conversationRecord, err := conveygo.ProtoToRecord(alias, key, timestamp, &conveygo.Conversation{
		Topic: expected,
	})
	testinggo.AssertNoError(t, err)
	messageHash, messageRecord, err := conveygo.ProtoToRecord(alias, key, timestamp, &conveygo.Message{
		Content: []byte("Test123"),
		Type:    conveygo.MediaType_TEXT_PLAIN,
	})
	testinggo.AssertNoError(t, err)
	testinggo.AssertNoError(t, s.NewConversation(conversationHash, conversationRecord, messageHash, messageRecord))
	listings, err := s.GetAllConversations(uint64(0))
	testinggo.AssertNoError(t, err)
	if len(listings) != 1 {
		t.Errorf("Wrong listings; expected '%d', got '%d'", 1, len(listings))
	}
	if listings[0].Topic != expected {
		t.Errorf("Wrong topic; expected '%s', got '%s'", expected, listings[0].Topic)
	}
}

func testConversationStore_GetAllConversations_Since(t *testing.T, s conveygo.ConversationStore, alias string, key *rsa.PrivateKey) {
	t.Helper()
	expected := "Foo"
	{
		timestamp := uint64(0)
		conversationHash, conversationRecord, err := conveygo.ProtoToRecord(alias, key, timestamp, &conveygo.Conversation{
			Topic: expected,
		})
		testinggo.AssertNoError(t, err)
		messageHash, messageRecord, err := conveygo.ProtoToRecord(alias, key, timestamp, &conveygo.Message{
			Content: []byte("Test123"),
			Type:    conveygo.MediaType_TEXT_PLAIN,
		})
		testinggo.AssertNoError(t, err)
		testinggo.AssertNoError(t, s.NewConversation(conversationHash, conversationRecord, messageHash, messageRecord))
	}
	expected = "Bar"
	{
		timestamp := bcgo.Timestamp()
		conversationHash, conversationRecord, err := conveygo.ProtoToRecord(alias, key, timestamp, &conveygo.Conversation{
			Topic: expected,
		})
		testinggo.AssertNoError(t, err)
		messageHash, messageRecord, err := conveygo.ProtoToRecord(alias, key, timestamp, &conveygo.Message{
			Content: []byte("Test123"),
			Type:    conveygo.MediaType_TEXT_PLAIN,
		})
		testinggo.AssertNoError(t, err)
		testinggo.AssertNoError(t, s.NewConversation(conversationHash, conversationRecord, messageHash, messageRecord))
	}
	listings, err := s.GetAllConversations(uint64(1))
	testinggo.AssertNoError(t, err)
	if len(listings) != 1 {
		t.Errorf("Wrong listings; expected '%d', got '%d'", 1, len(listings))
	}
	if listings[0].Topic != expected {
		t.Errorf("Wrong topic; expected '%s', got '%s'", expected, listings[0].Topic)
	}
}

func testConversationStore_GetRecentConversations_Empty(t *testing.T, s conveygo.ConversationStore) {
	t.Helper()
	listings, err := s.GetRecentConversations(1)
	testinggo.AssertNoError(t, err)
	if len(listings) != 0 {
		t.Errorf("Wrong listings; expected '%d', got '%d'", 0, len(listings))
	}
}

func testConversationStore_GetRecentConversations_NotEmpty(t *testing.T, s conveygo.ConversationStore, alias string, key *rsa.PrivateKey) {
	t.Helper()
	expected := "Foo"
	timestamp := bcgo.Timestamp()
	conversationHash, conversationRecord, err := conveygo.ProtoToRecord(alias, key, timestamp, &conveygo.Conversation{
		Topic: expected,
	})
	testinggo.AssertNoError(t, err)
	messageHash, messageRecord, err := conveygo.ProtoToRecord(alias, key, timestamp, &conveygo.Message{
		Content: []byte("Test123"),
		Type:    conveygo.MediaType_TEXT_PLAIN,
	})
	testinggo.AssertNoError(t, err)
	testinggo.AssertNoError(t, s.NewConversation(conversationHash, conversationRecord, messageHash, messageRecord))
	listings, err := s.GetRecentConversations(1)
	testinggo.AssertNoError(t, err)
	if len(listings) != 1 {
		t.Errorf("Wrong listings; expected '%d', got '%d'", 1, len(listings))
	}
	if listings[0].Topic != expected {
		t.Errorf("Wrong topic; expected '%s', got '%s'", expected, listings[0].Topic)
	}
}

func testConversationStore_GetRecentConversations_Limit(t *testing.T, s conveygo.ConversationStore, alias string, key *rsa.PrivateKey) {
	t.Helper()
	expected := "Foo"
	{
		timestamp := bcgo.Timestamp()
		conversationHash, conversationRecord, err := conveygo.ProtoToRecord(alias, key, timestamp, &conveygo.Conversation{
			Topic: expected,
		})
		testinggo.AssertNoError(t, err)
		messageHash, messageRecord, err := conveygo.ProtoToRecord(alias, key, timestamp, &conveygo.Message{
			Content: []byte("Test123"),
			Type:    conveygo.MediaType_TEXT_PLAIN,
		})
		testinggo.AssertNoError(t, err)
		testinggo.AssertNoError(t, s.NewConversation(conversationHash, conversationRecord, messageHash, messageRecord))
	}
	expected = "Bar"
	{
		timestamp := bcgo.Timestamp()
		conversationHash, conversationRecord, err := conveygo.ProtoToRecord(alias, key, timestamp, &conveygo.Conversation{
			Topic: expected,
		})
		testinggo.AssertNoError(t, err)
		messageHash, messageRecord, err := conveygo.ProtoToRecord(alias, key, timestamp, &conveygo.Message{
			Content: []byte("Test123"),
			Type:    conveygo.MediaType_TEXT_PLAIN,
		})
		testinggo.AssertNoError(t, err)
		testinggo.AssertNoError(t, s.NewConversation(conversationHash, conversationRecord, messageHash, messageRecord))
	}
	listings, err := s.GetRecentConversations(1)
	testinggo.AssertNoError(t, err)
	if len(listings) != 1 {
		t.Errorf("Wrong listings; expected '%d', got '%d'", 1, len(listings))
	}
	if listings[0].Topic != expected {
		t.Errorf("Wrong topic; expected '%s', got '%s'", expected, listings[0].Topic)
	}
}

func testMessageStore_AddMessage_Exists(t *testing.T, s conveygo.MessageStore, alias string, key *rsa.PrivateKey) {
	t.Helper()
	timestamp := bcgo.Timestamp()
	conversationHash, conversationRecord, err := conveygo.ProtoToRecord(alias, key, timestamp, &conveygo.Conversation{
		Topic: "Test123",
	})
	testinggo.AssertNoError(t, err)
	messageHash, messageRecord, err := conveygo.ProtoToRecord(alias, key, timestamp, &conveygo.Message{
		Content: []byte("Foo"),
		Type:    conveygo.MediaType_TEXT_PLAIN,
	})
	testinggo.AssertNoError(t, err)
	testinggo.AssertNoError(t, s.NewConversation(conversationHash, conversationRecord, messageHash, messageRecord))
	replyHash, replyRecord, err := conveygo.ProtoToRecord(alias, key, bcgo.Timestamp(), &conveygo.Message{
		Previous: messageHash,
		Content:  []byte("Bar"),
		Type:     conveygo.MediaType_TEXT_PLAIN,
	})
	testinggo.AssertNoError(t, err)
	testinggo.AssertNoError(t, s.AddMessage(conversationHash, replyHash, replyRecord))
}

func testMessageStore_AddMessage_NotExists(t *testing.T, s conveygo.MessageStore, alias string, key *rsa.PrivateKey) {
	t.Helper()
	replyHash, replyRecord, err := conveygo.ProtoToRecord(alias, key, bcgo.Timestamp(), &conveygo.Message{
		Content: []byte("FooBar"),
		Type:    conveygo.MediaType_TEXT_PLAIN,
	})
	testinggo.AssertNoError(t, err)
	err = s.AddMessage([]byte("ConversationDoesNotExist"), replyHash, replyRecord)
	testinggo.AssertError(t, "No such conversation: Q29udmVyc2F0aW9uRG9lc05vdEV4aXN0", err)
}

func testMessageStore_GetMessage_Exists(t *testing.T, s conveygo.MessageStore, alias string, key *rsa.PrivateKey) {
	t.Helper()
	timestamp := bcgo.Timestamp()
	conversationHash, conversationRecord, err := conveygo.ProtoToRecord(alias, key, timestamp, &conveygo.Conversation{
		Topic: "Test123",
	})
	testinggo.AssertNoError(t, err)
	messageHash, messageRecord, err := conveygo.ProtoToRecord(alias, key, timestamp, &conveygo.Message{
		Content: []byte("Foo"),
		Type:    conveygo.MediaType_TEXT_PLAIN,
	})
	testinggo.AssertNoError(t, err)
	testinggo.AssertNoError(t, s.NewConversation(conversationHash, conversationRecord, messageHash, messageRecord))
	results := make(map[string]*conveygo.Message)
	testinggo.AssertNoError(t, s.GetMessage(conversationHash, nil, func(hash []byte, timestamp uint64, author string, cost uint64, message *conveygo.Message) error {
		results[string(hash)] = message
		return nil
	}))
	if len(results) != 1 {
		t.Errorf("Incorrect number of results; expected '%d', got '%d'", 1, len(results))
	}
	message, ok := results[string(messageHash)]
	if !ok {
		t.Error("Results does not contain initial message")
	}
	content := string(message.Content)
	if content != "Foo" {
		t.Errorf("Incorrect message content; expected '%s', got '%s'", "Foo", content)
	}
}

func testMessageStore_GetMessage_Exists_Hash(t *testing.T, s conveygo.MessageStore, alias string, key *rsa.PrivateKey) {
	t.Helper()
	timestamp := bcgo.Timestamp()
	conversationHash, conversationRecord, err := conveygo.ProtoToRecord(alias, key, timestamp, &conveygo.Conversation{
		Topic: "Test123",
	})
	testinggo.AssertNoError(t, err)
	messageHash, messageRecord, err := conveygo.ProtoToRecord(alias, key, timestamp, &conveygo.Message{
		Content: []byte("Foo"),
		Type:    conveygo.MediaType_TEXT_PLAIN,
	})
	testinggo.AssertNoError(t, err)
	testinggo.AssertNoError(t, s.NewConversation(conversationHash, conversationRecord, messageHash, messageRecord))
	replyHash, replyRecord, err := conveygo.ProtoToRecord(alias, key, bcgo.Timestamp(), &conveygo.Message{
		Previous: messageHash,
		Content:  []byte("Bar"),
		Type:     conveygo.MediaType_TEXT_PLAIN,
	})
	testinggo.AssertNoError(t, err)
	testinggo.AssertNoError(t, s.AddMessage(conversationHash, replyHash, replyRecord))
	results := make(map[string]*conveygo.Message)
	testinggo.AssertNoError(t, s.GetMessage(conversationHash, replyHash, func(hash []byte, timestamp uint64, author string, cost uint64, message *conveygo.Message) error {
		results[string(hash)] = message
		return nil
	}))
	if len(results) != 1 {
		t.Errorf("Incorrect number of results; expected '%d', got '%d'", 1, len(results))
	}
	reply, ok := results[string(replyHash)]
	if !ok {
		t.Error("Results does not contain reply")
	}
	content := string(reply.Content)
	if content != "Bar" {
		t.Errorf("Incorrect reply content; expected '%s', got '%s'", "Bar", content)
	}
}

func testMessageStore_GetMessage_Exists_Reply(t *testing.T, s conveygo.MessageStore, alias string, key *rsa.PrivateKey) {
	t.Helper()
	timestamp := bcgo.Timestamp()
	conversationHash, conversationRecord, err := conveygo.ProtoToRecord(alias, key, timestamp, &conveygo.Conversation{
		Topic: "Test123",
	})
	testinggo.AssertNoError(t, err)
	messageHash, messageRecord, err := conveygo.ProtoToRecord(alias, key, timestamp, &conveygo.Message{
		Content: []byte("Foo"),
		Type:    conveygo.MediaType_TEXT_PLAIN,
	})
	testinggo.AssertNoError(t, err)
	testinggo.AssertNoError(t, s.NewConversation(conversationHash, conversationRecord, messageHash, messageRecord))
	replyHash, replyRecord, err := conveygo.ProtoToRecord(alias, key, bcgo.Timestamp(), &conveygo.Message{
		Previous: messageHash,
		Content:  []byte("Bar"),
		Type:     conveygo.MediaType_TEXT_PLAIN,
	})
	testinggo.AssertNoError(t, err)
	testinggo.AssertNoError(t, s.AddMessage(conversationHash, replyHash, replyRecord))
	results := make(map[string]*conveygo.Message)
	testinggo.AssertNoError(t, s.GetMessage(conversationHash, nil, func(hash []byte, timestamp uint64, author string, cost uint64, message *conveygo.Message) error {
		results[string(hash)] = message
		return nil
	}))
	if len(results) != 2 {
		t.Errorf("Incorrect number of results; expected '%d', got '%d'", 2, len(results))
	}
	message, ok := results[string(messageHash)]
	if !ok {
		t.Error("Results does not contain initial message")
	}
	content := string(message.Content)
	if content != "Foo" {
		t.Errorf("Incorrect message content; expected '%s', got '%s'", "Foo", content)
	}
	reply, ok := results[string(replyHash)]
	if !ok {
		t.Error("Results does not contain reply")
	}
	content = string(reply.Content)
	if content != "Bar" {
		t.Errorf("Incorrect reply content; expected '%s', got '%s'", "Bar", content)
	}
}

func testMessageStore_GetMessage_NotExists(t *testing.T, s conveygo.MessageStore) {
	t.Helper()
	results := make(map[string]*conveygo.Message)
	err := s.GetMessage([]byte("ConversationDoesNotExist"), nil, func(hash []byte, timestamp uint64, author string, cost uint64, message *conveygo.Message) error {
		results[string(hash)] = message
		return nil
	})
	testinggo.AssertError(t, "No such conversation: Q29udmVyc2F0aW9uRG9lc05vdEV4aXN0", err)
	if len(results) != 0 {
		t.Errorf("Incorrect number of results; expected '%d', got '%d'", 0, len(results))
	}
}

func testMessageStore_GetYield_Exists(t *testing.T, s conveygo.MessageStore, alias string, key *rsa.PrivateKey) {
	t.Helper()
	// TODO
}

func testMessageStore_GetYield_Exists_Reply(t *testing.T, s conveygo.MessageStore, alias string, key *rsa.PrivateKey) {
	t.Helper()
	// TODO
}

func testMessageStore_GetYield_NotExists(t *testing.T, s conveygo.MessageStore) {
	t.Helper()
	// TODO
}
