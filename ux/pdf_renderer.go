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
	"image"
	"os"
	"sync"
	"time"

	"github.com/richardwilkes/pdf"
	"github.com/richardwilkes/toolbox/errs"
	"github.com/richardwilkes/unison"
)

// PDFTableOfContents holds a table of contents entry.
type PDFTableOfContents struct {
	Title        string
	PageNumber   int
	PageLocation unison.Point
	Children     []*PDFTableOfContents
}

// PDFPage holds a rendered PDFRenderer page.
type PDFPage struct {
	Error      error
	PageNumber int
	Image      *unison.Image
	TOC        []*PDFTableOfContents
	Links      []*PDFLink
	Matches    []unison.Rect
}

// PDFLink holds a single link on a page. If PageNumber if >= 0, then this is an internal link and the URI will be empty.
type PDFLink struct {
	Bounds     unison.Rect
	PageNumber int
	URI        string
}

type pdfParams struct {
	sequence   int
	pageNumber int
	search     string
	scale      float32
}

// PDFRenderer holds a PDFRenderer page renderer.
type PDFRenderer struct {
	MaxSearchMatches     int
	PPI                  int
	DisplayScaleAdjust   float32
	doc                  *pdf.Document
	pageCount            int
	pageLoadedCallback   func()
	lock                 sync.RWMutex
	page                 *PDFPage
	lastRequest          *pdfParams
	sequence             int
	lastRenderedSequence int
	lastRenderRequest    time.Time
}

// NewPDFRenderer creates a new PDFRenderer page renderer.
func NewPDFRenderer(filePath string, pageLoadedCallback func()) (*PDFRenderer, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, errs.Wrap(err)
	}
	var doc *pdf.Document
	if doc, err = pdf.New(data, 0); err != nil {
		return nil, errs.Wrap(err)
	}
	display := unison.PrimaryDisplay()
	return &PDFRenderer{
		MaxSearchMatches:   100,
		PPI:                display.PPI(),
		DisplayScaleAdjust: 1 / display.ScaleX,
		doc:                doc,
		pageCount:          doc.PageCount(),
		pageLoadedCallback: pageLoadedCallback,
	}, nil
}

// PageCount returns the total page count.
func (p *PDFRenderer) PageCount() int {
	return p.pageCount
}

// CurrentPage returns the currently rendered page.
func (p *PDFRenderer) CurrentPage() *PDFPage {
	p.lock.RLock()
	defer p.lock.RUnlock()
	return p.page
}

// MostRecentPageNumber returns the most recent page number that has been asked to be rendered.
func (p *PDFRenderer) MostRecentPageNumber() int {
	p.lock.RLock()
	defer p.lock.RUnlock()
	if p.lastRequest != nil {
		return p.lastRequest.pageNumber
	}
	if p.page != nil {
		return p.page.PageNumber
	}
	return 0
}

// LoadPage requests the given page to be loaded and rendered at the specified scale.
func (p *PDFRenderer) LoadPage(pageNumber int, scale float32, search string) {
	if pageNumber < 0 || pageNumber >= p.pageCount {
		return
	}
	p.lock.Lock()
	defer p.lock.Unlock()
	if p.lastRequest != nil && p.lastRequest.pageNumber == pageNumber && p.lastRequest.scale == scale &&
		p.lastRequest.search == search {
		return
	}
	if p.lastRequest == nil || p.lastRenderedSequence == p.lastRequest.sequence {
		p.lastRenderRequest = time.Now()
	}
	p.sequence++
	p.lastRequest = &pdfParams{
		sequence:   p.sequence,
		pageNumber: pageNumber,
		scale:      scale,
		search:     search,
	}
	submitPDF(p, false)
}

// RequestRenderPriority attempts to bump this PDFRenderer's rendering to the head of the queue.
func (p *PDFRenderer) RequestRenderPriority() {
	p.lock.Lock()
	defer p.lock.Unlock()
	if p.lastRequest != nil {
		submitPDF(p, true)
	}
}

func (p *PDFRenderer) render(state *pdfParams) {
	if p.shouldAbortRender(state) {
		return
	}

	dpi := int(state.scale * float32(p.PPI) / p.DisplayScaleAdjust)
	toc := p.doc.TableOfContents(dpi)
	if p.shouldAbortRender(state) {
		return
	}

	page, err := p.doc.RenderPage(state.pageNumber, dpi, p.MaxSearchMatches, state.search)
	if err != nil {
		p.errorDuringRender(state, err)
		return
	}
	if p.shouldAbortRender(state) {
		return
	}

	var img *unison.Image
	img, err = unison.NewImageFromPixels(page.Image.Rect.Dx(), page.Image.Rect.Dy(), page.Image.Pix, p.DisplayScaleAdjust)
	if err != nil {
		p.errorDuringRender(state, err)
		return
	}
	p.lock.RLock()
	displayScaleAdjust := p.DisplayScaleAdjust
	p.lock.RUnlock()
	if p.shouldAbortRender(state) {
		return
	}
	pg := &PDFPage{
		PageNumber: state.pageNumber,
		Image:      img,
		TOC:        convertPDFTableOfContentsEntries(toc, displayScaleAdjust),
		Links:      convertPDFLinks(page.Links, displayScaleAdjust),
		Matches:    convertPDFMatches(page.SearchHits, displayScaleAdjust),
	}
	p.lock.Lock()
	if state.sequence != p.lastRequest.sequence {
		p.lock.Unlock()
		return
	}
	p.page = pg
	p.lastRenderedSequence = state.sequence
	p.lock.Unlock()
	p.pageLoadedCallback()
}

// RenderingFinished returns true if there is no rendering being done for this PDFRenderer at the moment.
func (p *PDFRenderer) RenderingFinished() (finished bool, pageNumber int, requested time.Time) {
	p.lock.RLock()
	defer p.lock.RUnlock()
	if p.lastRequest != nil {
		finished = p.lastRenderedSequence == p.lastRequest.sequence
		pageNumber = p.lastRequest.pageNumber
		requested = p.lastRenderRequest
	} else {
		finished = true
		if p.page != nil {
			pageNumber = p.page.PageNumber
		}
		requested = time.Now()
	}
	return
}

func (p *PDFRenderer) shouldAbortRender(state *pdfParams) bool {
	p.lock.RLock()
	defer p.lock.RUnlock()
	return state.sequence != p.lastRequest.sequence
}

func (p *PDFRenderer) errorDuringRender(state *pdfParams, err error) {
	p.lock.Lock()
	if state.sequence != p.lastRequest.sequence {
		p.lock.Unlock()
		return
	}
	p.page = &PDFPage{
		Error:      err,
		PageNumber: state.pageNumber,
	}
	p.lock.Unlock()
	p.pageLoadedCallback()
}

func convertPDFTableOfContentsEntries(entries []*pdf.TOCEntry, displayScaleAdjust float32) []*PDFTableOfContents {
	if len(entries) == 0 {
		return nil
	}
	toc := make([]*PDFTableOfContents, len(entries))
	for i, entry := range entries {
		toc[i] = &PDFTableOfContents{
			Title:        entry.Title,
			PageNumber:   entry.PageNumber,
			PageLocation: pointFromPDFPagePoint(entry.PageX, entry.PageY, displayScaleAdjust),
			Children:     convertPDFTableOfContentsEntries(entry.Children, displayScaleAdjust),
		}
	}
	return toc
}

func convertPDFLinks(pageLinks []*pdf.PageLink, displayScaleAdjust float32) []*PDFLink {
	if len(pageLinks) == 0 {
		return nil
	}
	links := make([]*PDFLink, len(pageLinks))
	for i, link := range pageLinks {
		links[i] = &PDFLink{
			Bounds:     rectFromPDFPageRect(link.Bounds, displayScaleAdjust),
			PageNumber: link.PageNumber,
			URI:        link.URI,
		}
	}
	return links
}

func convertPDFMatches(hits []image.Rectangle, displayScaleAdjust float32) []unison.Rect {
	if len(hits) == 0 {
		return nil
	}
	matches := make([]unison.Rect, len(hits))
	for i, hit := range hits {
		matches[i] = rectFromPDFPageRect(hit, displayScaleAdjust)
	}
	return matches
}

func pointFromPDFPagePoint(x, y int, displayScaleAdjust float32) unison.Point {
	return unison.NewPoint(float32(x)*displayScaleAdjust, float32(y)*displayScaleAdjust)
}

func rectFromPDFPageRect(r image.Rectangle, displayScaleAdjust float32) unison.Rect {
	return unison.NewRect(float32(r.Min.X)*displayScaleAdjust, float32(r.Min.Y)*displayScaleAdjust,
		float32(r.Dx())*displayScaleAdjust, float32(r.Dy())*displayScaleAdjust)
}