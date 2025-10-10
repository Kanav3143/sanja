package sanja

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNormalize(t *testing.T) {
	norm := NewNormalizer("MW")

	tests := []struct {
		name          string
		input         string
		expected      string
		expectError   bool
		errorContains string
	}{
		{
			name:     "local number with leading zero",
			input:    "0886392814",
			expected: "+265886392814",
		},
		{
			name:     "local number without leading zero",
			input:    "886392814",
			expected: "+265886392814",
		},
		{
			name:     "already normalized with plus",
			input:    "+265886392814",
			expected: "+265886392814",
		},
		{
			name:     "number with country code without plus",
			input:    "265886392814",
			expected: "+265886392814",
		},
		{
			name:     "number with multiple leading zeros",
			input:    "000886392814",
			expected: "+265886392814",
		},

		{
			name:     "number with spaces",
			input:    "088 639 2814",
			expected: "+265886392814",
		},
		{
			name:     "number with dashes",
			input:    "088-639-2814",
			expected: "+265886392814",
		},
		{
			name:     "number with parentheses",
			input:    "(088)6392814",
			expected: "+265886392814",
		},
		{
			name:     "number with mixed formatting",
			input:    "+265 (88) 639-2814",
			expected: "+265886392814",
		},

		{
			name:     "US number with country code",
			input:    "+12025550123",
			expected: "+12025550123",
		},
		{
			name:     "UK number with country code",
			input:    "+442079460000",
			expected: "+442079460000",
		},
		{
			name:     "South Africa number",
			input:    "+27821234567",
			expected: "+27821234567",
		},

		{
			name:          "empty string",
			input:         "",
			expectError:   true,
			errorContains: "invalid phone number",
		},
		{
			name:          "only special characters",
			input:         "+-() ",
			expectError:   true,
			errorContains: "invalid phone number",
		},
		{
			name:          "too short number",
			input:         "123",
			expectError:   true,
			errorContains: "invalid phone number",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := norm.Normalize(tc.input)

			if tc.expectError {
				assert.Error(t, err)
				if tc.errorContains != "" {
					assert.Contains(t, err.Error(), tc.errorContains)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, result)
			}
		})
	}
}

func TestNormalizeBulk(t *testing.T) {
	norm := NewNormalizer("MW")

	tests := []struct {
		name           string
		input          []string
		expectedCount  int
		errorCount     int
		expectedOutput []string
	}{
		{
			name: "mixed valid and invalid numbers",
			input: []string{
				"0886392814",
				"265886392814",
				"+265886392814",
				"00886392814",
				"",
				"123",
				"+447911123456",
			},
			expectedCount: 7,
			errorCount:    2,
			expectedOutput: []string{
				"+265886392814",
				"+265886392814",
				"+265886392814",
				"+265886392814",
				"",
				"",
				"+447911123456",
			},
		},
		{
			name: "all valid numbers",
			input: []string{
				"0886392814",
				"886392814",
				"+265886392814",
				"265886392814",
			},
			expectedCount: 4,
			errorCount:    0,
			expectedOutput: []string{
				"+265886392814",
				"+265886392814",
				"+265886392814",
				"+265886392814",
			},
		},
		{
			name: "all invalid numbers",
			input: []string{
				"",
				"abc",
				"123",
			},
			expectedCount: 3,
			errorCount:    3,
			expectedOutput: []string{
				"", "", "",
			},
		},
		{
			name:           "empty slice",
			input:          []string{},
			expectedCount:  0,
			errorCount:     0,
			expectedOutput: []string{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			normalizedNumbers, errors := norm.NormalizeBulk(tc.input)

			assert.Equal(t, len(tc.input), len(normalizedNumbers))
			assert.Equal(t, len(tc.input), len(errors))

			actualErrorCount := 0
			for _, err := range errors {
				if err != nil {
					actualErrorCount++
				}
			}

			assert.Equal(t, tc.errorCount, actualErrorCount)

			for i, expected := range tc.expectedOutput {
				if i < len(normalizedNumbers) {
					assert.Equal(t, expected, normalizedNumbers[i])
				}
			}
		})
	}
}

func TestNormalizeWithDifferentDefaultCountries(t *testing.T) {
	tests := []struct {
		name           string
		defaultCountry string
		input          string
		expected       string
	}{
		{
			name:           "US default country",
			defaultCountry: "US",
			input:          "5550123",
			expected:       "+15550123",
		},
		{
			name:           "UK default country",
			defaultCountry: "GB",
			input:          "2079460000",
			expected:       "+442079460000",
		},
		{
			name:           "ZA default country",
			defaultCountry: "ZA",
			input:          "821234567",
			expected:       "+27821234567",
		},
		{
			name:           "KE default country",
			defaultCountry: "KE",
			input:          "712345678",
			expected:       "+254712345678",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			norm := NewNormalizer(tc.defaultCountry)
			result, err := norm.Normalize(tc.input)
			assert.NoError(t, err)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestNormalizer_NoDefaultCountry(t *testing.T) {

	norm := NewNormalizer("INVALID")

	_, err := norm.Normalize("0886392814")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no default country set")
}

// Benchmark tests
func BenchmarkNormalize(b *testing.B) {
	norm := NewNormalizer("MW")
	phoneNumbers := []string{
		"0886392814",
		"265886392814",
		"+265886392814",
		"886392814",
	}

	for b.Loop() {
		for _, phone := range phoneNumbers {
			norm.Normalize(phone)
		}
	}
}

func BenchmarkNormalizeBulk(b *testing.B) {
	norm := NewNormalizer("MW")
	phoneNumbers := []string{
		"0886392814", "265886392814", "+265886392814", "886392814",
		"0999123456", "265999123456", "+265999123456", "999123456",
	}

	b.ResetTimer()
	for b.Loop() {
		norm.NormalizeBulk(phoneNumbers)
	}
}
