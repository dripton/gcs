/*
 * Copyright ©1998-2023 by Richard A. Wilkes. All rights reserved.
 *
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, version 2.0. If a copy of the MPL was not distributed with
 * this file, You can obtain one at http://mozilla.org/MPL/2.0/.
 *
 * This Source Code Form is "Incompatible With Secondary Licenses", as
 * defined by the Mozilla Public License, version 2.0.
 */

package gurps

import (
	"github.com/richardwilkes/unison/enums/align"
)

// PageRefCellAlias is used an alias to request the page reference cell, if any.
const PageRefCellAlias = -10

// CellData holds data for creating a cell's visual representation.
type CellData struct {
	Type              CellType
	Disabled          bool
	Dim               bool
	Checked           bool
	Alignment         align.Enum
	Primary           string
	Secondary         string
	Tooltip           string
	UnsatisfiedReason string
	TemplateInfo      string
	InlineTag         string
}

// ForSort returns a string that can be used to sort or search against for this data.
func (c *CellData) ForSort() string {
	switch c.Type {
	case TextCellType:
		if c.Secondary != "" {
			return c.Primary + "\n" + c.Secondary
		}
		return c.Primary
	case ToggleCellType:
		if c.Checked {
			return "√"
		}
	case PageRefCellType, TagsCellType, MarkdownCellType:
		return c.Primary
	}
	return ""
}
