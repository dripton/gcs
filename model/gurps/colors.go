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
	"context"
	"fmt"
	"io/fs"
	"strconv"
	"strings"
	"sync"

	"github.com/richardwilkes/gcs/v5/model/jio"
	"github.com/richardwilkes/json"
	"github.com/richardwilkes/toolbox"
	"github.com/richardwilkes/toolbox/errs"
	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/unison"
)

const colorsTypeKey = "theme_colors"

// Additional colors over and above what unison provides by default.
var (
	HeaderColor             = &unison.ThemeColor{Light: unison.RGB(43, 43, 43), Dark: unison.RGB(64, 64, 64)}
	HintColor               = &unison.ThemeColor{Light: unison.Grey, Dark: unison.RGB(64, 64, 64)}
	MarkerColor             = &unison.ThemeColor{Light: unison.RGB(252, 242, 196), Dark: unison.RGB(0, 51, 0)}
	OnHeaderColor           = &unison.ThemeColor{Light: unison.White, Dark: unison.Silver}
	OnMarkerColor           = &unison.ThemeColor{Light: unison.Black, Dark: unison.RGB(221, 221, 221)}
	OnOverloadedColor       = &unison.ThemeColor{Light: unison.White, Dark: unison.RGB(221, 221, 221)}
	OnPageColor             = &unison.ThemeColor{Light: unison.Black, Dark: unison.RGB(160, 160, 160)}
	OnPageStandoutColor     = &unison.ThemeColor{Light: unison.Black, Dark: unison.RGB(160, 160, 160)}
	OnSearchListColor       = &unison.ThemeColor{Light: unison.Black, Dark: unison.RGB(204, 204, 204)}
	OverloadedColor         = &unison.ThemeColor{Light: unison.RGB(192, 64, 64), Dark: unison.RGB(115, 37, 37)}
	PDFLinkHighlightColor   = &unison.ThemeColor{Light: unison.SpringGreen, Dark: unison.SpringGreen}
	PDFMarkerHighlightColor = &unison.ThemeColor{Light: unison.Yellow, Dark: unison.Yellow}
	PageColor               = &unison.ThemeColor{Light: unison.White, Dark: unison.RGB(16, 16, 16)}
	PageStandoutColor       = &unison.ThemeColor{Light: unison.RGB(221, 221, 221), Dark: unison.RGB(64, 64, 64)}
	SearchListColor         = &unison.ThemeColor{Light: unison.LightCyan, Dark: unison.RGB(0, 43, 43)}
)

var (
	colorsOnce    sync.Once
	currentColors []*ThemedColor
	factoryColors []*ThemedColor
)

// ThemedColor holds a themed color.
type ThemedColor struct {
	ID    string
	Title string
	Color *unison.ThemeColor
}

// Colors holds a set of themed colors.
type Colors struct {
	data map[string]*unison.ThemeColor // Just here for serialization
}

type colorsData struct {
	Type    string `json:"type"`
	Version int    `json:"version"`
	Colors
}

// CurrentColors returns the current theme.
func CurrentColors() []*ThemedColor {
	colorsOnce.Do(initColors)
	return currentColors
}

// FactoryColors returns the original theme before any modifications.
func FactoryColors() []*ThemedColor {
	colorsOnce.Do(initColors)
	return factoryColors
}

func initColors() {
	currentColors = []*ThemedColor{
		{ID: "primary", Title: i18n.Text("Primary"), Color: &unison.PrimaryTheme.Primary},
		{ID: "on_primary", Title: i18n.Text("On Primary"), Color: &unison.PrimaryTheme.OnPrimary},
		{ID: "secondary", Title: i18n.Text("Secondary"), Color: &unison.PrimaryTheme.Secondary},
		{ID: "on_secondary", Title: i18n.Text("On Secondary"), Color: &unison.PrimaryTheme.OnSecondary},
		{ID: "tertiary", Title: i18n.Text("Tertiary"), Color: &unison.PrimaryTheme.Tertiary},
		{ID: "on_tertiary", Title: i18n.Text("On Tertiary"), Color: &unison.PrimaryTheme.OnTertiary},
		{ID: "surface", Title: i18n.Text("Surface"), Color: &unison.PrimaryTheme.Surface},
		{ID: "on_surface", Title: i18n.Text("On Surface"), Color: &unison.PrimaryTheme.OnSurface},
		{ID: "surface_above", Title: i18n.Text("Surface Above"), Color: &unison.PrimaryTheme.SurfaceAbove},
		// {ID: "on_surface_above", Title: i18n.Text("On Surface Above"), Color: &unison.PrimaryTheme.OnSurfaceAbove},
		{ID: "surface_below", Title: i18n.Text("Surface Below"), Color: &unison.PrimaryTheme.SurfaceBelow},
		// {ID: "on_surface_below", Title: i18n.Text("On Surface Below"), Color: &unison.PrimaryTheme.OnSurfaceBelow},
		{ID: "error", Title: i18n.Text("Error"), Color: &unison.PrimaryTheme.Error},
		{ID: "on_error", Title: i18n.Text("On Error"), Color: &unison.PrimaryTheme.OnError},
		{ID: "warning", Title: i18n.Text("Warning"), Color: &unison.PrimaryTheme.Warning},
		{ID: "on_warning", Title: i18n.Text("On Warning"), Color: &unison.PrimaryTheme.OnWarning},
		{ID: "outline", Title: i18n.Text("Outline"), Color: &unison.PrimaryTheme.Outline},
		{ID: "outline_variant", Title: i18n.Text("Outline Variant"), Color: &unison.PrimaryTheme.OutlineVariant},

		{ID: "header", Title: i18n.Text("Header"), Color: HeaderColor},
		{ID: "on_header", Title: i18n.Text("On Header"), Color: OnHeaderColor},
		{ID: "hint", Title: i18n.Text("Hint"), Color: HintColor},
		{ID: "search_list", Title: i18n.Text("Search List"), Color: SearchListColor},
		{ID: "on_search_list", Title: i18n.Text("On Search List"), Color: OnSearchListColor},
		{ID: "marker", Title: i18n.Text("Marker"), Color: MarkerColor},
		{ID: "on_marker", Title: i18n.Text("On Marker"), Color: OnMarkerColor},
		{ID: "overloaded", Title: i18n.Text("Overloaded"), Color: OverloadedColor},
		{ID: "on_overloaded", Title: i18n.Text("On Overloaded"), Color: OnOverloadedColor},
		{ID: "page", Title: i18n.Text("Page"), Color: PageColor},
		{ID: "on_page", Title: i18n.Text("On Page"), Color: OnPageColor},
		{ID: "page_standout", Title: i18n.Text("Page Standout"), Color: PageStandoutColor},
		{ID: "on_page_standout", Title: i18n.Text("On Page Standout"), Color: OnPageStandoutColor},
		{ID: "pdf_link", Title: i18n.Text("PDF Link Highlight"), Color: PDFLinkHighlightColor},
		{ID: "pdf_marker", Title: i18n.Text("PDF Marker Highlight"), Color: PDFMarkerHighlightColor},

		{ID: "shadow", Title: i18n.Text("Shadow"), Color: &unison.PrimaryTheme.Shadow},
	}
	factoryColors = make([]*ThemedColor, len(currentColors))
	for i, c := range currentColors {
		factoryColors[i] = &ThemedColor{
			ID:    c.ID,
			Title: c.Title,
			Color: &unison.ThemeColor{
				Light: c.Color.Light,
				Dark:  c.Color.Dark,
			},
		}
	}
}

// NewColorsFromFS creates a new set of colors from a file. Any missing values will be filled in with defaults.
func NewColorsFromFS(fileSystem fs.FS, filePath string) (*Colors, error) {
	var current colorsData
	if err := jio.LoadFromFS(context.Background(), fileSystem, filePath, &current); err != nil {
		return nil, errs.Wrap(err)
	}
	switch current.Version {
	case 0:
		// During development of v5, forgot to add the type & version initially, so try and fix that up
		if current.Type == "" {
			current.Type = colorsTypeKey
			current.Version = CurrentDataVersion
		}
	case 1:
		current.Type = colorsTypeKey
		current.Version = CurrentDataVersion
	default:
	}
	if current.Type != colorsTypeKey {
		return nil, errs.New(unexpectedFileDataMsg())
	}
	if err := CheckVersion(current.Version); err != nil {
		return nil, err
	}
	return &current.Colors, nil
}

// Save writes the Colors to the file as JSON.
func (c *Colors) Save(filePath string) error {
	return jio.SaveToFile(context.Background(), filePath, &colorsData{
		Type:    colorsTypeKey,
		Version: CurrentDataVersion,
		Colors:  *c,
	})
}

// MarshalJSON implements json.Marshaler.
func (c *Colors) MarshalJSON() ([]byte, error) {
	c.data = make(map[string]*unison.ThemeColor, len(CurrentColors()))
	for _, one := range CurrentColors() {
		c.data[one.ID] = one.Color
	}
	return json.Marshal(&c.data)
}

// UnmarshalJSON implements json.Unmarshaler.
func (c *Colors) UnmarshalJSON(data []byte) error {
	c.data = make(map[string]*unison.ThemeColor, len(CurrentColors()))
	var err error
	toolbox.CallWithHandler(func() {
		err = json.Unmarshal(data, &c.data)
	}, func(e error) {
		err = e
	})
	if err != nil {
		var old map[string]string
		if err = json.Unmarshal(data, &old); err != nil {
			return errs.New("invalid color data")
		}
		c.data = make(map[string]*unison.ThemeColor, len(CurrentColors()))
		for _, fc := range CurrentColors() {
			local := *fc.Color
			c.data[fc.ID] = &local
			if cc, ok := old[fc.ID]; ok {
				var clr unison.Color
				if clr, err = unison.ColorDecode(cc); err != nil {
					if clr, err = unison.ColorDecode("rgb(" + cc + ")"); err != nil {
						lastComma := strings.LastIndexByte(cc, ',')
						if lastComma == -1 {
							continue
						}
						var v int
						if v, err = strconv.Atoi(cc[lastComma+1:]); err != nil {
							continue
						}
						if clr, err = unison.ColorDecode(fmt.Sprintf("rgba(%s,%f)", cc[:lastComma], float32(v)/255)); err != nil {
							continue
						}
					}
				}
				local.Light = clr
				local.Dark = clr
			}
		}
	}
	if c.data == nil {
		c.data = make(map[string]*unison.ThemeColor, len(CurrentColors()))
	}
	for _, one := range FactoryColors() {
		if _, ok := c.data[one.ID]; !ok {
			clr := *one.Color
			c.data[one.ID] = &clr
		}
	}
	return nil
}

// MakeCurrent applies these colors to the current theme color set and updates all windows.
func (c *Colors) MakeCurrent() {
	for _, one := range CurrentColors() {
		if v, ok := c.data[one.ID]; ok {
			*one.Color = *v
		}
	}
	unison.ThemeChanged()
}

// Reset to factory defaults.
func (c *Colors) Reset() {
	for _, one := range FactoryColors() {
		*c.data[one.ID] = *one.Color
	}
}

// ResetOne resets one color by ID to factory defaults.
func (c *Colors) ResetOne(id string) {
	for _, v := range FactoryColors() {
		if v.ID == id {
			*c.data[id] = *v.Color
			break
		}
	}
}
