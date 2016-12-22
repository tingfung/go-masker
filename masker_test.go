package masker_test

import (
	"encoding/json"
	"testing"

	"github.com/syoya/go-masker"
)

func TestReplacementMapper(t *testing.T) {
	type Case struct {
		spec        string
		replacement map[string]string
		expected    interface{}
	}

	input, err := json.Marshal(map[string]interface{}{
		"foo": 1,
		"bar": 2,
		"baz": 3,
	})
	if err != nil {
		t.Fatalf("Fail to marshal input with error `%s`", err)
	}

	cases := []Case{
		{
			spec:        "none mapper",
			replacement: map[string]string{},
			expected: map[string]interface{}{
				"foo": 1,
				"bar": 2,
				"baz": 3,
			},
		},
		{
			spec: "non-existing fields",
			replacement: map[string]string{
				"qux": "**********",
			},
			expected: map[string]interface{}{
				"foo": 1,
				"bar": 2,
				"baz": 3,
			},
		},
		{
			spec: "single mapper",
			replacement: map[string]string{
				"bar": "*****",
			},
			expected: map[string]interface{}{
				"foo": 1,
				"bar": "*****",
				"baz": 3,
			},
		},
		{
			spec: "multi-mapper",
			replacement: map[string]string{
				"foo": "",
				"baz": "MASKED",
			},
			expected: map[string]interface{}{
				"foo": "",
				"bar": 2,
				"baz": "MASKED",
			},
		},
	}

	for _, c := range cases {
		m, err := masker.New(masker.Options{Replacement: c.replacement})
		if err != nil {
			t.Errorf("Fail to create masker with error `%s`", err)
			continue
		}
		a := m.Mask(input)

		e, err := json.Marshal(c.expected)
		if err != nil {
			t.Errorf("Fail to marshal expected with error `%s`", err)
			continue
		}

		if string(a) != string(e) {
			t.Errorf("%s, masked JSON is expected `%s`, but actual `%s`", c.spec, string(e), string(a))
		}
	}
}

func TestReplacementDeep(t *testing.T) {
	m, err := masker.New(masker.Options{
		Replacement: map[string]string{
			"password": "**********",
		},
	})
	if err != nil {
		t.Fatal("Fail to create masker with error `%s`", err)
	}

	type Case struct {
		spec     string
		input    interface{}
		expected interface{}
	}

	cases := []Case{
		{
			spec: "Mask() should replace a value in shallow object",
			input: map[string]interface{}{
				"qux":      "this shouldn't be masked",
				"password": "this should be masked",
			},
			expected: map[string]interface{}{
				"qux":      "this shouldn't be masked",
				"password": "**********",
			},
		},
		{
			spec: "Mask() should replace a value in deep object",
			input: map[string]interface{}{
				"foo": map[string]interface{}{
					"bar": map[string]interface{}{
						"baz": map[string]interface{}{
							"qux":      "this shouldn't be masked",
							"password": "this should be masked",
						},
					},
				},
			},
			expected: map[string]interface{}{
				"foo": map[string]interface{}{
					"bar": map[string]interface{}{
						"baz": map[string]interface{}{
							"qux":      "this shouldn't be masked",
							"password": "**********",
						},
					},
				},
			},
		},
		{
			spec: "Mask() should replace a value in shallow array",
			input: []interface{}{
				map[string]interface{}{
					"qux":      "this shouldn't be masked",
					"password": "this should be masked",
				},
			},
			expected: []interface{}{
				map[string]interface{}{
					"qux":      "this shouldn't be masked",
					"password": "**********",
				},
			},
		},
		{
			spec: "Mask() should replace a value in deep array",
			input: []interface{}{
				[]interface{}{
					[]interface{}{
						map[string]interface{}{
							"qux":      "this shouldn't be masked",
							"password": "this should be masked",
						},
					},
				},
			},
			expected: []interface{}{
				[]interface{}{
					[]interface{}{
						map[string]interface{}{
							"qux":      "this shouldn't be masked",
							"password": "**********",
						},
					},
				},
			},
		},
		{
			spec: "Mask() should replace values in complex data",
			input: map[string]interface{}{
				"qux0":     "this shouldn't be masked",
				"password": "this should be masked",
				"qux1":     "this shouldn't be masked",
				"foo": map[string]interface{}{
					"qux3":     "this shouldn't be masked",
					"password": "this should be masked",
					"qux4":     "this shouldn't be masked",
					"bar": []interface{}{
						"this shouldn't be masked",
						map[string]interface{}{
							"qux6": "this shouldn't be masked",
							"baz": []interface{}{
								"this shouldn't be masked",
								map[string]interface{}{
									"qux9":     "this shouldn't be masked",
									"password": "this should be masked",
									"qux10":    "this shouldn't be masked",
								},
								"this shouldn't be masked",
							},
							"qux7":     "this shouldn't be masked",
							"password": "this should be masked",
							"qux8":     "this shouldn't be masked",
						},
						"this shouldn't be masked",
					},
					"qux5": "this shouldn't be masked",
				},
				"qux2": "this shouldn't be masked",
			},
			expected: map[string]interface{}{
				"qux0":     "this shouldn't be masked",
				"password": "**********",
				"qux1":     "this shouldn't be masked",
				"foo": map[string]interface{}{
					"qux3":     "this shouldn't be masked",
					"password": "**********",
					"qux4":     "this shouldn't be masked",
					"bar": []interface{}{
						"this shouldn't be masked",
						map[string]interface{}{
							"qux6": "this shouldn't be masked",
							"baz": []interface{}{
								"this shouldn't be masked",
								map[string]interface{}{
									"qux9":     "this shouldn't be masked",
									"password": "**********",
									"qux10":    "this shouldn't be masked",
								},
								"this shouldn't be masked",
							},
							"qux7":     "this shouldn't be masked",
							"password": "**********",
							"qux8":     "this shouldn't be masked",
						},
						"this shouldn't be masked",
					},
					"qux5": "this shouldn't be masked",
				},
				"qux2": "this shouldn't be masked",
			},
		},
		{
			spec: "Mask() should replace a value in object containing nil",
			input: map[string]interface{}{
				"qux":      "this shouldn't be masked",
				"password": "this should be masked",
				"oops":     nil,
			},
			expected: map[string]interface{}{
				"qux":      "this shouldn't be masked",
				"password": "**********",
				"oops":     nil,
			},
		},
		{
			spec: "Mask() should replace a value in array containing nil",
			input: []interface{}{
				nil,
				map[string]interface{}{
					"qux":      "this shouldn't be masked",
					"password": "this should be masked",
				},
			},
			expected: []interface{}{
				nil,
				map[string]interface{}{
					"qux":      "this shouldn't be masked",
					"password": "**********",
				},
			},
		},
	}

	for _, c := range cases {
		i, err := json.Marshal(c.input)
		if err != nil {
			t.Errorf("Fail to marshal input with error `%s`", err)
			continue
		}

		var actual interface{}
		if err := json.Unmarshal(m.Mask(i), &actual); err != nil {
			t.Errorf("Fail to unmarshal actual with error `%s`", err)
			continue
		}
		a, err := json.Marshal(actual)
		if err != nil {
			t.Errorf("Fail to marshal actual with error `%s`", err)
			continue
		}

		e, err := json.Marshal(c.expected)
		if err != nil {
			t.Errorf("Fail to marshal expected with error `%s`", err)
			continue
		}

		if string(a) != string(e) {
			t.Errorf("%s, masked JSON is expected `%s`, but actual `%s`", c.spec, string(e), string(a))
		}
	}
}

func TestTrancation(t *testing.T) {
	type Case struct {
		spec       string
		truncation masker.Truncation
		input      interface{}
		expected   interface{}
	}

	cases := []Case{
		{
			spec:       "Mask() shouldn't truncate a plain string less than specified length",
			truncation: masker.Truncation{Length: 20, Omission: "..."},
			input:      "1234567890123456789",
			expected:   "1234567890123456789",
		},
		{
			spec:       "Mask() shouldn't truncate a plain string same as specified length",
			truncation: masker.Truncation{Length: 20, Omission: "..."},
			input:      "12345678901234567890",
			expected:   "12345678901234567890",
		},
		{
			spec:       "Mask() should truncate a plain string over specified length",
			truncation: masker.Truncation{Length: 20, Omission: "..."},
			input:      "123456789012345678901",
			expected:   "12345678901234567890...",
		},
		{
			spec:       "Mask() should truncate a value in shallow object",
			truncation: masker.Truncation{Length: 20, Omission: "..."},
			input: map[string]interface{}{
				"a": "abcde",
				"b": "this should be truncated",
			},
			expected: map[string]interface{}{
				"a": "abcde",
				"b": "this should be trunc...",
			},
		},
		{
			spec:       "Mask() should truncate a value in deep object",
			truncation: masker.Truncation{Length: 20, Omission: "..."},
			input: map[string]interface{}{
				"foo": map[string]interface{}{
					"bar": map[string]interface{}{
						"baz": map[string]interface{}{
							"a": "abcde",
							"b": "this should be truncated",
						},
					},
				},
			},
			expected: map[string]interface{}{
				"foo": map[string]interface{}{
					"bar": map[string]interface{}{
						"baz": map[string]interface{}{
							"a": "abcde",
							"b": "this should be trunc...",
						},
					},
				},
			},
		},
		{
			spec:       "Mask() should truncate a value in shallow array",
			truncation: masker.Truncation{Length: 20, Omission: "..."},
			input: []string{
				"abcde",
				"this should be truncated",
			},
			expected: []string{
				"abcde",
				"this should be trunc...",
			},
		},
		{
			spec:       "Mask() should truncate a value in deep array",
			truncation: masker.Truncation{Length: 20, Omission: "..."},
			input: []interface{}{
				[]interface{}{
					[]interface{}{
						[]string{
							"abcde",
							"this should be truncated",
						},
					},
				},
			},
			expected: []interface{}{
				[]interface{}{
					[]interface{}{
						[]string{
							"abcde",
							"this should be trunc...",
						},
					},
				},
			},
		},
		{
			spec:       "Mask() should truncate values in complex data",
			truncation: masker.Truncation{Length: 20, Omission: "..."},
			input: map[string]interface{}{
				"qux0": "abcde",
				"qux1": "this should be truncated",
				"foo": map[string]interface{}{
					"qux3": "abcde",
					"qux4": "this should be truncated",
					"bar": []interface{}{
						"abcde",
						map[string]interface{}{
							"qux6": "abcde",
							"baz": []interface{}{
								"abcde",
								map[string]interface{}{
									"qux9":  "abcde",
									"qux10": "this should be truncated",
								},
								"abcde",
							},
							"qux7": "abcde",
							"qux8": "this should be truncated",
						},
						"abcde",
					},
					"qux5": "abcde",
				},
				"qux2": "abcde",
			},
			expected: map[string]interface{}{
				"qux0": "abcde",
				"qux1": "this should be trunc...",
				"foo": map[string]interface{}{
					"qux3": "abcde",
					"qux4": "this should be trunc...",
					"bar": []interface{}{
						"abcde",
						map[string]interface{}{
							"qux6": "abcde",
							"baz": []interface{}{
								"abcde",
								map[string]interface{}{
									"qux9":  "abcde",
									"qux10": "this should be trunc...",
								},
								"abcde",
							},
							"qux7": "abcde",
							"qux8": "this should be trunc...",
						},
						"abcde",
					},
					"qux5": "abcde",
				},
				"qux2": "abcde",
			},
		},
	}

	for _, c := range cases {
		m, err := masker.New(masker.Options{Truncation: c.truncation})
		if err != nil {
			t.Errorf("Fail to create masker with error `%s`", err)
			continue
		}

		i, err := json.Marshal(c.input)
		if err != nil {
			t.Errorf("Fail to marshal input with error `%s`", err)
			continue
		}

		var actual interface{}
		if err := json.Unmarshal(m.Mask(i), &actual); err != nil {
			t.Errorf("Fail to unmarshal actual with error `%s`", err)
			continue
		}
		a, err := json.Marshal(actual)
		if err != nil {
			t.Errorf("Fail to marshal actual with error `%s`", err)
			continue
		}

		e, err := json.Marshal(c.expected)
		if err != nil {
			t.Errorf("Fail to marshal expected with error `%s`", err)
			continue
		}

		if string(a) != string(e) {
			t.Errorf("%s, masked JSON is expected `%s`, but actual `%s`", c.spec, e, a)
		}
	}
}
