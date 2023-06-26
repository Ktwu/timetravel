package entity

import "github.com/google/go-cmp/cmp"

type Record struct {
	ID   int               `json:"id"`
	Data map[string]string `json:"data"`

	// v2 fields
	Version int `json:"version,omitempty"`
}

func (d *Record) Sanitize(apiVersion int) {
	if apiVersion < 2 {
		d.Version = 0
	}
}

func (d *Record) ApplyUpdate(updates map[string]*string) bool {
	var didChange bool
	for key, value := range updates {
		if value == nil {
			// During deletion, a real change is made if the key previously existed
			if !didChange {
				_, didChange = d.Data[key]
			}
			delete(d.Data, key)
		} else {
			// During addition/update, a real change is made if the previous value didn't exist
			if !didChange {
				prevValue, _ := d.Data[key]
				didChange = !cmp.Equal(prevValue, value)
			}
			d.Data[key] = *value
		}
	}
	return didChange
}

// Returns a map that negates any changes the updates argument would
// make to the record once applied.
func (d Record) InverseUpdate(
	updates map[string]*string,
) map[string]*string {
	updateReversal := make(map[string]*string)

	for key, value := range updates {
		prevValue, exists := d.Data[key]
		if exists && !cmp.Equal(prevValue, value) {
			updateReversal[key] = &prevValue
		}
		if !exists && value != nil {
			updateReversal[key] = nil
		}
	}

	return updateReversal
}

func (d *Record) Copy() Record {
	values := d.Data

	newMap := map[string]string{}
	for key, value := range values {
		newMap[key] = value
	}

	return Record{
		ID:      d.ID,
		Data:    newMap,
		Version: d.Version,
	}
}
