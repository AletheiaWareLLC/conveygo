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

package pdf

import (
	"bytes"
	"github.com/AletheiaWareLLC/conveygo"
	"github.com/AletheiaWareLLC/conveygo/pdf/graphics"
	"github.com/AletheiaWareLLC/pdfgo"
	"github.com/AletheiaWareLLC/pdfgo/font"
	pdfgraphics "github.com/AletheiaWareLLC/pdfgo/graphics"
)

func AddEntries(p *pdfgo.PDF, host string, entries []*conveygo.DigestEntry, fonts map[string]font.Font) error {
	// Resources
	fs := p.NewDictionaryObject()
	for id, font := range fonts {
		fs.AddNameObjectEntry(id, font.GetReference())
	}
	resources := p.NewDictionaryObject()
	resources.AddNameObjectEntry("Font", pdfgo.NewObjectReference(fs))

	pageWidth := 595.28
	pageHeight := 841.89

	// Contents
	sizes := []float64{806, 496, 310, 186, 124, 62, 62, 0}
	contentWidth := sizes[1]
	contentHeight := sizes[0]
	marginX := (pageWidth - contentWidth) / 2
	marginY := (pageHeight - contentHeight) / 2
	bounds := &pdfgraphics.Rectangle{
		Left:   marginX,
		Right:  marginX + contentWidth,
		Top:    marginY + contentHeight,
		Bottom: marginY,
	}

	layout := &pdfgraphics.FibonacciLayout{
		Sizes: sizes[1:],
	}
	for i := 0; i < conveygo.DIGEST_LIMIT; i++ {
		layout.Add(&graphics.DigestEntryBox{
			Host:  host,
			Level: i + 1,
			Entry: entries[i],
			Fonts: fonts,
		})
	}
	layout.Add(&pdfgraphics.GravityLayout{
		Gravity: pdfgraphics.Middle,
		Box: &pdfgraphics.TextBox{
			Text:       []rune("Convey"),
			FontId:     "F1",
			Font:       fonts["F1"],
			FontSize:   18,
			FontColour: graphics.DARK_SKY_BLUE,
			Align:      pdfgraphics.Center,
		},
	})
	layout.Add(&graphics.LogoBox{
		Colour: graphics.DARK_SKY_BLUE,
	})

	if err := layout.SetBounds(bounds); err != nil {
		return err
	}

	var buffer bytes.Buffer
	if err := layout.Write(p, &buffer); err != nil {
		return err
	}
	contents := p.NewStreamObject()
	contents.Data = buffer.Bytes()
	p.AddPage(pageWidth, pageHeight, pdfgo.NewObjectReference(resources), pdfgo.NewObjectReference(contents))
	return nil
}
