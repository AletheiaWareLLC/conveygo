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
	"github.com/AletheiaWareLLC/aliasgo"
	"github.com/AletheiaWareLLC/bcgo"
	"github.com/AletheiaWareLLC/conveygo"
	"github.com/AletheiaWareLLC/testinggo"
	"io/ioutil"
	"os"
	"path"
	"testing"
)

func makeBCStore(t *testing.T, alias string, key *rsa.PrivateKey, keystore string) *conveygo.BCStore {
	t.Helper()
	dir, err := ioutil.ReadDir(keystore)
	testinggo.AssertNoError(t, err)
	for _, d := range dir {
		os.RemoveAll(path.Join(keystore, d.Name()))
	}
	store := &conveygo.BCStore{
		Node: &bcgo.Node{
			Alias:    alias,
			Key:      key,
			Cache:    bcgo.NewMemoryCache(1),
			Network:  bcgo.NewTCPNetwork(),
			Channels: make(map[string]*bcgo.Channel),
		},
		Listener: nil,
		KeyStore: keystore,
	}
	store.Node.AddChannel(aliasgo.OpenAliasChannel())
	store.Node.AddChannel(conveygo.OpenConversationChannel())
	return store
}

func TestBCStore(t *testing.T) {
	dir, err := ioutil.TempDir("", "keystore")
	testinggo.AssertNoError(t, err)
	defer os.RemoveAll(dir)
	aliasA := "Alice"
	keyA, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		t.Error("Could not generate key:", err)
	}
	aliasB := "Bob"
	emailB := "bob@example.com"
	paymentB := "payment1234"
	passwordB := []byte("password1234")
	keyB, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		t.Error("Could not generate key:", err)
	}
	t.Run("AddKey", func(t *testing.T) {
		t.Run("Exists", func(t *testing.T) {
			testUserStore_AddKey_Exists(t, makeBCStore(t, aliasA, keyA, dir), aliasB, passwordB, keyB)
		})
		t.Run("NotExists", func(t *testing.T) {
			testUserStore_AddKey_NotExists(t, makeBCStore(t, aliasA, keyA, dir), aliasB, passwordB, keyB)
		})
	})
	t.Run("GetKey", func(t *testing.T) {
		t.Run("Exists", func(t *testing.T) {
			testUserStore_GetKey_Exists(t, makeBCStore(t, aliasA, keyA, dir), aliasB, passwordB, keyB)
		})
		t.Run("NotExists", func(t *testing.T) {
			testUserStore_GetKey_NotExists(t, makeBCStore(t, aliasA, keyA, dir), aliasB, passwordB, keyB)
		})
	})
	t.Run("HasKey", func(t *testing.T) {
		t.Run("Exists", func(t *testing.T) {
			testUserStore_HasKey_Exists(t, makeBCStore(t, aliasA, keyA, dir), aliasB, passwordB, keyB)
		})
		t.Run("NotExists", func(t *testing.T) {
			testUserStore_HasKey_NotExists(t, makeBCStore(t, aliasA, keyA, dir), aliasB, passwordB, keyB)
		})
	})
	t.Run("RegisterAlias", func(t *testing.T) {
		t.Run("Exists", func(t *testing.T) {
			testUserStore_RegisterAlias_Exists(t, makeBCStore(t, aliasA, keyA, dir), aliasB, emailB, paymentB, keyB)
		})
		t.Run("NotExists", func(t *testing.T) {
			testUserStore_RegisterAlias_NotExists(t, makeBCStore(t, aliasA, keyA, dir), aliasB, emailB, paymentB, keyB)
		})
	})
	t.Run("RegisterCustomer", func(t *testing.T) {
		t.Run("Exists", func(t *testing.T) {
			testUserStore_RegisterCustomer_Exists(t, makeBCStore(t, aliasA, keyA, dir), aliasB, emailB, paymentB, keyB)
		})
		t.Run("NotExists", func(t *testing.T) {
			testUserStore_RegisterCustomer_NotExists(t, makeBCStore(t, aliasA, keyA, dir), aliasB, emailB, paymentB, keyB)
		})
	})
	t.Run("GetRegistration", func(t *testing.T) {
		t.Run("Exists", func(t *testing.T) {
			testUserStore_GetRegistration_Exists(t, makeBCStore(t, aliasA, keyA, dir), aliasB, emailB, paymentB, keyB)
		})
		t.Run("NotExists", func(t *testing.T) {
			testUserStore_GetRegistration_NotExists(t, makeBCStore(t, aliasA, keyA, dir), aliasB, emailB, paymentB, keyB)
		})
	})
	t.Run("SubscribeCustomer", func(t *testing.T) {
		t.Run("Exists", func(t *testing.T) {
			testUserStore_SubscribeCustomer_Exists(t, makeBCStore(t, aliasA, keyA, dir), aliasB, emailB, paymentB, keyB)
		})
		t.Run("NotExists", func(t *testing.T) {
			testUserStore_SubscribeCustomer_NotExists(t, makeBCStore(t, aliasA, keyA, dir), aliasB, emailB, paymentB, keyB)
		})
	})
	t.Run("GetSubscription", func(t *testing.T) {
		t.Run("Exists", func(t *testing.T) {
			testUserStore_GetSubscription_Exists(t, makeBCStore(t, aliasA, keyA, dir), aliasB, emailB, paymentB, keyB)
		})
		t.Run("NotExists", func(t *testing.T) {
			testUserStore_GetSubscription_NotExists(t, makeBCStore(t, aliasA, keyA, dir), aliasB, emailB, paymentB, keyB)
		})
	})
	t.Run("NewConversation", func(t *testing.T) {
		testConversationStore_NewConversation(t, makeBCStore(t, aliasA, keyA, dir), aliasB, keyB)
	})
	t.Run("GetConversation", func(t *testing.T) {
		t.Run("Exists", func(t *testing.T) {
			testConversationStore_GetConversation_Exists(t, makeBCStore(t, aliasA, keyA, dir), aliasB, keyB)
		})
		t.Run("NotExists", func(t *testing.T) {
			testConversationStore_GetConversation_NotExists(t, makeBCStore(t, aliasA, keyA, dir))
		})
	})
	t.Run("GetAllConversations", func(t *testing.T) {
		t.Run("Empty", func(t *testing.T) {
			testConversationStore_GetAllConversations_Empty(t, makeBCStore(t, aliasA, keyA, dir))
		})
		t.Run("NotEmpty", func(t *testing.T) {
			testConversationStore_GetAllConversations_NotEmpty(t, makeBCStore(t, aliasA, keyA, dir), aliasB, keyB)
		})
		t.Run("From", func(t *testing.T) {
			testConversationStore_GetAllConversations_From(t, makeBCStore(t, aliasA, keyA, dir), aliasB, keyB)
		})
		t.Run("To", func(t *testing.T) {
			testConversationStore_GetAllConversations_To(t, makeBCStore(t, aliasA, keyA, dir), aliasB, keyB)
		})
	})
	t.Run("GetRecentConversations", func(t *testing.T) {
		t.Run("Empty", func(t *testing.T) {
			testConversationStore_GetRecentConversations_Empty(t, makeBCStore(t, aliasA, keyA, dir))
		})
		t.Run("NotEmpty", func(t *testing.T) {
			testConversationStore_GetRecentConversations_NotEmpty(t, makeBCStore(t, aliasA, keyA, dir), aliasB, keyB)
		})
		t.Run("Limit", func(t *testing.T) {
			testConversationStore_GetRecentConversations_Limit(t, makeBCStore(t, aliasA, keyA, dir), aliasB, keyB)
		})
	})
	t.Run("AddMessage", func(t *testing.T) {
		t.Run("Exists", func(t *testing.T) {
			testMessageStore_AddMessage_Exists(t, makeBCStore(t, aliasA, keyA, dir), aliasB, keyB)
		})
		t.Run("NotExists", func(t *testing.T) {
			testMessageStore_AddMessage_NotExists(t, makeBCStore(t, aliasA, keyA, dir), aliasB, keyB)
		})
	})
	t.Run("MineBlockEntry", func(t *testing.T) {
		// TODO testinggo.AssertNoError(t, s.MineBlockEntry(channel, entry))
	})
	t.Run("GetMessage", func(t *testing.T) {
		t.Run("Exists", func(t *testing.T) {
			testMessageStore_GetMessage_Exists(t, makeBCStore(t, aliasA, keyA, dir), aliasB, keyB)
		})
		t.Run("Exists_Hash", func(t *testing.T) {
			testMessageStore_GetMessage_Exists_Hash(t, makeBCStore(t, aliasA, keyA, dir), aliasB, keyB)
		})
		t.Run("Exists_Reply", func(t *testing.T) {
			testMessageStore_GetMessage_Exists_Reply(t, makeBCStore(t, aliasA, keyA, dir), aliasB, keyB)
		})
		t.Run("NotExists", func(t *testing.T) {
			testMessageStore_GetMessage_NotExists(t, makeBCStore(t, aliasA, keyA, dir))
		})
	})
	t.Run("GetYield", func(t *testing.T) {
		t.Run("Exists", func(t *testing.T) {
			testMessageStore_GetYield_Exists(t, makeBCStore(t, aliasA, keyA, dir), aliasB, keyB)
		})
		t.Run("Exists_Reply", func(t *testing.T) {
			testMessageStore_GetYield_Exists_Reply(t, makeBCStore(t, aliasA, keyA, dir), aliasB, keyB)
		})
		t.Run("NotExists", func(t *testing.T) {
			testMessageStore_GetYield_NotExists(t, makeBCStore(t, aliasA, keyA, dir))
		})
	})
}
