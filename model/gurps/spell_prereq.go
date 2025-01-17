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
	"github.com/richardwilkes/gcs/v5/model/fxp"
	"github.com/richardwilkes/gcs/v5/model/gurps/enums/prereq"
	"github.com/richardwilkes/gcs/v5/model/gurps/enums/spellcmp"
	"github.com/richardwilkes/toolbox/xio"
)

var _ Prereq = &SpellPrereq{}

// SpellPrereq holds a prerequisite for a spell.
type SpellPrereq struct {
	Parent            *PrereqList     `json:"-"`
	Type              prereq.Type     `json:"type"`
	SubType           spellcmp.Type   `json:"sub_type"`
	Has               bool            `json:"has"`
	QualifierCriteria StringCriteria  `json:"qualifier,omitempty"`
	QuantityCriteria  NumericCriteria `json:"quantity,omitempty"`
}

// NewSpellPrereq creates a new SpellPrereq.
func NewSpellPrereq() *SpellPrereq {
	return &SpellPrereq{
		Type:    prereq.Spell,
		SubType: spellcmp.Name,
		QualifierCriteria: StringCriteria{
			StringCriteriaData: StringCriteriaData{
				Compare: IsString,
			},
		},
		QuantityCriteria: NumericCriteria{
			NumericCriteriaData: NumericCriteriaData{
				Compare:   AtLeastNumber,
				Qualifier: fxp.One,
			},
		},
		Has: true,
	}
}

// PrereqType implements Prereq.
func (s *SpellPrereq) PrereqType() prereq.Type {
	return s.Type
}

// ParentList implements Prereq.
func (s *SpellPrereq) ParentList() *PrereqList {
	return s.Parent
}

// Clone implements Prereq.
func (s *SpellPrereq) Clone(parent *PrereqList) Prereq {
	clone := *s
	clone.Parent = parent
	return &clone
}

// FillWithNameableKeys implements Prereq.
func (s *SpellPrereq) FillWithNameableKeys(m map[string]string) {
	if s.SubType.UsesStringCriteria() {
		Extract(s.QualifierCriteria.Qualifier, m)
	}
}

// ApplyNameableKeys implements Prereq.
func (s *SpellPrereq) ApplyNameableKeys(m map[string]string) {
	if s.SubType.UsesStringCriteria() {
		s.QualifierCriteria.Qualifier = Apply(s.QualifierCriteria.Qualifier, m)
	}
}

// Satisfied implements Prereq.
func (s *SpellPrereq) Satisfied(entity *Entity, exclude any, tooltip *xio.ByteBuffer, prefix string, _ *bool) bool {
	var techLevel *string
	if sp, ok := exclude.(*Spell); ok {
		techLevel = sp.TechLevel
	}
	count := 0
	colleges := make(map[string]bool)
	Traverse(func(sp *Spell) bool {
		if exclude == sp || sp.Points == 0 {
			return false
		}
		if techLevel != nil && sp.TechLevel != nil && *techLevel != *sp.TechLevel {
			return false
		}
		switch s.SubType {
		case spellcmp.Name:
			if s.QualifierCriteria.Matches(sp.Name) {
				count++
			}
		case spellcmp.Tag:
			for _, one := range sp.Tags {
				if s.QualifierCriteria.Matches(one) {
					count++
					break
				}
			}
		case spellcmp.College:
			for _, one := range sp.College {
				if s.QualifierCriteria.Matches(one) {
					count++
					break
				}
			}
		case spellcmp.CollegeCount:
			for _, one := range sp.College {
				colleges[one] = true
			}
		case spellcmp.Any:
			count++
		}
		return false
	}, false, true, entity.Spells...)
	if s.SubType == spellcmp.CollegeCount {
		count = len(colleges)
	}
	satisfied := s.QuantityCriteria.Matches(fxp.From(count))
	if !s.Has {
		satisfied = !satisfied
	}
	if !satisfied && tooltip != nil {
		tooltip.WriteString(prefix)
		tooltip.WriteString(HasText(s.Has))
		tooltip.WriteByte(' ')
		if s.SubType == spellcmp.CollegeCount {
			tooltip.WriteString("college count which ")
			tooltip.WriteString(s.QuantityCriteria.String())
		} else {
			tooltip.WriteString(s.QuantityCriteria.String())
			if s.QuantityCriteria.Qualifier == fxp.One {
				tooltip.WriteString(" spell ")
			} else {
				tooltip.WriteString(" spells ")
			}
			if s.SubType == spellcmp.Any {
				tooltip.WriteString("of any kind")
			} else {
				switch s.SubType {
				case spellcmp.Name:
					tooltip.WriteString("whose name ")
				case spellcmp.Tag:
					tooltip.WriteString("whose tag ")
				case spellcmp.College:
					tooltip.WriteString("whose college ")
				}
				tooltip.WriteString(s.QualifierCriteria.String())
			}
		}
	}
	return satisfied
}
