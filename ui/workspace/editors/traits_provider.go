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

package editors

import (
	"github.com/richardwilkes/gcs/v5/model/gurps"
	"github.com/richardwilkes/gcs/v5/model/gurps/gid"
	"github.com/richardwilkes/gcs/v5/model/jio"
	"github.com/richardwilkes/gcs/v5/res"
	"github.com/richardwilkes/gcs/v5/ui/widget"
	"github.com/richardwilkes/gcs/v5/ui/widget/ntable"
	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/toolbox/log/jot"
	"github.com/richardwilkes/toolbox/txt"
	"github.com/richardwilkes/unison"
	"golang.org/x/exp/maps"
)

var (
	traitListColMap = map[int]int{
		0: gurps.TraitDescriptionColumn,
		1: gurps.TraitPointsColumn,
		2: gurps.TraitTagsColumn,
		3: gurps.TraitReferenceColumn,
	}
	traitPageColMap = map[int]int{
		0: gurps.TraitDescriptionColumn,
		1: gurps.TraitPointsColumn,
		2: gurps.TraitReferenceColumn,
	}
	_ ntable.TableProvider[*gurps.Trait] = &traitsProvider{}
)

type traitsProvider struct {
	table    *unison.Table[*ntable.Node[*gurps.Trait]]
	colMap   map[int]int
	provider gurps.TraitListProvider
	forPage  bool
}

// NewTraitsProvider creates a new table provider for traits.
func NewTraitsProvider(provider gurps.TraitListProvider, forPage bool) ntable.TableProvider[*gurps.Trait] {
	p := &traitsProvider{
		provider: provider,
		forPage:  forPage,
	}
	if forPage {
		p.colMap = traitPageColMap
	} else {
		p.colMap = traitListColMap
	}
	return p
}

func (p *traitsProvider) RefKey() string {
	return gurps.BlockLayoutTraitsKey
}

func (p *traitsProvider) Tags() []string {
	set := make(map[string]struct{})
	gurps.Traverse(func(trait *gurps.Trait) bool {
		for _, tag := range trait.Tags {
			set[tag] = struct{}{}
		}
		return false
	}, false, false, p.RootData()...)
	tags := maps.Keys(set)
	txt.SortStringsNaturalAscending(tags)
	return tags
}

func (p *traitsProvider) FilterByTag(tag string) []*gurps.Trait {
	var traits []*gurps.Trait
	gurps.Traverse(func(trait *gurps.Trait) bool {
		for _, one := range trait.Tags {
			if one == tag {
				traits = append(traits, trait)
				break
			}
		}
		return false
	}, false, false, p.RootData()...)
	return traits
}

func (p *traitsProvider) SetTable(table *unison.Table[*ntable.Node[*gurps.Trait]]) {
	p.table = table
}

func (p *traitsProvider) RootRowCount() int {
	return len(p.provider.TraitList())
}

func (p *traitsProvider) RootRows() []*ntable.Node[*gurps.Trait] {
	data := p.provider.TraitList()
	rows := make([]*ntable.Node[*gurps.Trait], 0, len(data))
	for _, one := range data {
		rows = append(rows, ntable.NewNode[*gurps.Trait](p.table, nil, p.colMap, one, p.forPage))
	}
	return rows
}

func (p *traitsProvider) SetRootRows(rows []*ntable.Node[*gurps.Trait]) {
	p.provider.SetTraitList(ntable.ExtractNodeDataFromList(rows))
}

func (p *traitsProvider) RootData() []*gurps.Trait {
	return p.provider.TraitList()
}

func (p *traitsProvider) SetRootData(data []*gurps.Trait) {
	p.provider.SetTraitList(data)
}

func (p *traitsProvider) Entity() *gurps.Entity {
	return p.provider.Entity()
}

func (p *traitsProvider) DragKey() string {
	return gid.Trait
}

func (p *traitsProvider) DragSVG() *unison.SVG {
	return res.GCSTraitsSVG
}

func (p *traitsProvider) DropShouldMoveData(from, to *unison.Table[*ntable.Node[*gurps.Trait]]) bool {
	return from == to
}

func (p *traitsProvider) ProcessDropData(_, _ *unison.Table[*ntable.Node[*gurps.Trait]]) {
}

func (p *traitsProvider) AltDropSupport() *ntable.AltDropSupport {
	return &ntable.AltDropSupport{
		DragKey: gid.TraitModifier,
		Drop: func(rowIndex int, data any) {
			if tableDragData, ok := data.(*unison.TableDragData[*ntable.Node[*gurps.TraitModifier]]); ok {
				entity := p.Entity()
				rows := make([]*gurps.TraitModifier, 0, len(tableDragData.Rows))
				for _, row := range tableDragData.Rows {
					rows = append(rows, row.Data().Clone(entity, nil, false))
				}
				rowData := p.table.RowFromIndex(rowIndex).Data()
				rowData.Modifiers = append(rowData.Modifiers, rows...)
				p.table.SyncToModel()
				if entity != nil {
					if rebuilder := unison.Ancestor[widget.Rebuildable](p.table); rebuilder != nil {
						rebuilder.Rebuild(true)
					}
					ntable.ProcessModifiers(p.table, rows)
					ntable.ProcessNameables(p.table, rows)
				}
			}
		},
	}
}

func (p *traitsProvider) ItemNames() (singular, plural string) {
	return i18n.Text("Trait"), i18n.Text("Traits")
}

func (p *traitsProvider) Headers() []unison.TableColumnHeader[*ntable.Node[*gurps.Trait]] {
	var headers []unison.TableColumnHeader[*ntable.Node[*gurps.Trait]]
	for i := 0; i < len(p.colMap); i++ {
		switch p.colMap[i] {
		case gurps.TraitDescriptionColumn:
			headers = append(headers, NewHeader[*gurps.Trait](i18n.Text("Trait"), "", p.forPage))
		case gurps.TraitPointsColumn:
			headers = append(headers, NewHeader[*gurps.Trait](i18n.Text("Pts"), i18n.Text("Points"), p.forPage))
		case gurps.TraitTagsColumn:
			headers = append(headers, NewHeader[*gurps.Trait](i18n.Text("Tags"), "", p.forPage))
		case gurps.TraitReferenceColumn:
			headers = append(headers, NewPageRefHeader[*gurps.Trait](p.forPage))
		default:
			jot.Fatalf(1, "invalid trait column: %d", p.colMap[i])
		}
	}
	return headers
}

func (p *traitsProvider) SyncHeader(_ []unison.TableColumnHeader[*ntable.Node[*gurps.Trait]]) {
}

func (p *traitsProvider) HierarchyColumnIndex() int {
	for k, v := range p.colMap {
		if v == gurps.TraitDescriptionColumn {
			return k
		}
	}
	return 0
}

func (p *traitsProvider) ExcessWidthColumnIndex() int {
	return p.HierarchyColumnIndex()
}

func (p *traitsProvider) OpenEditor(owner widget.Rebuildable, table *unison.Table[*ntable.Node[*gurps.Trait]]) {
	ntable.OpenEditor[*gurps.Trait](table, func(item *gurps.Trait) { EditTrait(owner, item) })
}

func (p *traitsProvider) CreateItem(owner widget.Rebuildable, table *unison.Table[*ntable.Node[*gurps.Trait]], variant ntable.ItemVariant) {
	item := gurps.NewTrait(p.Entity(), nil, variant == ntable.ContainerItemVariant)
	ntable.InsertItems[*gurps.Trait](owner, table, p.provider.TraitList, p.provider.SetTraitList,
		func(_ *unison.Table[*ntable.Node[*gurps.Trait]]) []*ntable.Node[*gurps.Trait] { return p.RootRows() }, item)
	EditTrait(owner, item)
}

func (p *traitsProvider) Serialize() ([]byte, error) {
	return jio.SerializeAndCompress(p.provider.TraitList())
}

func (p *traitsProvider) Deserialize(data []byte) error {
	var rows []*gurps.Trait
	if err := jio.DecompressAndDeserialize(data, &rows); err != nil {
		return err
	}
	p.provider.SetTraitList(rows)
	return nil
}

func (p *traitsProvider) ContextMenuItems() []ntable.ContextMenuItem {
	var list []ntable.ContextMenuItem
	list = append(list, ntable.TraitExtraContextMenuItems...)
	return append(list, ntable.DefaultContextMenuItems...)
}
