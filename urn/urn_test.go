package urn

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	tests := []struct {
		desc        string
		namespace   Match
		partition   Match
		service     Match
		region      Match
		identifier  Match
		resource    Match
		expected    *URN
		expectedErr bool
	}{
		{
			desc:       "returns a resource URN",
			namespace:  "urn",
			partition:  "cloud",
			service:    "machine",
			region:     "eu-central-1",
			identifier: "1234567890",
			resource:   "ulysses",
			expected: &URN{
				Namespace:  "urn",
				Partition:  "cloud",
				Service:    "machine",
				Region:     "eu-central-1",
				Identifier: "1234567890",
				Resource:   "ulysses",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			urn, err := New(tc.namespace, tc.partition, tc.service, tc.region, tc.identifier, tc.resource)

			if tc.expectedErr {
				require.Error(t, err)
			}
			assert.Equal(t, tc.expected, urn)
		})
	}
}

func BenchmarkNew(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		New("urn", "cloud", "machine", "eu-central-1", "1234567890", "ulysses")
	}
}

func BenchmarkParse(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		Parse("urn:cloud:machine:eu-central-1:1234567890:ulysses")
	}
}

func TestMatch(t *testing.T) {
	tests := []struct {
		desc     string
		urn      *URN
		other    *URN
		expected bool
	}{
		{
			desc: "returns true when the URNs are equal",
			urn: &URN{
				Namespace:  "urn",
				Partition:  "cloud",
				Service:    "machine",
				Region:     "eu-central-1",
				Identifier: "1234567890",
				Resource:   "ulysses",
			},
			other: &URN{
				Namespace:  "urn",
				Partition:  "cloud",
				Service:    "machine",
				Region:     "eu-central-1",
				Identifier: "1234567890",
				Resource:   "ulysses",
			},
			expected: true,
		},
		{
			desc: "returns true when the URNs are equal and the other has a wildcard",
			urn: &URN{
				Namespace:  "urn",
				Partition:  "cloud",
				Service:    "machine",
				Region:     "eu-central-1",
				Identifier: "1234567890",
				Resource:   "ulysses",
			},
			other: &URN{
				Namespace:  "urn",
				Partition:  "cloud",
				Service:    "machine",
				Region:     "eu-central-1",
				Identifier: "1234567890",
				Resource:   "*",
			},
			expected: true,
		},
		{
			desc: "returns true when the URNs are equal and the other has a wildcard",
			urn: &URN{
				Namespace:  "urn",
				Partition:  "cloud",
				Service:    "machine",
				Region:     "eu-central-1",
				Identifier: "1234567890",
				Resource:   "ulysses",
			},
			other: &URN{
				Namespace:  "urn",
				Partition:  "cloud",
				Service:    "machine",
				Region:     "eu-central-1",
				Identifier: "*",
				Resource:   "*",
			},
			expected: true,
		},
		{
			desc: "returns true when the URNs are equal and the other has a wildcard",
			urn: &URN{
				Namespace:  "urn",
				Partition:  "cloud",
				Service:    "machine",
				Region:     "eu-central-1",
				Identifier: "1234567890",
				Resource:   "*",
			},
			other: &URN{
				Namespace:  "urn",
				Partition:  "cloud",
				Service:    "machine",
				Region:     "eu-central-1",
				Identifier: "1234567890",
				Resource:   "ulysses",
			},
			expected: false,
		},
		{
			desc: "returns true when the URNs are equal and the other has a wildcard",
			urn: &URN{
				Namespace:  "urn",
				Partition:  "cloud",
				Service:    "machine",
				Region:     "",
				Identifier: "1234567890",
				Resource:   "*",
			},
			other: &URN{
				Namespace:  "urn",
				Partition:  "cloud",
				Service:    "machine",
				Region:     "",
				Identifier: "1234567890",
				Resource:   "ulysses",
			},
			expected: false,
		},
		{
			desc: "returns true when the URNs are equal and the other has a wildcard",
			urn: &URN{
				Namespace:  "urn",
				Partition:  "cloud",
				Service:    "machine",
				Region:     "",
				Identifier: "1234567890",
				Resource:   "ulysses",
			},
			other: &URN{
				Namespace:  "urn",
				Partition:  "cloud",
				Service:    "machine",
				Region:     "",
				Identifier: "1234567890",
				Resource:   "*",
			},
			expected: true,
		},
		{
			desc: "returns false when the URNs are equal and the other has a wildcard",
			urn: &URN{
				Namespace:  "urn",
				Partition:  "cloud",
				Service:    "machine",
				Region:     "",
				Identifier: "12345678910",
				Resource:   "ulysses",
			},
			other: &URN{
				Namespace:  "urn",
				Partition:  "cloud",
				Service:    "machine",
				Region:     "",
				Identifier: "1234567890",
				Resource:   "ulysses",
			},
			expected: false,
		},
		{
			desc: "returns true when the URNs are equal and the other has a wildcard",
			urn: &URN{
				Namespace:  "urn",
				Partition:  "cloud",
				Service:    "machine",
				Region:     "de-central-1",
				Identifier: "1234567890",
				Resource:   "ulysses",
			},
			other: &URN{
				Namespace:  Wildcard,
				Partition:  Wildcard,
				Service:    Wildcard,
				Region:     Wildcard,
				Identifier: Wildcard,
				Resource:   Wildcard,
			},
			expected: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			assert.Equal(t, tc.expected, tc.urn.Match(tc.other))
		})
	}
}

func TestParse(t *testing.T) {
	tests := []struct {
		desc        string
		urnStr      string
		expected    *URN
		expectedErr bool
	}{
		{
			desc:        "returns an ErrorInvalid when identifier is longer than 256 chars",
			urnStr:      "urn:collection:::aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaas:123456",
			expectedErr: true,
		},
		{
			desc:   "returns a wildcard URN",
			urnStr: "*",
			expected: &URN{
				Namespace:  Wildcard,
				Partition:  Wildcard,
				Service:    Wildcard,
				Region:     Wildcard,
				Identifier: Wildcard,
				Resource:   Wildcard,
			},
			expectedErr: false,
		},
		{
			desc:   "returns a wildcard URN",
			urnStr: "urn:*:*::1234567890:*",
			expected: &URN{
				Namespace:  "urn",
				Partition:  Wildcard,
				Service:    Wildcard,
				Region:     Wildcard,
				Identifier: "1234567890",
				Resource:   Wildcard,
			},
			expectedErr: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			urn, err := Parse(tc.urnStr)

			if tc.expectedErr {
				require.Error(t, err)
			}

			assert.Equal(t, tc.expected, urn)
		})
	}
}
