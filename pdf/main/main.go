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

package main

import (
	"encoding/base64"
	"flag"
	"github.com/AletheiaWareLLC/bcgo"
	"github.com/AletheiaWareLLC/conveygo"
	"github.com/AletheiaWareLLC/conveygo/pdf"
	"github.com/AletheiaWareLLC/pdfgo"
	"github.com/AletheiaWareLLC/pdfgo/font"
	"log"
	"os"
	"path"
)

var (
	mock          = flag.Bool("mock", false, "mock digest entries")
	host          = flag.String("host", "test-convey.aletheiaware.com", "Convey host")
	fontfamily    = flag.String("fontfamily", "Times", "ttf font family")
	fontdirectory = flag.String("fontdirectory", "/usr/share/fonts/", "ttf font directory")
)

func main() {
	var err error

	flag.Parse()

	var entries []*conveygo.DigestEntry
	if *mock {
		entries = GetMockDigestEntries()
	} else {
		entries, err = GetDigestEntries(*host)
		if err != nil {
			log.Fatal(err)
		}
	}

	if len(entries) < conveygo.DIGEST_LIMIT {
		log.Fatal("Too few entries for digest")
	}

	writer := os.Stdout
	args := flag.Args()
	if len(args) > 0 {
		log.Println("Writing:", args[0])
		file, err := os.OpenFile(args[0], os.O_CREATE|os.O_TRUNC|os.O_WRONLY, os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		writer = file
	}

	p := pdfgo.NewPDF()

	var fonts map[string]font.Font
	switch *fontfamily {
	case "Courier":
		fonts, err = LoadCoreFont(p, map[string]string{
			"F1": "Courier-Bold",
			"F2": "Courier-Oblique",
			"F3": "Courier",
		})
	case "Helvetica":
		fonts, err = LoadCoreFont(p, map[string]string{
			"F1": "Helvetica-Bold",
			"F2": "Helvetica-Oblique",
			"F3": "Helvetica",
		})
	case "Times":
		fonts, err = LoadCoreFont(p, map[string]string{
			"F1": "Times-Bold",
			"F2": "Times-Italic",
			"F3": "Times-Roman",
		})
	default:
		fonts, err = LoadTTFFont(p, map[string]string{
			"F1": path.Join(*fontdirectory, *fontfamily+" Bold.ttf"),
			"F2": path.Join(*fontdirectory, *fontfamily+" Italic.ttf"),
			"F3": path.Join(*fontdirectory, *fontfamily+".ttf"),
		})
	}
	if err != nil {
		log.Fatal(err)
	}

	err = pdf.AddEntries(p, *host, entries, fonts)
	if err != nil {
		log.Fatal(err)
	}

	err = p.Write(writer)
	if err != nil {
		log.Fatal(err)
	}
}

func LoadCoreFont(p *pdfgo.PDF, ids map[string]string) (map[string]font.Font, error) {
	fonts := make(map[string]font.Font)
	for id, name := range ids {
		font, err := font.NewCoreFont(p, name)
		if err != nil {
			return nil, err
		}
		fonts[id] = font
	}
	return fonts, nil
}

func LoadTTFFont(p *pdfgo.PDF, ids map[string]string) (map[string]font.Font, error) {
	fonts := make(map[string]font.Font)
	for id, file := range ids {
		font, err := font.NewTrueTypeFont(p, file)
		if err != nil {
			return nil, err
		}
		fonts[id] = font
	}
	return fonts, nil
}

func GetDigestEntries(host string) ([]*conveygo.DigestEntry, error) {
	rootDir, err := bcgo.GetRootDirectory()
	if err != nil {
		return nil, err
	}
	log.Println("Root Directory:", rootDir)

	cacheDir, err := bcgo.GetCacheDirectory(rootDir)
	if err != nil {
		return nil, err
	}
	log.Println("Cache Directory:", cacheDir)

	cache, err := bcgo.NewFileCache(cacheDir)
	if err != nil {
		return nil, err
	}

	peers, err := bcgo.GetPeers(rootDir)
	if err != nil {
		return nil, err
	}
	peers = append(peers, host)
	log.Println("Peers:", peers)

	network := &bcgo.TcpNetwork{
		Peers: peers,
	}

	node, err := bcgo.GetNode(rootDir, cache, network)
	if err != nil {
		return nil, err
	}

	conversations := conveygo.OpenConversationChannel()

	if err := conversations.LoadCachedHead(cache); err != nil {
		log.Println(err)
	}

	if err := conversations.Pull(cache, network); err != nil {
		log.Println(err)
	}

	node.AddChannel(conversations)

	if err := bcgo.Iterate(conversations.Name, conversations.Head, nil, cache, network, func(h []byte, b *bcgo.Block) error {
		for _, entry := range b.Entry {
			name := conveygo.CONVEY_PREFIX_MESSAGE + base64.RawURLEncoding.EncodeToString(entry.RecordHash)
			channel := bcgo.OpenPoWChannel(name, bcgo.THRESHOLD_G)

			if err := channel.LoadCachedHead(cache); err != nil {
				log.Println(err)
			}

			if err := channel.Pull(cache, network); err != nil {
				log.Println(err)
			}

			node.AddChannel(channel)
		}
		return nil
	}); err != nil {
		return nil, err
	}

	messages := &conveygo.BCStore{
		Node:     node,
		Listener: &bcgo.PrintingMiningListener{Output: os.Stdout},
	}

	return conveygo.GetDigestEntries(messages, 0, bcgo.Timestamp())
}

func GetMockDigestEntries() []*conveygo.DigestEntry {
	return []*conveygo.DigestEntry{
		&conveygo.DigestEntry{
			Topic:     "First Article",
			Timestamp: bcgo.TimestampToString(bcgo.Timestamp()),
			Author:    "Alice",
			Yield:     55,
			Message: &conveygo.Message{
				Content: []byte(LOREM_IPSUM),
				Type:    conveygo.MediaType_TEXT_PLAIN,
			},
		},
		&conveygo.DigestEntry{
			Topic:     "Second Article",
			Timestamp: bcgo.TimestampToString(bcgo.Timestamp()),
			Author:    "Bob",
			Yield:     34,
			Message: &conveygo.Message{
				Content: []byte(LOREM_IPSUM),
				Type:    conveygo.MediaType_TEXT_PLAIN,
			},
		},
		&conveygo.DigestEntry{
			Topic:     "Long Article is Loooooooooooooooong",
			Timestamp: bcgo.TimestampToString(bcgo.Timestamp()),
			Author:    "Charlie",
			Yield:     21,
			Message: &conveygo.Message{
				Content: []byte(LOREM_IPSUM),
				Type:    conveygo.MediaType_TEXT_PLAIN,
			},
		},
		&conveygo.DigestEntry{
			Topic:     "Blah",
			Timestamp: bcgo.TimestampToString(bcgo.Timestamp()),
			Author:    "Daniel",
			Yield:     13,
			Message: &conveygo.Message{
				Content: []byte(LOREM_IPSUM),
				Type:    conveygo.MediaType_TEXT_PLAIN,
			},
		},
	}
}

const (
	LOREM_IPSUM = `Lorem ipsum dolor sit amet, consectetur adipiscing elit. Pellentesque quis consectetur nisi. Suspendisse id interdum felis. Sed egestas eget tellus eu pharetra. Praesent pulvinar sed massa id placerat. Etiam sem libero, semper vitae consequat ut, volutpat id mi. Mauris volutpat pellentesque convallis. Curabitur rutrum venenatis orci nec ornare. Maecenas quis pellentesque neque. Aliquam consectetur dapibus nulla, id maximus odio ultrices ac. Sed luctus at felis sed faucibus. Cras leo augue, congue in velit ut, mattis rhoncus lectus.

Praesent viverra, mauris ut ullamcorper semper, leo urna auctor lectus, vitae vehicula mi leo quis lorem. Nullam condimentum, massa at tempor feugiat, metus enim lobortis velit, eget suscipit eros ipsum quis tellus. Aenean fermentum diam vel felis dictum semper. Duis nisl orci, tincidunt ut leo quis, luctus vehicula diam. Sed velit justo, congue id augue eu, euismod dapibus lacus. Proin sit amet imperdiet sapien. Mauris erat urna, fermentum et quam rhoncus, fringilla consequat ante. Vivamus consectetur molestie odio, ac rutrum erat finibus a. Suspendisse id maximus felis. Sed mauris odio, mattis eget mi eu, consequat tempus purus.

Nulla facilisi. In a condimentum dolor. Class aptent taciti sociosqu ad litora torquent per conubia nostra, per inceptos himenaeos. Maecenas ultrices, justo vel commodo cursus, sapien neque rhoncus turpis, vitae gravida lorem nisl in tortor. Vestibulum et consectetur sem. Morbi dapibus maximus vulputate. Maecenas nunc lacus, posuere ut gravida id, rhoncus ut velit.

Mauris ullamcorper nec leo quis elementum. Nam vel purus consequat, sodales lacus auctor, interdum massa. Vestibulum non efficitur elit. Mauris laoreet nunc ultricies purus condimentum, egestas hendrerit mauris mollis. Aliquam tincidunt eros sed leo eleifend, varius consequat leo volutpat. Integer eget ultricies lorem, ac vulputate velit. Maecenas eu magna mauris.

Nullam eu mattis dolor. Sed sit amet ipsum gravida, pretium justo eget, mattis est. Cras viverra aliquet velit, a faucibus urna luctus vel. Donec vehicula turpis ligula, non auctor justo tempus nec. In libero orci, tempus vitae ante eu, convallis dapibus nulla. Donec egestas volutpat elit vel semper. Lorem ipsum dolor sit amet, consectetur adipiscing elit. Curabitur enim orci, scelerisque ut viverra at, condimentum sed augue. Duis ut nulla dapibus, mattis ex ac, maximus neque. Integer a augue justo. Ut facilisis enim diam, condimentum interdum elit pulvinar tincidunt. Donec et ipsum ac turpis tristique efficitur. Morbi sodales, odio et vehicula eleifend, quam metus commodo lacus, nec gravida lectus augue eget neque. Quisque tortor lectus, viverra non elit et, hendrerit ornare purus. Cras pretium, felis vel gravida fermentum, libero metus dapibus ex, nec mollis ante enim eget dolor.

Nam nibh nulla, ullamcorper id tortor at, lacinia molestie lacus. Sed eget orci ac sem efficitur placerat mattis id arcu. Praesent leo neque, accumsan eget venenatis vel, accumsan at justo. Maecenas venenatis blandit varius. Fusce ut luctus velit. Nullam quis risus enim. Mauris auctor tempor fermentum. Nam nec leo rutrum nibh aliquam ultrices. Nunc ut molestie lacus, vel feugiat tellus. Suspendisse fringilla eu diam eget interdum.

Nunc eget scelerisque nunc. Quisque faucibus mi in magna ullamcorper porttitor. Ut sodales semper est, ut laoreet felis viverra dapibus. Fusce convallis finibus dapibus. Ut tincidunt at urna non imperdiet. Vivamus efficitur lorem a eros cursus, et tempor lectus malesuada. Sed ac lectus quis enim condimentum commodo eu eget lectus.

Pellentesque sollicitudin nisi at sapien egestas aliquet. Nulla sit amet maximus urna, in laoreet odio. In quis leo fringilla, mollis nisl et, aliquam mi. Nulla gravida tempus justo in pulvinar. Vestibulum nunc libero, accumsan id molestie sed, elementum a enim. Vivamus volutpat fermentum risus, quis faucibus nulla laoreet sed. Quisque id posuere velit. Nam blandit non orci eu mattis. Donec dui ipsum, scelerisque ut vulputate sed, cursus vel velit. Maecenas bibendum neque elit, ac placerat ex pulvinar in. Morbi ultricies est ac malesuada hendrerit. Pellentesque et mauris aliquet neque congue ultricies a vitae libero.

Fusce vitae malesuada ipsum, eget viverra nisi. Nullam varius rhoncus pellentesque. In in ligula est. Suspendisse id felis gravida, cursus justo ac, ultrices tellus. Suspendisse potenti. Quisque sit amet enim vitae dolor pulvinar tempus sit amet eget magna. Aliquam condimentum sapien eu lectus feugiat fringilla. Curabitur fringilla, metus id egestas condimentum, massa nibh lacinia velit, eu pulvinar lacus mi sed sem. Nulla volutpat tincidunt lacinia. Nunc tristique ipsum et nulla finibus egestas. Etiam luctus tempor metus a consectetur.

Ut ac pulvinar purus. Pellentesque tellus quam, condimentum at odio id, viverra porttitor tortor. Aliquam vulputate accumsan mattis. Maecenas quis enim in lorem elementum vehicula. Proin accumsan nec enim in pharetra. Phasellus id enim ligula. Vestibulum ante ipsum primis in faucibus orci luctus et ultrices posuere cubilia Curae; Mauris urna sem, euismod id tempor at, consectetur eu sem. Quisque vestibulum hendrerit diam. Duis gravida tristique velit et sodales. Etiam feugiat blandit diam tempor tincidunt.
`
)
