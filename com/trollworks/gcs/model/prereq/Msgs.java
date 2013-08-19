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

import com.trollworks.toolkit.io.TKMessages;

/** All localized strings in this package should be accessed via this class. */
class Msgs {
	/** Used by {@link CMAdvantagePrereq}. */
	public static String	NAME_PART;
	/** Used by {@link CMAdvantagePrereq}. */
	public static String	NOTES_PART;

	/** Used by {@link CMAttributePrereq}. */
	public static String	DESCRIPTION;
	/** Used by {@link CMAttributePrereq}. */
	public static String	COMBINED;

	/** Used by {@link CMContainedWeightPrereq}. */
	public static String	CONTAINED_WEIGHT;

	/** Used by {@link CMPrereqList}. */
	public static String	REQUIRES_ALL;
	/** Used by {@link CMPrereqList}. */
	public static String	REQUIRES_ANY;

	/** Used by {@link CMSkillPrereq}. */
	public static String	SKILL_NAME_PART;
	/** Used by {@link CMSkillPrereq}. */
	public static String	SPECIALIZATION_PART;
	/** Used by {@link CMSkillPrereq}. */
	public static String	LEVEL_AND_TL_PART;

	/** Used by {@link CMSpellPrereq}. */
	public static String	ONE_SPELL;
	/** Used by {@link CMSpellPrereq}. */
	public static String	MULTIPLE_SPELLS;
	/** Used by {@link CMSpellPrereq}. */
	public static String	WHOSE_NAME;
	/** Used by {@link CMSpellPrereq}. */
	public static String	OF_ANY_KIND;
	/** Used by {@link CMSpellPrereq}. */
	public static String	WHOSE_COLLEGE;
	/** Used by {@link CMSpellPrereq}. */
	public static String	COLLEGE_COUNT;

	/** Used by {@link CMAdvantagePrereq} and {@link CMSkillPrereq}. */
	public static String	LEVEL_PART;

	/**
	 * Used by {@link CMAdvantagePrereq}, {@link CMSkillPrereq}, {@link CMSpellPrereq},
	 * {@link CMContainedWeightPrereq} and {@link CMAttributePrereq}.
	 */
	public static String	HAS;
	/**
	 * Used by {@link CMAdvantagePrereq}, {@link CMSkillPrereq}, {@link CMSpellPrereq},
	 * {@link CMContainedWeightPrereq} and {@link CMAttributePrereq}.
	 */
	public static String	DOES_NOT_HAVE;

	static {
		TKMessages.initialize(Msgs.class);
	}
}
