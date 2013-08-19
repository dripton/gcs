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

package com.trollworks.gcs.model;

import com.trollworks.toolkit.undo.TKSimpleUndo;
import com.trollworks.toolkit.undo.TKUndoException;

/** Provides undo support for character fields. */
public class CMCharacterFieldUndo extends TKSimpleUndo {
	private CMCharacter	mCharacter;
	private String		mName;
	private String		mID;
	private Object		mBefore;
	private Object		mAfter;

	/**
	 * Create a new character field undo edit.
	 * 
	 * @param character The character to provide an undo edit for.
	 * @param name The name of the undo edit.
	 * @param id The ID of the field being changed.
	 * @param before The original value.
	 * @param after The new value.
	 */
	public CMCharacterFieldUndo(CMCharacter character, String name, String id, Object before, Object after) {
		super();
		mCharacter = character;
		mName = name;
		mID = id;
		mBefore = before;
		mAfter = after;
	}

	@Override public String getName() {
		return mName;
	}

	@Override public void apply(boolean forUndo) throws TKUndoException {
		super.apply(forUndo);
		mCharacter.setValueForID(mID, forUndo ? mBefore : mAfter);
	}
}
