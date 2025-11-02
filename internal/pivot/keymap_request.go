package pivot

// KeymapRequest merepresentasikan permintaan padanan kunci yang harus dicatat
// di database pivot ketika target menghasilkan nilai kunci baru. Struktur ini
// membantu checker dan executor saling bertukar konteks mengenai nama map,
// sumber data, serta kolom target yang nantinya akan diisi oleh hasil RETURNING.
type KeymapRequest struct {
	MapName   string `json:"map_name"`
	SrcKey    string `json:"src_key"`
	SrcColumn string `json:"src_column,omitempty"`
	TgtColumn string `json:"tgt_column,omitempty"`
	SrcTable  string `json:"src_table,omitempty"`
	TgtTable  string `json:"tgt_table,omitempty"`
}
