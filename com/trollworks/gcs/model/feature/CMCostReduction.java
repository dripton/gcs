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

package com.trollworks.gcs.model.feature;

import com.trollworks.gcs.model.CMCharacter;
import com.trollworks.toolkit.io.xml.TKXMLNodeType;
import com.trollworks.toolkit.io.xml.TKXMLReader;
import com.trollworks.toolkit.io.xml.TKXMLWriter;

import java.io.IOException;
import java.util.HashMap;
import java.util.HashSet;

/** Describes a cost reduction. */
public class CMCostReduction implements CMFeature {
	/** The XML tag. */
	public static final String	TAG_ROOT		= "cost_reduction"; //$NON-NLS-1$
	private static final String	TAG_ATTRIBUTE	= "attribute";		//$NON-NLS-1$
	private static final String	TAG_PERCENTAGE	= "percentage";	//$NON-NLS-1$
	/** The ST attribute. */
	public static final String	ST				= "ST";			//$NON-NLS-1$
	/** The DX attribute. */
	public static final String	DX				= "DX";			//$NON-NLS-1$
	/** The IQ attribute. */
	public static final String	IQ				= "IQ";			//$NON-NLS-1$
	/** The HT attribute. */
	public static final String	HT				= "HT";			//$NON-NLS-1$
	private String				mAttribute;
	private int					mPercentage;

	/** Creates a new cost reduction. */
	public CMCostReduction() {
		mAttribute = ST;
		mPercentage = 40;
	}

	/**
	 * Creates a clone of the specified cost reduction.
	 * 
	 * @param other The bonus to clone.
	 */
	public CMCostReduction(CMCostReduction other) {
		mAttribute = other.mAttribute;
		mPercentage = other.mPercentage;
	}

	/**
	 * Loads a {@link CMCostReduction}.
	 * 
	 * @param reader The XML reader to use.
	 * @throws IOException
	 */
	public CMCostReduction(TKXMLReader reader) throws IOException {
		this();
		load(reader);
	}

	@Override public boolean equals(Object obj) {
		if (this == obj) {
			return true;
		}
		if (obj instanceof CMCostReduction) {
			CMCostReduction other = (CMCostReduction) obj;

			return mPercentage == other.mPercentage && mAttribute.equals(other.mAttribute);
		}
		return false;
	}

	/** @return The percentage to use. */
	public int getPercentage() {
		return mPercentage;
	}

	/** @param percentage The percentage to use. */
	public void setPercentage(int percentage) {
		mPercentage = percentage;
	}

	/** @return The attribute this cost reduction applies to. */
	public String getAttribute() {
		return mAttribute;
	}

	/** @param attribute The attribute. */
	public void setAttribute(String attribute) {
		if (ST.equals(attribute)) {
			mAttribute = ST;
		} else if (DX.equals(attribute)) {
			mAttribute = DX;
		} else if (IQ.equals(attribute)) {
			mAttribute = IQ;
		} else if (HT.equals(attribute)) {
			mAttribute = HT;
		} else {
			mAttribute = ST;
		}
	}

	/**
	 * Loads a cost reduction.
	 * 
	 * @param reader The XML reader to use.
	 * @throws IOException
	 */
	protected void load(TKXMLReader reader) throws IOException {
		String marker = reader.getMarker();

		do {
			if (reader.next() == TKXMLNodeType.START_TAG) {
				String name = reader.getName();

				if (TAG_ATTRIBUTE.equals(name)) {
					setAttribute(reader.readText());
				} else if (TAG_PERCENTAGE.equals(name)) {
					setPercentage(reader.readInteger(0));
				} else {
					reader.skipTag(name);
				}
			}
		} while (reader.withinMarker(marker));
	}

	public String getXMLTag() {
		return TAG_ROOT;
	}

	public String getKey() {
		return CMCharacter.ATTRIBUTES_PREFIX + mAttribute;
	}

	public CMFeature cloneFeature() {
		return new CMCostReduction(this);
	}

	public void save(TKXMLWriter out) {
		out.startSimpleTagEOL(TAG_ROOT);
		out.simpleTag(TAG_ATTRIBUTE, mAttribute);
		out.simpleTag(TAG_PERCENTAGE, mPercentage);
		out.endTagEOL(TAG_ROOT, true);
	}

	public void fillWithNameableKeys(HashSet<String> set) {
		// Nothing to do.
	}

	public void applyNameableKeys(HashMap<String, String> map) {
		// Nothing to do.
	}
}
