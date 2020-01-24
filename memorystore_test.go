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
	"github.com/AletheiaWareLLC/conveygo"
	"testing"
)

func TestMemoryStore(t *testing.T) {
	alias := "Alice"
	email := "alice@example.com"
	payment := "payment1234"
	password := []byte("password1234")
	key, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		t.Error("Could not generate key:", err)
	}
	t.Run("AddKey", func(t *testing.T) {
		t.Run("Exists", func(t *testing.T) {
			testUserStore_AddKey_Exists(t, conveygo.NewMemoryStore(), alias, password, key)
		})
		t.Run("NotExists", func(t *testing.T) {
			testUserStore_AddKey_NotExists(t, conveygo.NewMemoryStore(), alias, password, key)
		})
	})
	t.Run("GetKey", func(t *testing.T) {
		t.Run("Exists", func(t *testing.T) {
			testUserStore_GetKey_Exists(t, conveygo.NewMemoryStore(), alias, password, key)
		})
		t.Run("NotExists", func(t *testing.T) {
			testUserStore_GetKey_NotExists(t, conveygo.NewMemoryStore(), alias, password, key)
		})
	})
	t.Run("HasKey", func(t *testing.T) {
		t.Run("Exists", func(t *testing.T) {
			testUserStore_HasKey_Exists(t, conveygo.NewMemoryStore(), alias, password, key)
		})
		t.Run("NotExists", func(t *testing.T) {
			testUserStore_HasKey_NotExists(t, conveygo.NewMemoryStore(), alias, password, key)
		})
	})
	t.Run("RegisterAlias", func(t *testing.T) {
		t.Run("Exists", func(t *testing.T) {
			testUserStore_RegisterAlias_Exists(t, conveygo.NewMemoryStore(), alias, email, payment, key)
		})
		t.Run("NotExists", func(t *testing.T) {
			testUserStore_RegisterAlias_NotExists(t, conveygo.NewMemoryStore(), alias, email, payment, key)
		})
	})
	t.Run("RegisterCustomer", func(t *testing.T) {
		t.Run("Exists", func(t *testing.T) {
			testUserStore_RegisterCustomer_Exists(t, conveygo.NewMemoryStore(), alias, email, payment, key)
		})
		t.Run("NotExists", func(t *testing.T) {
			testUserStore_RegisterCustomer_NotExists(t, conveygo.NewMemoryStore(), alias, email, payment, key)
		})
	})
	t.Run("GetRegistration", func(t *testing.T) {
		t.Run("Exists", func(t *testing.T) {
			testUserStore_GetRegistration_Exists(t, conveygo.NewMemoryStore(), alias, email, payment, key)
		})
		t.Run("NotExists", func(t *testing.T) {
			testUserStore_GetRegistration_NotExists(t, conveygo.NewMemoryStore(), alias, email, payment, key)
		})
	})
	t.Run("SubscribeCustomer", func(t *testing.T) {
		t.Run("Exists", func(t *testing.T) {
			testUserStore_SubscribeCustomer_Exists(t, conveygo.NewMemoryStore(), alias, email, payment, key)
		})
		t.Run("NotExists", func(t *testing.T) {
			testUserStore_SubscribeCustomer_NotExists(t, conveygo.NewMemoryStore(), alias, email, payment, key)
		})
	})
	t.Run("GetSubscription", func(t *testing.T) {
		t.Run("Exists", func(t *testing.T) {
			testUserStore_GetSubscription_Exists(t, conveygo.NewMemoryStore(), alias, email, payment, key)
		})
		t.Run("NotExists", func(t *testing.T) {
			testUserStore_GetSubscription_NotExists(t, conveygo.NewMemoryStore(), alias, email, payment, key)
		})
	})
	t.Run("NewConversation", func(t *testing.T) {
		testConversationStore_NewConversation(t, conveygo.NewMemoryStore(), alias, key)
	})
	t.Run("GetConversation", func(t *testing.T) {
		t.Run("Exists", func(t *testing.T) {
			testConversationStore_GetConversation_Exists(t, conveygo.NewMemoryStore(), alias, key)
		})
		t.Run("NotExists", func(t *testing.T) {
			testConversationStore_GetConversation_NotExists(t, conveygo.NewMemoryStore())
		})
	})
	t.Run("GetAllConversations", func(t *testing.T) {
		t.Run("Empty", func(t *testing.T) {
			testConversationStore_GetAllConversations_Empty(t, conveygo.NewMemoryStore())
		})
		t.Run("NotEmpty", func(t *testing.T) {
			testConversationStore_GetAllConversations_NotEmpty(t, conveygo.NewMemoryStore(), alias, key)
		})
		t.Run("Since", func(t *testing.T) {
			testConversationStore_GetAllConversations_Since(t, conveygo.NewMemoryStore(), alias, key)
		})
	})
	t.Run("GetRecentConversations", func(t *testing.T) {
		t.Run("Empty", func(t *testing.T) {
			testConversationStore_GetRecentConversations_Empty(t, conveygo.NewMemoryStore())
		})
		t.Run("NotEmpty", func(t *testing.T) {
			testConversationStore_GetRecentConversations_NotEmpty(t, conveygo.NewMemoryStore(), alias, key)
		})
		t.Run("Limit", func(t *testing.T) {
			testConversationStore_GetRecentConversations_Limit(t, conveygo.NewMemoryStore(), alias, key)
		})
	})
	t.Run("AddMessage", func(t *testing.T) {
		t.Run("Exists", func(t *testing.T) {
			testMessageStore_AddMessage_Exists(t, conveygo.NewMemoryStore(), alias, key)
		})
		t.Run("NotExists", func(t *testing.T) {
			testMessageStore_AddMessage_NotExists(t, conveygo.NewMemoryStore(), alias, key)
		})
	})
	t.Run("GetMessage", func(t *testing.T) {
		t.Run("Exists", func(t *testing.T) {
			testMessageStore_GetMessage_Exists(t, conveygo.NewMemoryStore(), alias, key)
		})
		t.Run("Exists_Hash", func(t *testing.T) {
			testMessageStore_GetMessage_Exists_Hash(t, conveygo.NewMemoryStore(), alias, key)
		})
		t.Run("Exists_Reply", func(t *testing.T) {
			testMessageStore_GetMessage_Exists_Reply(t, conveygo.NewMemoryStore(), alias, key)
		})
		t.Run("NotExists", func(t *testing.T) {
			testMessageStore_GetMessage_NotExists(t, conveygo.NewMemoryStore())
		})
	})
	t.Run("GetYield", func(t *testing.T) {
		t.Run("Exists", func(t *testing.T) {
			testMessageStore_GetYield_Exists(t, conveygo.NewMemoryStore(), alias, key)
		})
		t.Run("Exists_Reply", func(t *testing.T) {
			testMessageStore_GetYield_Exists_Reply(t, conveygo.NewMemoryStore(), alias, key)
		})
		t.Run("NotExists", func(t *testing.T) {
			testMessageStore_GetYield_NotExists(t, conveygo.NewMemoryStore())
		})
	})
}
