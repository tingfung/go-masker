package masker_test

import (
	"encoding/json"
	"testing"

	"github.com/syoya/go-masker"
)

func TestMapper(t *testing.T) {
	type Case struct {
		description string
		inputMapper map[string]string
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
			description: "none mapper",
			inputMapper: map[string]string{},
			expected: map[string]interface{}{
				"foo": 1,
				"bar": 2,
				"baz": 3,
			},
		},
		{
			description: "none mapper",
			inputMapper: map[string]string{
				"qux": "**********",
			},
			expected: map[string]interface{}{
				"foo": 1,
				"bar": 2,
				"baz": 3,
			},
		},
		{
			description: "single mapper",
			inputMapper: map[string]string{
				"bar": "*****",
			},
			expected: map[string]interface{}{
				"foo": 1,
				"bar": "*****",
				"baz": 3,
			},
		},
		{
			description: "multi-mapper",
			inputMapper: map[string]string{
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
		m, err := masker.New(c.inputMapper)
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
			t.Errorf("Test with %s, masked JSON is expected `%s`, but actual `%s`", c.description, string(e), string(a))
		}
	}
}

func TestMaskWithDepth(t *testing.T) {
	m, err := masker.New(map[string]string{
		"password": "**********",
	})
	if err != nil {
		t.Fatal("Fail to create masker with error `%s`", err)
	}

	type Case struct {
		description string
		input       interface{}
		expected    interface{}
	}

	cases := []Case{
		{
			description: "shallow object",
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
			description: "deep object",
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
			description: "shallow array",
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
			description: "deep array",
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
			description: "complex data",
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
			t.Errorf("Test with %s, masked JSON is expected `%s`, but actual `%s`", c.description, string(e), string(a))
		}
	}
}
