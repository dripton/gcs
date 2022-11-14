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

package ux

import (
	"github.com/richardwilkes/gcs/v5/model"
	"github.com/richardwilkes/gcs/v5/model/jio"
	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/unison"
)

// MiscPanel holds the contents of the miscellaneous block on the sheet.
type MiscPanel struct {
	unison.Panel
	entity    *model.Entity
	targetMgr *TargetMgr
	prefix    string
}

// NewMiscPanel creates a new miscellaneous panel.
func NewMiscPanel(entity *model.Entity, targetMgr *TargetMgr) *MiscPanel {
	m := &MiscPanel{
		entity:    entity,
		targetMgr: targetMgr,
		prefix:    targetMgr.NextPrefix(),
	}
	m.Self = m
	m.SetLayout(&unison.FlexLayout{
		Columns:  2,
		HSpacing: 4,
	})
	m.SetLayoutData(&unison.FlexLayoutData{
		HAlign: unison.FillAlignment,
		VAlign: unison.FillAlignment,
	})
	m.SetBorder(unison.NewCompoundBorder(&TitledBorder{Title: i18n.Text("Miscellaneous")}, unison.NewEmptyBorder(unison.Insets{
		Top:    1,
		Left:   2,
		Bottom: 1,
		Right:  2,
	})))
	m.DrawCallback = func(gc *unison.Canvas, rect unison.Rect) {
		gc.DrawRect(rect, unison.ContentColor.Paint(gc, rect, unison.Fill))
	}

	m.AddChild(NewPageLabelEnd(i18n.Text("Created")))
	m.AddChild(NewNonEditablePageField(func(f *NonEditablePageField) {
		if text := m.entity.CreatedOn.String(); text != f.Text {
			f.Text = text
			MarkForLayoutWithinDockable(f)
		}
	}))

	m.AddChild(NewPageLabelEnd(i18n.Text("Modified")))
	m.AddChild(NewNonEditablePageField(func(f *NonEditablePageField) {
		if text := m.entity.ModifiedOn.String(); text != f.Text {
			f.Text = text
			MarkForLayoutWithinDockable(f)
		}
	}))

	title := i18n.Text("Player")
	m.AddChild(NewPageLabelEnd(title))
	m.AddChild(NewStringPageFieldNoGrab(m.targetMgr, m.prefix+"player", title,
		func() string { return m.entity.Profile.PlayerName },
		func(s string) { m.entity.Profile.PlayerName = s }))

	return m
}

// UpdateModified updates the current modification timestamp.
func (m *MiscPanel) UpdateModified() {
	m.entity.ModifiedOn = jio.Now()
}

// SetTextAndMarkModified sets the field to the given text, selects it, requests focus, then calls MarkModified().
func SetTextAndMarkModified(field *unison.Field, text string) {
	field.SetText(text)
	field.SelectAll()
	field.RequestFocus()
	field.Parent().MarkForLayoutAndRedraw()
	MarkModified(field)
}
