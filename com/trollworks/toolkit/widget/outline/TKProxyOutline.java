/* ***** BEGIN LICENSE BLOCK *****
 * Version: MPL 1.1
 *
 * The contents of this file are subject to the Mozilla Public License Version
 * 1.1 (the "License"); you may not use this file except in compliance with
 * the License. You may obtain a copy of the License at
 * http://www.mozilla.org/MPL/
 *
 * Software distributed under the License is distributed on an "AS IS" basis,
 * WITHOUT WARRANTY OF ANY KIND, either express or implied. See the License
 * for the specific language governing rights and limitations under the
 * License.
 *
 * The Original Code is GURPS Character Sheet.
 *
 * The Initial Developer of the Original Code is Richard A. Wilkes.
 * Portions created by the Initial Developer are Copyright (C) 1998-2002,
 * 2005-2007 the Initial Developer. All Rights Reserved.
 *
 * Contributor(s):
 *
 * ***** END LICENSE BLOCK ***** */

package com.trollworks.toolkit.widget.outline;

import com.trollworks.toolkit.widget.TKPanel;
import com.trollworks.toolkit.widget.menu.TKMenuItem;

import java.awt.Color;
import java.awt.Point;
import java.awt.Rectangle;
import java.awt.dnd.DropTargetDragEvent;
import java.awt.dnd.DropTargetDropEvent;
import java.awt.dnd.DropTargetEvent;
import java.awt.event.ActionEvent;
import java.util.List;

/** An outline which gets its data and selection from another {@link TKOutline}. */
public class TKProxyOutline extends TKOutline {
	private TKOutline	mOutlineToProxy;

	/**
	 * Creates a new proxy outline.
	 * 
	 * @param outlineToProxy The {@link TKOutline} to proxy.
	 */
	public TKProxyOutline(TKOutline outlineToProxy) {
		super(outlineToProxy.getModel(), outlineToProxy.showIndent());
		mOutlineToProxy = outlineToProxy;
		mOutlineToProxy.addProxy(this);
	}

	/** @return The outline being proxied. */
	@Override public TKOutline getRealOutline() {
		return mOutlineToProxy;
	}

	@Override protected void addProxy(TKProxyOutline proxy) {
		assert false : "You cannot add a proxy to a proxy"; //$NON-NLS-1$
	}

	@Override protected void convertDragRowsToSelf(List<TKRow> list) {
		mOutlineToProxy.convertDragRowsToSelf(list);
	}

	@Override protected boolean isDragAcceptable(DropTargetDragEvent dtde) {
		return mOutlineToProxy.isDragAcceptable(dtde);
	}

	@Override public TKPanel getCurrentEditor() {
		return mOutlineToProxy.getCurrentEditor();
	}

	@Override public void stopEditing() {
		mOutlineToProxy.stopEditing();
	}

	@Override public Color getBackground(int rowIndex, boolean selected, boolean active) {
		return mOutlineToProxy.getBackground(rowIndex, selected, active);
	}

	@Override public Point getCellLocationOfEditor() {
		return mOutlineToProxy.getCellLocationOfEditor();
	}

	@Override public void contentSizeMayHaveChanged() {
		mOutlineToProxy.contentSizeMayHaveChanged();
	}

	@Override public void notifyActionListeners(ActionEvent event) {
		if (mOutlineToProxy != null) {
			mOutlineToProxy.notifyActionListeners(event);
		} else {
			super.notifyActionListeners(event);
		}
	}

	@Override public String getActionCommand() {
		if (mOutlineToProxy != null) {
			return mOutlineToProxy.getActionCommand();
		}
		return super.getActionCommand();
	}

	@Override public String getSelectionChangedActionCommand() {
		if (mOutlineToProxy != null) {
			return mOutlineToProxy.getSelectionChangedActionCommand();
		}
		return super.getSelectionChangedActionCommand();
	}

	@Override public String getDeleteSelectionActionCommand() {
		if (mOutlineToProxy != null) {
			return mOutlineToProxy.getDeleteSelectionActionCommand();
		}
		return super.getDeleteSelectionActionCommand();
	}

	@Override public String getPotentialContentSizeChangeActionCommand() {
		if (mOutlineToProxy != null) {
			return mOutlineToProxy.getPotentialContentSizeChangeActionCommand();
		}
		return super.getPotentialContentSizeChangeActionCommand();
	}

	@Override public void repaint(Rectangle bounds) {
		if (mOutlineToProxy != null) {
			mOutlineToProxy.repaint(bounds);
		}
	}

	@Override void repaintHeader(Rectangle bounds) {
		if (mOutlineToProxy != null) {
			mOutlineToProxy.repaintHeader(bounds);
		}
	}

	@Override public void repaintSelection() {
		mOutlineToProxy.repaintSelection();
	}

	@Override protected void repaintProxyRow(TKRow row) {
		mOutlineToProxy.repaintProxyRow(row);
	}

	/**
	 * Repaint the specified area.
	 * 
	 * @param bounds The area to repaint.
	 */
	protected void repaintProxy(Rectangle bounds) {
		super.repaint(bounds);
	}

	@Override protected void rowsWereDropped() {
		mOutlineToProxy.rowsWereDropped();
	}

	@Override protected int dragEnterRow(DropTargetDragEvent dtde) {
		mOutlineToProxy.dragEnterRow(dtde, this);
		return super.dragEnterRow(dtde);
	}

	@Override protected void dragExitRow(DropTargetEvent dte) {
		mOutlineToProxy.dragExitRow(dte, this);
		super.dragExitRow(dte);
	}

	@Override protected void dropRow(DropTargetDropEvent dtde) {
		mOutlineToProxy.dropRow(dtde, this);
		super.dropRow(dtde);
	}

	@Override public boolean adjustMenuItem(String command, TKMenuItem item) {
		return mOutlineToProxy.adjustMenuItem(command, item);
	}

	@Override public boolean obeyCommand(String command, TKMenuItem item) {
		return mOutlineToProxy.obeyCommand(command, item);
	}
}
