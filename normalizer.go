package sanja

import (
	"fmt"
	"strings"
)

const (
	MinimumLocalDigitsLength int = 7
	MaximumLocalDigitsLength int = 12
)

// Normalizer hanldes phone numbder normalization
type Normalizer struct {
	countries      []Country
	defaultCountry *Country
	codeMap        map[string]*Country
}

// NewNormalizer creatres a new phone number normalizer
func NewNormalizer(defaultCountryA2 string) *Normalizer {
	var n = Normalizer{
		countries: getCountries(),
		codeMap:   make(map[string]*Country),
	}

	for i := range n.countries {
		country := &n.countries[i]

		codes := splitDialingCodes(country.DialingCode)

		for _, code := range codes {
			n.codeMap[code] = country
		}
	}

	// set default country
	n.defaultCountry = n.GetCountryByA2(defaultCountryA2)

	return &n
}

// Normalize adds the country code prefix to a phone number if missing
func (n *Normalizer) Normalize(phone string) (string, error) {
	if n.defaultCountry == nil {
		return phone, ErrDefaultCountryNotFound
	}

	phone = cleanPhone(phone)

	if phone == "" || len(phone) < MinimumLocalDigitsLength {
		return "", ErrInvalidPhoneNumber
	}

	if strings.HasPrefix(phone, "+") {
		return phone, nil
	}

	// Check it already starts with a country code without +
	if n.hasCountryCode(phone) {
		return fmt.Sprintf("+%s", phone), nil
	}

	phone = strings.TrimLeft(phone, "0") // removes leading zeros in local numbers

	defaultCode := splitDialingCodes(n.defaultCountry.DialingCode)[0]

	return fmt.Sprintf("+%s%s", defaultCode, phone), nil
}

// NormalizeBulk normalizes multiple phone numbers
func (n *Normalizer) NormalizeBulk(phones []string) ([]string, []error) {
	var normalized = make([]string, len(phones))
	var errors = make([]error, len(phones))

	for i, phone := range phones {
		normalized[i], errors[i] = n.Normalize(phone)
	}

	return normalized, errors
}

// hasCountryCode checks if a phone number already starts with a valid country code
func (n *Normalizer) hasCountryCode(phone string) bool {
	cleanPhone := strings.TrimLeft(phone, "0")
	for i := 1; i <= 4 && i <= len(phone); i++ {
		prefix := cleanPhone[:i]

		if country, exists := n.codeMap[prefix]; exists {
			return country.A2 == n.defaultCountry.A2
		}
	}

	return false
}
