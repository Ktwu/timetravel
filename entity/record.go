package entity

type Record struct {
	ID   int               `json:"id"`
	Data map[string]string `json:"data"`

	// v2 fields
	Version int				`json:"version,omitempty"`
}

func (d *Record) Sanitize(apiVersion int) {
	if apiVersion < 2 {
		d.Version = 0
	}
}

func (d *Record) Copy() Record {
	values := d.Data

	newMap := map[string]string{}
	for key, value := range values {
		newMap[key] = value
	}

	return Record{
		ID:   d.ID,
		Data: newMap,
		Version: d.Version,
	}
}
