package sanja

import (
	"fmt"
	"strings"
)

// Country represents a country with its dialing codes
type Country struct {
	Name        string
	A2          string
	A3          string
	NumCode     int
	DialingCode string
}

// GetCountryByA2 returns a countyr by its A2 code
func (n *Normalizer) GetCountryByA2(a2 string) *Country {
	for i := range n.countries {
		if n.countries[i].A2 == a2 {
			return &n.countries[i]
		}
	}

	return nil
}

// GetCountryByCode returns a country by its dialing code
func (n *Normalizer) GetCountryByCode(code string) *Country {
	return n.codeMap[code]
}

// ValidatePhoneNumber validates if a phone number matches a specific country
func (n *Normalizer) ValidatePhoneNumber(phone, countryA2 string) error {
	country := n.GetCountryByA2(countryA2)
	if country == nil {
		return fmt.Errorf("%w: country: %s ", ErrUnknownCountry, countryA2)
	}

	cleanPhone := cleanPhone(phone)
	expectedCodes := splitDialingCodes(country.DialingCode)

	for _, code := range expectedCodes {
		if strings.HasPrefix(cleanPhone, fmt.Sprintf("+%s", code)) || strings.HasPrefix(cleanPhone, code) {
			return nil
		}
	}

	return fmt.Errorf("%w: country name: %s", ErrPhoneNumberAndCountryCodeMismatch, country.Name)
}
