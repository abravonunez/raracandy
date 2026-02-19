package money

import (
	"testing"
)

func TestBCDConversion(t *testing.T) {
	tests := []struct {
		decimal byte
		bcd     byte
	}{
		{0, 0x00},
		{23, 0x23},
		{99, 0x99},
		{50, 0x50},
		{12, 0x12},
		{5, 0x05},
	}

	for _, tt := range tests {
		// Test decimal to BCD
		bcd := decimalToBCD(tt.decimal)
		if bcd != tt.bcd {
			t.Errorf("decimalToBCD(%d) = 0x%02X, want 0x%02X", tt.decimal, bcd, tt.bcd)
		}

		// Test BCD to decimal
		dec := bcdToDecimal(tt.bcd)
		if dec != tt.decimal {
			t.Errorf("bcdToDecimal(0x%02X) = %d, want %d", tt.bcd, dec, tt.decimal)
		}
	}
}

func TestBCDRoundtrip(t *testing.T) {
	for i := byte(0); i <= 99; i++ {
		bcd := decimalToBCD(i)
		dec := bcdToDecimal(bcd)
		if dec != i {
			t.Errorf("Roundtrip failed for %d: got %d", i, dec)
		}
	}
}

func TestFormatMoney(t *testing.T) {
	tests := []struct {
		amount   uint32
		expected string
	}{
		{0, "¥0"},
		{123, "¥123"},
		{1234, "¥1,234"},
		{12345, "¥12,345"},
		{123456, "¥123,456"},
		{999999, "¥999,999"},
	}

	for _, tt := range tests {
		result := FormatMoney(tt.amount)
		if result != tt.expected {
			t.Errorf("FormatMoney(%d) = %s, want %s", tt.amount, result, tt.expected)
		}
	}
}
