package masker

import "encoding/json"

type Masker struct {
	mapper map[string]*json.RawMessage
}

func New(m map[string]string) (Masker, error) {
	mapper := make(map[string]*json.RawMessage)
	for k, v := range m {
		b, err := json.Marshal(v)
		if err != nil {
			return Masker{}, err
		}
		r := json.RawMessage(b)
		mapper[k] = &r
	}
	return Masker{mapper}, nil
}

func (m Masker) Mask(i []byte) []byte {
	{
		o, err := m.maskAsObject(i)
		if err == nil {
			return o
		}
	}
	{
		o, err := m.maskAsArray(i)
		if err == nil {
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
		if r, ok := m.mapper[k]; ok {
			o[k] = r
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
		r := json.RawMessage(m.Mask(*v))
		a[i] = &r
	}
	return json.Marshal(a)
}
