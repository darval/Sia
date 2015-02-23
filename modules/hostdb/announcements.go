package hostdb

import (
	"strings"

	"github.com/NebulousLabs/Sia/consensus"
	"github.com/NebulousLabs/Sia/encoding"
	"github.com/NebulousLabs/Sia/modules"
)

// findHostAnnouncements returns a list of the host announcements found within
// a given block.
func findHostAnnouncements(height consensus.BlockHeight, b consensus.Block) (announcements []modules.HostEntry) {
	for _, t := range b.Transactions {
		for _, data := range t.ArbitraryData {
			// the HostAnnouncement must be prefaced by the standard host announcement string
			if !strings.HasPrefix(data, modules.PrefixHostAnnouncement) {
				continue
			}

			// decode the HostAnnouncement
			var ha modules.HostAnnouncement
			encAnnouncement := []byte(strings.TrimPrefix(data, modules.PrefixHostAnnouncement))
			err := encoding.Unmarshal(encAnnouncement, &ha)
			if err != nil {
				continue
			}

			// Add the announcement to the slice being returned.
			announcements = append(announcements, modules.HostEntry{
				IPAddress: ha.IPAddress,
			})
		}
	}

	return
}

// threadedInsertFromAnnouncement requests a host's hosting parameters, and inserts
// the resulting HostEntry into the database.
func (hdb *HostDB) threadedInsertFromAnnouncement(entry modules.HostEntry) {
	err := entry.IPAddress.RPC("HostSettings", nil, &entry.HostSettings)
	if err != nil {
		return
	}
	hdb.Insert(entry)
}