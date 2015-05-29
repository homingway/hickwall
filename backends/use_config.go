package backends

import (
	"fmt"
	"github.com/oliveagle/hickwall/config"
	"github.com/oliveagle/hickwall/newcore"
)

func UseConfigCreateBackends(rconf *config.RuntimeConfig) ([]newcore.Publication, error) {
	var pubs []newcore.Publication

	if rconf == nil {
		return nil, fmt.Errorf("runtime config is nil")
	}

	// create file transport
	if rconf.Client.Transport_file != nil {
		b, err := NewFileBackend("file", rconf.Client.Transport_file)
		if err != nil {
			return nil, err
		}
		pubs = append(pubs, b)
	}

	if rconf.Client.Transport_dummy != nil {
		b := newcore.MustNewDummyBackend(
			rconf.Client.Transport_dummy.Name,
			rconf.Client.Transport_dummy.Jamming,
			rconf.Client.Transport_dummy.Printting)
		pubs = append(pubs, b)
	}

	return pubs[:], nil

}
