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
	"github.com/AletheiaWareLLC/cryptogo"
	"github.com/golang/protobuf/proto"
)

const (
	CONVEY_HOUR           = "Convey-Hour"         // Hourly Validation Chain
	CONVEY_DAY            = "Convey-Day"          // Daily Validation Chain
	CONVEY_WEEK           = "Convey-Week"         // Weekly Validation Chain
	CONVEY_YEAR           = "Convey-Year"         // Yearly Validation Chain
	CONVEY_DECADE         = "Convey-Decade"       // Decennially Validation Chain
	CONVEY_CENTURY        = "Convey-Century"      // Centennially Validation Chain
	CONVEY_CHARGE         = "Convey-Charge"       // financego.Charge Chain
	CONVEY_INVOICE        = "Convey-Invoice"      // financego.Invoice Chain
	CONVEY_REGISTRATION   = "Convey-Registration" // financego.Registration Chain
	CONVEY_SUBSCRIPTION   = "Convey-Subscription" // financego.Subscription Chain
	CONVEY_CONVERSATION   = "Convey-Conversation" // conveygo.Conversation Chain
	CONVEY_TRANSACTION    = "Convey-Transaction"  // conveygo.Transaction Chain
	CONVEY_PREFIX         = "Convey-"
	CONVEY_PREFIX_MESSAGE = "Convey-Message-" // conveygo.Message Chain
	CONVEY_PREFIX_TAG     = "Convey-Tag-"     // conveygo.Tag Chain
)

func GetConveyHosts() []string {
	if bcgo.IsDebug() {
		return []string{
			"test-convey.aletheiaware.com",
		}
	}
	return []string{
		"convey-nyc.aletheiaware.com",
		"convey-sfo.aletheiaware.com",
	}
}

func OpenHourChannel() *bcgo.Channel {
	return bcgo.OpenPoWChannel(CONVEY_HOUR, bcgo.THRESHOLD_PERIOD_HOUR)
}

func OpenDayChannel() *bcgo.Channel {
	return bcgo.OpenPoWChannel(CONVEY_DAY, bcgo.THRESHOLD_PERIOD_DAY)
}

func OpenWeekChannel() *bcgo.Channel {
	return bcgo.OpenPoWChannel(CONVEY_WEEK, bcgo.THRESHOLD_PERIOD_WEEK)
}

func OpenYearChannel() *bcgo.Channel {
	return bcgo.OpenPoWChannel(CONVEY_YEAR, bcgo.THRESHOLD_PERIOD_YEAR)
}

func OpenDecadeChannel() *bcgo.Channel {
	return bcgo.OpenPoWChannel(CONVEY_DECADE, bcgo.THRESHOLD_PERIOD_DECADE)
}

func OpenCenturyChannel() *bcgo.Channel {
	return bcgo.OpenPoWChannel(CONVEY_CENTURY, bcgo.THRESHOLD_PERIOD_CENTURY)
}

func OpenChargeChannel() *bcgo.Channel {
	// TODO(v2) add validator to ensure Message Payload can be unmarshalled as protobuf
	return bcgo.OpenPoWChannel(CONVEY_CHARGE, bcgo.THRESHOLD_G)
}

func OpenInvoiceChannel() *bcgo.Channel {
	// TODO(v2) add validator to ensure Message Payload can be unmarshalled as protobuf
	return bcgo.OpenPoWChannel(CONVEY_INVOICE, bcgo.THRESHOLD_G)
}

func OpenRegistrationChannel() *bcgo.Channel {
	// TODO(v2) add validator to ensure Message Payload can be unmarshalled as protobuf
	return bcgo.OpenPoWChannel(CONVEY_REGISTRATION, bcgo.THRESHOLD_G)
}

func OpenSubscriptionChannel() *bcgo.Channel {
	// TODO(v2) add validator to ensure Message Payload can be unmarshalled as protobuf
	return bcgo.OpenPoWChannel(CONVEY_SUBSCRIPTION, bcgo.THRESHOLD_G)
}

func OpenConversationChannel() *bcgo.Channel {
	// TODO(v2) add validator to ensure Message Payload can be unmarshalled as protobuf
	return bcgo.OpenPoWChannel(CONVEY_CONVERSATION, bcgo.THRESHOLD_G)
}

func OpenTransactionChannel() *bcgo.Channel {
	// TODO(v2) add validator to ensure Message Payload can be unmarshalled as protobuf
	transactions := bcgo.OpenPoWChannel(CONVEY_TRANSACTION, bcgo.THRESHOLD_G)
	transactions.AddValidator(&TransactionValidator{})
	return transactions
}

func OpenMessageChannel(conversationId string) *bcgo.Channel {
	// TODO(v2) add validator to ensure Message Payload can be unmarshalled as protobuf
	return bcgo.OpenPoWChannel(CONVEY_PREFIX_MESSAGE+conversationId, bcgo.THRESHOLD_G)
}

func OpenTagChannel(messageId string) *bcgo.Channel {
	// TODO(v2) add validator to ensure Message Payload can be unmarshalled as protobuf
	return bcgo.OpenPoWChannel(CONVEY_PREFIX_TAG+messageId, bcgo.THRESHOLD_G)
}

func ConversationEntryToListing(entry *bcgo.BlockEntry) (*Listing, error) {
	record := entry.Record
	// Unmarshal Protobuf
	c := &Conversation{}
	if err := proto.Unmarshal(record.Payload, c); err != nil {
		return nil, err
	}
	// Create Listing
	return &Listing{
		Hash:      entry.RecordHash,
		Timestamp: record.Timestamp,
		Author:    record.Creator,
		Topic:     c.Topic,
		Cost:      Cost(record),
	}, nil
}

func ProtoToRecord(alias string, key *rsa.PrivateKey, timestamp uint64, protobuf proto.Message) ([]byte, *bcgo.Record, error) {
	// Marshal Protobuf
	data, err := proto.Marshal(protobuf)
	if err != nil {
		return nil, nil, err
	}

	// Create Record
	_, record, err := bcgo.CreateRecord(timestamp, alias, key, nil, nil, data)
	if err != nil {
		return nil, nil, err
	}

	hash, err := cryptogo.HashProtobuf(record)
	if err != nil {
		return nil, nil, err
	}

	return hash, record, nil
}
