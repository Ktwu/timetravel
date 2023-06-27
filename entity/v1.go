package entity

type RecordV1 struct {
	ID   int               `json:"id"`
	Data map[string]string `json:"data"`
}

func (d *Record) IntoV1() RecordV1 {
	return RecordV1{
		ID:   d.ID,
		Data: d.Data,
	}
}
