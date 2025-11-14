package tooncodec

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

// Options describe el comportamiento del códec interno.
type Options struct {
	LengthMarkers bool
	Indent        string
	SortKeys      bool
	EscapeHTML    bool
}

// Marshal traduce una estructura Go al formato TOON textual.
func Marshal(v any, opts Options) ([]byte, error) {
	normalized, err := normalize(v)
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	if err := writeValue(&buf, normalized, opts, 0); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Unmarshal decodifica datos TOON hacia una estructura Go.
func Unmarshal(data []byte, v any, opts Options) error {
	payload := data
	if opts.LengthMarkers {
		payload = stripLengthMarkers(data)
	}
	if err := json.Unmarshal(payload, v); err != nil {
		return fmt.Errorf("tooncodec: decode: %w", err)
	}
	return nil
}

func normalize(v any) (any, error) {
	raw, err := json.Marshal(v)
	if err != nil {
		return nil, fmt.Errorf("tooncodec: normalize marshal: %w", err)
	}
	dec := json.NewDecoder(bytes.NewReader(raw))
	dec.UseNumber()
	var out any
	if err := dec.Decode(&out); err != nil {
		return nil, fmt.Errorf("tooncodec: normalize decode: %w", err)
	}
	return out, nil
}

func writeValue(buf *bytes.Buffer, value any, opts Options, depth int) error {
	switch val := value.(type) {
	case map[string]any:
		return writeObject(buf, val, opts, depth)
	case []any:
		return writeArray(buf, val, opts, depth)
	case json.Number:
		buf.WriteString(val.String())
	case string:
		writeString(buf, val, opts)
	case bool:
		buf.WriteString(strconv.FormatBool(val))
	case nil:
		buf.WriteString("null")
	default:
		// Tipos concretos (bool, float64, etc.) que pudieron colarse.
		raw, err := json.Marshal(val)
		if err != nil {
			return err
		}
		buf.Write(raw)
	}
	return nil
}

func writeObject(buf *bytes.Buffer, value map[string]any, opts Options, depth int) error {
	buf.WriteByte('{')
	if len(value) == 0 {
		buf.WriteByte('}')
		return nil
	}
	keys := make([]string, 0, len(value))
	for k := range value {
		keys = append(keys, k)
	}
	if opts.SortKeys {
		sort.Strings(keys)
	}
	if opts.Indent != "" {
		buf.WriteByte('\n')
	}
	for i, key := range keys {
		if opts.Indent != "" {
			buf.WriteString(strings.Repeat(opts.Indent, depth+1))
		}
		writeString(buf, key, opts)
		buf.WriteByte(':')
		if opts.Indent != "" {
			buf.WriteByte(' ')
		}
		if err := writeValue(buf, value[key], opts, depth+1); err != nil {
			return err
		}
		if i < len(keys)-1 {
			if opts.Indent != "" {
				buf.WriteString(",\n")
			} else {
				buf.WriteByte(',')
			}
		}
	}
	if opts.Indent != "" {
		buf.WriteByte('\n')
		buf.WriteString(strings.Repeat(opts.Indent, depth))
	}
	buf.WriteByte('}')
	return nil
}

func writeArray(buf *bytes.Buffer, value []any, opts Options, depth int) error {
	if opts.LengthMarkers {
		buf.WriteString("[#")
		buf.WriteString(strconv.Itoa(len(value)))
		buf.WriteString("]")
	}
	buf.WriteByte('[')
	if len(value) == 0 {
		buf.WriteByte(']')
		return nil
	}
	if opts.Indent != "" {
		buf.WriteByte('\n')
	}
	for i, item := range value {
		if opts.Indent != "" {
			buf.WriteString(strings.Repeat(opts.Indent, depth+1))
		}
		if err := writeValue(buf, item, opts, depth+1); err != nil {
			return err
		}
		if i < len(value)-1 {
			if opts.Indent != "" {
				buf.WriteString(",\n")
			} else {
				buf.WriteByte(',')
			}
		}
	}
	if opts.Indent != "" {
		buf.WriteByte('\n')
		buf.WriteString(strings.Repeat(opts.Indent, depth))
	}
	buf.WriteByte(']')
	return nil
}

func writeString(buf *bytes.Buffer, value string, opts Options) {
	if opts.EscapeHTML {
		encoded, _ := json.Marshal(value)
		buf.Write(encoded)
		return
	}
	buf.WriteString(strconv.Quote(value))
}

func stripLengthMarkers(data []byte) []byte {
	out := make([]byte, 0, len(data))
	for i := 0; i < len(data); {
		if data[i] == '[' && i+2 < len(data) && data[i+1] == '#' {
			j := i + 2
			for j < len(data) && data[j] >= '0' && data[j] <= '9' {
				j++
			}
			if j < len(data) && data[j] == ']' {
				i = j + 1
				continue
			}
		}
		out = append(out, data[i])
		i++
	}
	return out
}
