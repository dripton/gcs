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

package com.trollworks.toolkit.utility.units;

/** Holds a value and {@link TKWeightUnits} pair. */
public class TKWeightValue extends TKUnitsValue<TKWeightUnits> {
	/**
	 * Creates a new {@link TKUnitsValue}.
	 * 
	 * @param value The value to use.
	 * @param units The {@link TKUnits} to use.
	 */
	public TKWeightValue(double value, TKWeightUnits units) {
		super(value, units);
	}

	/**
	 * Creates a new {@link TKUnitsValue}.
	 * 
	 * @param other A {@link TKUnitsValue} to clone.
	 */
	public TKWeightValue(TKWeightValue other) {
		super(other);
	}

	@Override public TKWeightUnits getDefaultUnits() {
		return TKWeightUnits.POUNDS;
	}
}
