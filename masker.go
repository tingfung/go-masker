package masker

import (
	"encoding/json"
	"strings"
)

type Replacement map[string]string

func (r Replacement) toRawMessageMap() (map[string]*json.RawMessage, error) {
	m := make(map[string]*json.RawMessage)
	for k, v := range r {
		b, err := json.Marshal(v)
		if err != nil {
			return m, err
		}
		raw := json.RawMessage(b)
		m[strings.ToLower(k)] = &raw
	}
	return m, nil
}

type Truncation struct {
	Length   int
	Omission string
}

func (t Truncation) valid() bool {
	return t.Length > 0
}

type Options struct {
	Replacement Replacement
	Truncation  Truncation
}

func New(o Options) (Masker, error) {
	m := Masker{}
	r, err := o.Replacement.toRawMessageMap()
	if err != nil {
		return m, nil
	}
	m.replacement = r
	m.truncation = o.Truncation
	return m, err
}

type Masker struct {
	replacement map[string]*json.RawMessage
	truncation  Truncation
}

func (m Masker) Mask(i []byte) []byte {
	if o, err := m.maskAsObject(i); err == nil {
		return o
	}
	if o, err := m.maskAsArray(i); err == nil {
		return o
	}
	if m.truncation.valid() {
		if o, err := m.maskAsString(i); err == nil {
			return o
		}
	}
	return i
}

func (m Masker) maskAsObject(i []byte) ([]byte, error) {
	var o map[string]*json.RawMessage
	if err := json.Unmarshal(i, &o); err != nil {
		return i, err
	}

	for k, v := range o {
		lk := strings.ToLower(k)
		if r, ok := m.replacement[lk]; ok {
			o[k] = r
			continue
		}
		if v == nil {
			continue
		}
		r := json.RawMessage(m.Mask(*v))
		o[k] = &r
	}

	return json.Marshal(o)
}

func (m Masker) maskAsArray(i []byte) ([]byte, error) {
	var a []*json.RawMessage
	if err := json.Unmarshal(i, &a); err != nil {
		return i, err
	}

	for i, v := range a {
		if v == nil {
			continue
		}
		r := json.RawMessage(m.Mask(*v))
		a[i] = &r
	}
	return json.Marshal(a)
}

func (m Masker) maskAsString(i []byte) ([]byte, error) {
	var s string
	if err := json.Unmarshal(i, &s); err != nil {
		return i, err
	}
	if len(s) <= m.truncation.Length {
		return i, nil
	}
	return json.Marshal(s[:m.truncation.Length] + m.truncation.Omission)
}
