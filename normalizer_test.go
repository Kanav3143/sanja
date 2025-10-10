package sanja

import (
	"fmt"
	"testing"

	assert2 "github.com/stretchr/testify/assert"
)

func TestNormalize(t *testing.T) {
	assert := assert2.New(t)

	norm := NewNormalizer("MW")

	expectedNumber := "+265886392814"

	normalizedPhone, err := norm.Normalize("0886392814")

	assert.Equal(expectedNumber, normalizedPhone)
	assert.ErrorIs(err, nil)
}

func TestNormalizeBulk(t *testing.T) {
	// assert := assert2.New(t)

	// TODO: add more test cases

	norm := NewNormalizer("MW")

	phoneNumbers := []string{
		"0886392814",
		"265886392814",
		"+265886392814",
		"00886392814",
	}

	normalizedNumbers, errors := norm.NormalizeBulk(phoneNumbers)

	if len(errors) > 0 {
		for _, err := range errors {
			if err != nil {
				fmt.Printf("failed to normalize phone number: Error, %v\n", err)
			}
		}
	}

	for _, phone := range normalizedNumbers {
		fmt.Printf("Phone: %s\n", phone)
	}
}
