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

package com.trollworks.toolkit.io.cmdline;

import com.trollworks.toolkit.io.TKMessages;

/** All localized strings in this package should be accessed via this class. */
class Msgs {
	/** Used by {@link TKCmdLine}. */
	public static String	AVAILABLE_OPTIONS;
	/** Used by {@link TKCmdLine}. */
	public static String	UNEXPECTED_OPTION;
	/** Used by {@link TKCmdLine}. */
	public static String	UNEXPECTED_OPTION_ARGUMENT;
	/** Used by {@link TKCmdLine}. */
	public static String	MISSING_OPTION_ARGUMENT;
	/** Used by {@link TKCmdLine}. */
	public static String	HELP_DESCRIPTION;
	/** Used by {@link TKCmdLine}. */
	public static String	VERSION_DESCRIPTION;
	/** Used by {@link TKCmdLine}. */
	public static String	REFERENCE_DESCRIPTION;

	static {
		TKMessages.initialize(Msgs.class);
	}
}
