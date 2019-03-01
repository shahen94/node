/*
 * Copyright (C) 2019 The "MysteriumNetwork/node" Authors.
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package tequilapi

import "strings"

// WhitelistingCorsPolicy allows customizing CORS (Cross-Origin Resource Sharing) behaviour - whitelisting only specific domains
type WhitelistingCorsPolicy struct {
	DefaultTrustedOrigin  string
	AllowedOriginSuffixes []string
}

// AllowedOrigin returns the same request origin if it is allowed, or default origin if it is not allowed
func (policy WhitelistingCorsPolicy) AllowedOrigin(requestOrigin string) string {
	if policy.isOriginAllowed(requestOrigin) {
		return requestOrigin
	}

	return policy.DefaultTrustedOrigin
}

func (policy WhitelistingCorsPolicy) isOriginAllowed(origin string) bool {
	for _, allowedSuffix := range policy.AllowedOriginSuffixes {
		if strings.HasSuffix(origin, allowedSuffix) {
			return true
		}
	}
	return false
}