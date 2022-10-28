/*
 * Copyright ©1998-2022 by Richard A. Wilkes. All rights reserved.
 *
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, version 2.0. If a copy of the MPL was not distributed with
 * this file, You can obtain one at http://mozilla.org/MPL/2.0/.
 *
 * This Source Code Form is "Incompatible With Secondary Licenses", as
 * defined by the Mozilla Public License, version 2.0.
 */

package sheet

import (
	"fmt"
	"strconv"

	"github.com/richardwilkes/gcs/v5/model/gurps"
	"github.com/richardwilkes/gcs/v5/model/gurps/datafile"
	"github.com/richardwilkes/gcs/v5/model/theme"
	"github.com/richardwilkes/gcs/v5/res"
	"github.com/richardwilkes/gcs/v5/ui/widget"
	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/unison"
)

var _ unison.ColorProvider = &encRowColor{}

// EncumbrancePanel holds the contents of the encumbrance block on the sheet.
type EncumbrancePanel struct {
	unison.Panel
	entity     *gurps.Entity
	row        []unison.Paneler
	current    int
	overloaded bool
}

// NewEncumbrancePanel creates a new encumbrance panel.
func NewEncumbrancePanel(entity *gurps.Entity) *EncumbrancePanel {
	p := &EncumbrancePanel{entity: entity}
	p.Self = p
	p.SetLayout(&unison.FlexLayout{
		Columns:  9,
		HSpacing: 4,
	})
	p.SetLayoutData(&unison.FlexLayoutData{
		HAlign: unison.FillAlignment,
		VAlign: unison.FillAlignment,
		HGrab:  true,
	})
	p.SetBorder(&widget.TitledBorder{Title: i18n.Text("Encumbrance, Move & Dodge")})
	p.DrawCallback = func(gc *unison.Canvas, rect unison.Rect) {
		r := p.Children()[0].FrameRect()
		r.X = rect.X
		r.Width = rect.Width
		gc.DrawRect(r, theme.HeaderColor.Paint(gc, r, unison.Fill))
		p.current = int(entity.EncumbranceLevel(true))
		p.overloaded = entity.WeightCarried(false) > entity.MaximumCarry(datafile.ExtraHeavy)
		for i, row := range p.row {
			var ink unison.Ink
			switch {
			case p.current == i:
				if p.overloaded {
					ink = theme.OverloadedColor
				} else {
					ink = theme.MarkerColor
				}
			case i&1 == 1:
				ink = unison.BandingColor
			default:
				ink = unison.ContentColor
			}
			r = row.AsPanel().FrameRect()
			r.X = rect.X
			r.Width = rect.Width
			gc.DrawRect(r, ink.Paint(gc, r, unison.Fill))
		}
	}

	p.AddChild(widget.NewPageHeader(i18n.Text("Level"), 3))
	p.AddChild(widget.NewInteriorSeparator())
	p.AddChild(widget.NewPageHeader(i18n.Text("Max Load"), 1))
	p.AddChild(widget.NewInteriorSeparator())
	p.AddChild(widget.NewPageHeader(i18n.Text("Move"), 1))
	p.AddChild(widget.NewInteriorSeparator())
	p.AddChild(widget.NewPageHeader(i18n.Text("Dodge"), 1))

	for i, enc := range datafile.AllEncumbrance {
		rowColor := &encRowColor{
			owner: p,
			index: i,
		}
		p.AddChild(p.createMarker(entity, enc, rowColor))
		p.AddChild(p.createLevelField(enc, rowColor))
		name := widget.NewPageLabel(enc.String())
		name.OnBackgroundInk = rowColor
		name.SetLayoutData(&unison.FlexLayoutData{
			HAlign: unison.FillAlignment,
			VAlign: unison.MiddleAlignment,
			HGrab:  true,
		})
		p.AddChild(name)
		p.row = append(p.row, name)
		if i == 0 {
			p.addSeparator()
		}
		p.AddChild(p.createMaxCarryField(enc, rowColor))
		if i == 0 {
			p.addSeparator()
		}
		p.AddChild(p.createMoveField(enc, rowColor))
		if i == 0 {
			p.addSeparator()
		}
		p.AddChild(p.createDodgeField(enc, rowColor))
	}

	return p
}

func (p *EncumbrancePanel) createMarker(entity *gurps.Entity, enc datafile.Encumbrance, rowColor *encRowColor) *unison.Label {
	marker := widget.NewPageLabel("")
	marker.OnBackgroundInk = rowColor
	marker.SetBorder(unison.NewEmptyBorder(unison.Insets{Left: 4}))
	baseline := marker.Font.Baseline()
	marker.Drawable = &unison.DrawableSVG{
		SVG:  res.WeightSVG,
		Size: unison.Size{Width: baseline, Height: baseline},
	}
	marker.DrawCallback = func(gc *unison.Canvas, rect unison.Rect) {
		if enc == entity.EncumbranceLevel(true) {
			marker.DefaultDraw(gc, rect)
		}
	}
	return marker
}

func (p *EncumbrancePanel) createLevelField(enc datafile.Encumbrance, rowColor *encRowColor) *widget.NonEditablePageField {
	field := widget.NewNonEditablePageFieldEnd(func(f *widget.NonEditablePageField) {})
	field.OnBackgroundInk = rowColor
	field.Text = strconv.Itoa(int(enc))
	field.SetLayoutData(&unison.FlexLayoutData{
		HAlign: unison.FillAlignment,
		VAlign: unison.MiddleAlignment,
	})
	return field
}

func (p *EncumbrancePanel) createMaxCarryField(enc datafile.Encumbrance, rowColor *encRowColor) *widget.NonEditablePageField {
	field := widget.NewNonEditablePageFieldEnd(func(f *widget.NonEditablePageField) {
		if text := p.entity.SheetSettings.DefaultWeightUnits.Format(p.entity.MaximumCarry(enc)); text != f.Text {
			f.Text = text
			widget.MarkForLayoutWithinDockable(f)
		}
	})
	field.OnBackgroundInk = rowColor
	field.Tooltip = unison.NewTooltipWithText(fmt.Sprintf(i18n.Text("The maximum load that can be carried and still remain within the %s encumbrance level"), enc.String()))
	return field
}

func (p *EncumbrancePanel) createMoveField(enc datafile.Encumbrance, rowColor *encRowColor) *widget.NonEditablePageField {
	field := widget.NewNonEditablePageFieldEnd(func(f *widget.NonEditablePageField) {
		if text := strconv.Itoa(p.entity.Move(enc)); text != f.Text {
			f.Text = text
			widget.MarkForLayoutWithinDockable(f)
		}
	})
	field.OnBackgroundInk = rowColor
	field.Tooltip = unison.NewTooltipWithText(fmt.Sprintf(i18n.Text("The ground movement rate for the %s encumbrance level"), enc.String()))
	return field
}

func (p *EncumbrancePanel) createDodgeField(enc datafile.Encumbrance, rowColor *encRowColor) *widget.NonEditablePageField {
	field := widget.NewNonEditablePageFieldEnd(func(f *widget.NonEditablePageField) {
		if text := strconv.Itoa(p.entity.Dodge(enc)); text != f.Text {
			f.Text = text
			widget.MarkForLayoutWithinDockable(f)
		}
	})
	field.OnBackgroundInk = rowColor
	field.Tooltip = unison.NewTooltipWithText(fmt.Sprintf(i18n.Text("The dodge for the %s encumbrance level"), enc.String()))
	field.SetBorder(unison.NewEmptyBorder(unison.Insets{Right: 4}))
	return field
}

func (p *EncumbrancePanel) addSeparator() {
	sep := unison.NewSeparator()
	sep.Vertical = true
	sep.LineInk = unison.InteriorDividerColor
	sep.SetLayoutData(&unison.FlexLayoutData{
		VSpan:  len(datafile.AllEncumbrance),
		HAlign: unison.MiddleAlignment,
		VAlign: unison.FillAlignment,
		VGrab:  true,
	})
	p.AddChild(sep)
}

type encRowColor struct {
	owner *EncumbrancePanel
	index int
}

func (c *encRowColor) GetColor() unison.Color {
	switch {
	case c.owner.current == c.index:
		if c.owner.overloaded {
			return theme.OnOverloadedColor.GetColor()
		}
		return theme.OnMarkerColor.GetColor()
	case c.index&1 == 1:
		return unison.OnBandingColor.GetColor()
	default:
		return unison.OnContentColor.GetColor()
	}
}

func (c *encRowColor) Paint(canvas *unison.Canvas, rect unison.Rect, style unison.PaintStyle) *unison.Paint {
	return c.GetColor().Paint(canvas, rect, style)
}
