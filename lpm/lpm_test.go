/*
 * Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
 * or more contributor license agreements. Licensed under the Apache License 2.0.
 * See the file "LICENSE" for details.
 */

package lpm

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestGetRightmostSetBit(t *testing.T) {
	tests := map[string]struct {
		input    uint64
		expected uint64
	}{
		"1":   {input: 0b1, expected: 0b1},
		"2":   {input: 0b10, expected: 0b10},
		"3":   {input: 0b11, expected: 0b1},
		"160": {input: 0b10100000, expected: 0b100000},
	}

	for name, test := range tests {
		name := name
		test := test
		t.Run(name, func(t *testing.T) {
			output := getRightmostSetBit(test.input)
			if output != test.expected {
				t.Fatalf("Expected %d (0b%b) but got %d (0b%b)",
					test.expected, test.expected, output, output)
			}
		})
	}
}

func TestCalculatePrefixList(t *testing.T) {
	tests := map[string]struct {
		start  uint64
		end    uint64
		err    bool
		expect []Prefix
	}{
		"4k to 0": {start: 4096, end: 0, err: true},
		"10 to 22": {start: 0b1010, end: 0b10110,
			expect: []Prefix{{0b1010, 63}, {0b1100, 62}, {0b10000, 62},
				{0b10100, 63}}},
		"4k to 16k": {start: 4096, end: 16384,
			expect: []Prefix{{0x1000, 52}, {0x2000, 51}}},
		"0x55ff3f68a000 to 0x55ff3f740000": {start: 0x55ff3f68a000, end: 0x55ff3f740000,
			expect: []Prefix{{0x55ff3f68a000, 51}, {0x55ff3f68c000, 50},
				{0x55ff3f690000, 48}, {0x55ff3f6a0000, 47},
				{0x55ff3f6c0000, 46}, {0x55ff3f700000, 46}}},
		"0x7f5b6ef4f000 to 0x7f5b6ef5d000": {start: 0x7f5b6ef4f000, end: 0x7f5b6ef5d000,
			expect: []Prefix{{0x7f5b6ef4f000, 52}, {0x7f5b6ef50000, 49},
				{0x7f5b6ef58000, 50}, {0x7f5b6ef5c000, 52}}},
	}

	for name, test := range tests {
		name := name
		test := test
		t.Run(name, func(t *testing.T) {
			prefixes, err := CalculatePrefixList(test.start, test.end)
			if err != nil {
				if test.err {
					// We received and expected an error. So we can return here.
					return
				}
				t.Fatalf("Unexpected error: %v", err)
			}
			if test.err {
				t.Fatalf("Expected an error but got none")
			}
			if diff := cmp.Diff(test.expect, prefixes); diff != "" {
				t.Fatalf("CalculatePrefixList() mismatching prefixes (-want +got):\n%s", diff)
			}
		})
	}
}
