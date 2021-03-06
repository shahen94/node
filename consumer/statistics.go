/*
 * Copyright (C) 2017 The "MysteriumNetwork/node" Authors.
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

package consumer

// calcStatDiff takes in the old and the new values of statistics, returns the calculated delta
func calcStatDiff(old, new uint64) (res uint64) {
	if old > new {
		return new
	}
	return new - old
}

// DiffWithNew calculates the difference in bytes between the old stats and new
func (ss SessionStatistics) DiffWithNew(new SessionStatistics) SessionStatistics {
	return SessionStatistics{
		BytesSent:     calcStatDiff(ss.BytesSent, new.BytesSent),
		BytesReceived: calcStatDiff(ss.BytesReceived, new.BytesReceived),
	}
}

// AddUpStatistics adds up the given statistics with the diff and returns new stats
func AddUpStatistics(stats, diff SessionStatistics) SessionStatistics {
	return SessionStatistics{
		BytesReceived: stats.BytesReceived + diff.BytesReceived,
		BytesSent:     stats.BytesSent + diff.BytesSent,
	}
}

// SessionStatistics represents statistics, generated by bytescount middleware
type SessionStatistics struct {
	BytesSent, BytesReceived uint64
}
