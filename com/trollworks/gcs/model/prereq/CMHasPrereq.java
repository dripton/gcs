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

package com.trollworks.gcs.model.prereq;

import com.trollworks.toolkit.io.xml.TKXMLReader;
import com.trollworks.toolkit.io.xml.TKXMLWriter;

/**
 * An abstract prerequisite class for whether or not the specific item is present.
 */
public abstract class CMHasPrereq extends CMPrereq {
	/** The "has" attribute name. */
	protected static final String	ATTRIBUTE_HAS	= "has";	//$NON-NLS-1$
	private boolean					mHas;

	/**
	 * Creates a new prerequisite.
	 * 
	 * @param parent The owning prerequisite list, if any.
	 */
	public CMHasPrereq(CMPrereqList parent) {
		super(parent);
		mHas = true;
	}

	/**
	 * Creates a copy of the specified prerequisite.
	 * 
	 * @param parent The owning prerequisite list, if any.
	 * @param prereq The prerequisite to clone.
	 */
	protected CMHasPrereq(CMPrereqList parent, CMHasPrereq prereq) {
		super(parent);
		mHas = prereq.mHas;
	}

	@Override public boolean equals(Object obj) {
		if (this == obj) {
			return true;
		}
		if (obj instanceof CMHasPrereq) {
			return mHas == ((CMHasPrereq) obj).mHas;
		}
		return false;
	}

	/**
	 * Loads the "has" attribute.
	 * 
	 * @param reader The XML reader to load from.
	 */
	protected void loadHasAttribute(TKXMLReader reader) {
		mHas = reader.isAttributeSet(ATTRIBUTE_HAS);
	}

	/**
	 * Writes the "has" attribute to the stream.
	 * 
	 * @param out The XML writer to use.
	 */
	protected void saveHasAttribute(TKXMLWriter out) {
		out.writeAttribute(ATTRIBUTE_HAS, mHas);
	}

	/**
	 * @return <code>true</code> if the specified criteria should exist, <code>false</code> if
	 *         it should not.
	 */
	public boolean has() {
		return mHas;
	}

	/**
	 * @param has <code>true</code> if the specified criteria should exist, <code>false</code>
	 *            if it should not.
	 */
	public void has(boolean has) {
		mHas = has;
	}
}
