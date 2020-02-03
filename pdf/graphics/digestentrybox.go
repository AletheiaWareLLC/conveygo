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

package graphics

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/AletheiaWareLLC/conveygo"
	"github.com/AletheiaWareLLC/pdfgo"
	"github.com/AletheiaWareLLC/pdfgo/font"
	pdfgraphics "github.com/AletheiaWareLLC/pdfgo/graphics"
)

const ENTRY_PADDING = 10

type DigestEntryBox struct {
	pdfgraphics.Rectangle
	Host   string
	Level  int
	Entry  *conveygo.DigestEntry
	Fonts  map[string]font.Font
	Layout pdfgraphics.Layout
}

func (b *DigestEntryBox) SetBounds(bounds *pdfgraphics.Rectangle) error {
	b.Left = bounds.Left
	b.Top = bounds.Top
	b.Right = bounds.Right
	b.Bottom = bounds.Bottom

	t := b.Entry.Message.GetType()
	switch t {
	case conveygo.MediaType_TEXT_PLAIN:
		b.Layout = &pdfgraphics.ListLayout{
			Direction: pdfgraphics.TopBottom,
			Padding:   ENTRY_PADDING,
		}
		b.Layout.Add(&pdfgraphics.TextBox{
			Text:       []rune(b.Entry.Topic),
			FontId:     "F1",
			Font:       b.Fonts["F1"],
			FontSize:   32 - (2 * float64(b.Level)),
			FontColour: DARK_SKY_BLUE,
			Align:      pdfgraphics.Center,
		})
		b.Layout.Add(&pdfgraphics.TextBox{
			Text:       []rune(fmt.Sprintf("%s %s %d", b.Entry.Timestamp, b.Entry.Author, b.Entry.Yield)),
			FontId:     "F2",
			Font:       b.Fonts["F2"],
			FontSize:   12 - float64(b.Level),
			FontColour: LIGHT_SKY_BLUE,
			Align:      pdfgraphics.Center,
		})
		b.Layout.Add(&pdfgraphics.TextBox{
			Text:       []rune(string(b.Entry.Message.Content)),
			FontId:     "F3",
			Font:       b.Fonts["F3"],
			FontSize:   16 - float64(b.Level),
			FontColour: BLACK,
			Align:      pdfgraphics.JustifiedLeft,
		})
	default:
		return errors.New(fmt.Sprintf(conveygo.ERROR_UNRECOGNIZED_MEDIA_TYPE, t))
	}
	b.Layout.SetBounds(&pdfgraphics.Rectangle{
		Left:   bounds.Left + ENTRY_PADDING,
		Top:    bounds.Top - ENTRY_PADDING,
		Right:  bounds.Right - ENTRY_PADDING,
		Bottom: bounds.Bottom + ENTRY_PADDING,
	})
	return nil
}

func (b *DigestEntryBox) Write(p *pdfgo.PDF, buffer *bytes.Buffer) error {
	if b.Entry.Hash != "" {
		// Hyperlink
		p.AddAnnotation(pdfgo.NewHyperlink(b.Left, b.Bottom, b.Right, b.Top, fmt.Sprintf("https://%s/conversation?hash=%s", b.Host, b.Entry.Hash)))
	}

	// Border
	// buffer.WriteString(fmt.Sprintf("%s %s %s %s re S\n", FloatToString(b.Left), FloatToString(b.Bottom), FloatToString(b.GetWidth()), FloatToString(b.GetHeight())))

	return b.Layout.Write(p, buffer)
}
