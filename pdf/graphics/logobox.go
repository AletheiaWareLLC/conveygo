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
	"fmt"
	"github.com/AletheiaWareLLC/pdfgo"
	pdfgraphics "github.com/AletheiaWareLLC/pdfgo/graphics"
)

const MAGIC = 0.551784

type LogoBox struct {
	pdfgraphics.Rectangle
	Colour []float64
}

func (b *LogoBox) SetBounds(bounds *pdfgraphics.Rectangle) error {
	b.Left = bounds.Left
	b.Top = bounds.Top
	b.Right = bounds.Right
	b.Bottom = bounds.Bottom
	return nil
}

func (b *LogoBox) Write(p *pdfgo.PDF, buffer *bytes.Buffer) error {
	buffer.WriteString("q\n")

	// Hyperlink
	p.AddAnnotation(pdfgo.NewHyperlink(b.Left, b.Bottom, b.Right, b.Top, "https://aletheiaware.com"))

	s := pdfgraphics.FloatToString

	// Border
	// buffer.WriteString(fmt.Sprintf("%s %s %s %s re S\n", s(b.Left), s(b.Bottom), s(b.GetWidth()), s(b.GetHeight())))

	buffer.WriteString(fmt.Sprintf("%s %s %s RG\n", s(b.Colour[0]), s(b.Colour[1]), s(b.Colour[2])))
	buffer.WriteString("2 w\n")

	width := b.GetWidth()
	height := b.GetHeight()
	limit := width
	if height < width {
		limit = height
	}

	// Reduce by 25% for padding
	limit *= 0.75

	// Circle
	radius := limit / 2
	magic := radius * MAGIC
	// Center
	cx := b.Left + (width / 2)
	cy := b.Bottom + (height / 2)
	// Point 1
	x1 := cx - radius
	y1 := cy
	// Point 2
	x2 := cx
	y2 := cy + radius
	// Point 3
	x3 := cx + radius
	y3 := cy
	// Point 4
	x4 := cx
	y4 := cy - radius

	// Control 1
	cx1 := x1
	cy1 := y1 + magic
	// Control 2
	cx2 := x2 - magic
	cy2 := y2
	// Control 3
	cx3 := x2 + magic
	cy3 := y2
	// Control 4
	cx4 := x3
	cy4 := y3 + magic
	// Control 5
	cx5 := x3
	cy5 := y3 - magic
	// Control 6
	cx6 := x4 + magic
	cy6 := y4
	// Control 7
	cx7 := x4 - magic
	cy7 := y4
	// Control 8
	cx8 := x1
	cy8 := y1 - magic

	buffer.WriteString(fmt.Sprintf("%s %s m\n", s(x1), s(y1)))
	buffer.WriteString(fmt.Sprintf("%s %s %s %s %s %s c\n", s(cx1), s(cy1), s(cx2), s(cy2), s(x2), s(y2)))
	buffer.WriteString(fmt.Sprintf("%s %s %s %s %s %s c\n", s(cx3), s(cy3), s(cx4), s(cy4), s(x3), s(y3)))
	buffer.WriteString(fmt.Sprintf("%s %s %s %s %s %s c\n", s(cx5), s(cy5), s(cx6), s(cy6), s(x4), s(y4)))
	buffer.WriteString(fmt.Sprintf("%s %s %s %s %s %s c S\n", s(cx7), s(cy7), s(cx8), s(cy8), s(x1), s(y1)))

	pd1 := limit / 10
	pd2 := pd1 + pd1

	px1 := x1
	py1 := y1
	px2 := px1 + pd2
	py2 := py1 + pd2
	px3 := px2 + pd2
	py3 := py1
	px4 := px3 + pd2
	py4 := py1 - pd2
	px5 := px4 + pd1
	py5 := py1 - pd1
	px6 := px5 + pd1
	py6 := py1 - pd2
	px7 := px6 + pd2
	py7 := py1

	// Path
	buffer.WriteString(fmt.Sprintf("%s %s m\n", s(px1), s(py1)))
	buffer.WriteString(fmt.Sprintf("%s %s l\n", s(px2), s(py2)))
	buffer.WriteString(fmt.Sprintf("%s %s l\n", s(px3), s(py3)))
	buffer.WriteString(fmt.Sprintf("%s %s l\n", s(px4), s(py4)))
	buffer.WriteString(fmt.Sprintf("%s %s l\n", s(px5), s(py5)))
	buffer.WriteString(fmt.Sprintf("%s %s l\n", s(px6), s(py6)))
	buffer.WriteString(fmt.Sprintf("%s %s l S\n", s(px7), s(py7)))
	buffer.WriteString("Q\n")
	return nil
}
