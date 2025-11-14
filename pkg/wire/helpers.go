package wire

import "fmt"

// JSONToTOON convierte un payload JSON a TOON usando el Transformer provisto.
func JSONToTOON[T any](t *Transformer, jsonData []byte) ([]byte, error) {
	trans := ensureTransformer(t)
	var payload T
	if err := trans.Unmarshal(jsonData, FormatJSON, &payload); err != nil {
		return nil, fmt.Errorf("json->toon: %w", err)
	}
	data, err := trans.Marshal(payload, FormatTOON)
	if err != nil {
		return nil, fmt.Errorf("json->toon marshal: %w", err)
	}
	return data, nil
}

// TOONToJSON convierte un payload TOON a JSON usando el Transformer provisto.
func TOONToJSON[T any](t *Transformer, toonData []byte) ([]byte, error) {
	trans := ensureTransformer(t)
	var payload T
	if err := trans.Unmarshal(toonData, FormatTOON, &payload); err != nil {
		return nil, fmt.Errorf("toon->json: %w", err)
	}
	data, err := trans.Marshal(payload, FormatJSON)
	if err != nil {
		return nil, fmt.Errorf("toon->json marshal: %w", err)
	}
	return data, nil
}

func ensureTransformer(t *Transformer) *Transformer {
	if t != nil {
		return t
	}
	return NewTransformer()
}
