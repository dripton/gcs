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

package com.trollworks.gcs.model.criteria;

import com.trollworks.toolkit.collections.TKEnumExtractor;

/** The allowed string comparison types. */
public enum CMStringCompareType {
	/** The comparison for "is anything". */
	IS_ANYTHING(Msgs.IS_ANYTHING) {
		@Override public boolean matches(String qualifier, String data) {
			return true;
		}
	},
	/** The comparison for "is". */
	IS(Msgs.IS) {
		@Override public boolean matches(String qualifier, String data) {
			return data.equalsIgnoreCase(qualifier);
		}
	},
	/** The comparison for "is not". */
	IS_NOT(Msgs.IS_NOT) {
		@Override public boolean matches(String qualifier, String data) {
			return !data.equalsIgnoreCase(qualifier);
		}
	},
	/** The comparison for "contains". */
	CONTAINS(Msgs.CONTAINS) {
		@Override public boolean matches(String qualifier, String data) {
			return data.toLowerCase().indexOf(qualifier.toLowerCase()) != -1;
		}
	},
	/** The comparison for "does not contain". */
	DOES_NOT_CONTAIN(Msgs.DOES_NOT_CONTAIN) {
		@Override public boolean matches(String qualifier, String data) {
			return data.toLowerCase().indexOf(qualifier.toLowerCase()) == -1;
		}
	},
	/** The comparison for "starts with". */
	STARTS_WITH(Msgs.STARTS_WITH) {
		@Override public boolean matches(String qualifier, String data) {
			return data.toLowerCase().startsWith(qualifier.toLowerCase());
		}
	},
	/** The comparison for "does not start with". */
	DOES_NOT_START_WITH(Msgs.DOES_NOT_START_WITH) {
		@Override public boolean matches(String qualifier, String data) {
			return !data.toLowerCase().startsWith(qualifier.toLowerCase());
		}
	},
	/** The comparison for "ends with". */
	ENDS_WITH(Msgs.ENDS_WITH) {
		@Override public boolean matches(String qualifier, String data) {
			return data.toLowerCase().endsWith(qualifier.toLowerCase());
		}
	},
	/** The comparison for "does not end with". */
	DOES_NOT_END_WITH(Msgs.DOES_NOT_END_WITH) {
		@Override public boolean matches(String qualifier, String data) {
			return !data.toLowerCase().endsWith(qualifier.toLowerCase());
		}
	};

	private String	mTitle;

	private CMStringCompareType(String title) {
		mTitle = title;
	}

	@Override public String toString() {
		return mTitle;
	}

	/**
	 * @param qualifier The qualifier.
	 * @return The description of this comparison type.
	 */
	public String describe(String qualifier) {
		StringBuilder builder = new StringBuilder();

		builder.append(mTitle);
		builder.append(" \""); //$NON-NLS-1$
		builder.append(qualifier);
		builder.append('"');
		return builder.toString();
	}

	/**
	 * Performs a comparison.
	 * 
	 * @param qualifier The qualifier to use in conjuction with this {@link CMStringCompareType}.
	 * @param data The data to check.
	 * @return Whether the data matches the criteria or not.
	 */
	public abstract boolean matches(String qualifier, String data);

	/**
	 * @param buffer The buffer to load from.
	 * @return The units representing the buffer's description.
	 */
	public static final CMStringCompareType get(String buffer) {
		CMStringCompareType result = (CMStringCompareType) TKEnumExtractor.extract(buffer, values());

		if (result == null) {
			// Check a few others, for legacy reasons
			if ("starts".equalsIgnoreCase(buffer)) { //$NON-NLS-1$
				result = STARTS_WITH;
			} else if ("ends".equalsIgnoreCase(buffer)) { //$NON-NLS-1$
				result = ENDS_WITH;
			} else {
				result = IS;
			}
		}

		return result;
	}
}
