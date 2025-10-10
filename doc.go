// Package sanja provides phone number normalization and formatting.
//
// Overview:
// Sanja helps normalize phone numbers to international format with
// country code handling. It supports 250+ countries and bulk processing.
//
// Example:
//
//	norm := sanja.NewNormalizer("US")
//	normalized, err := norm.Normalize("555-123-4567")
//	// Result: "+15551234567"
//
// Features:
//   - Local to international number conversion
//   - Bulk phone number processing
//   - Country code detection
//   - Comprehensive country data
package sanja
