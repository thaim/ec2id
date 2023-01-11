package main

import (
	"bytes"
	"testing"
)

func TestMainPrintIds(t *testing.T) {
	cases := []struct {
		name    string
		ids     []string
		allFlag bool
		expect  string
	}{
		{
			name:    "print nothing with empty id",
			ids:     nil,
			allFlag: false,
			expect:  "",
		},
		{
			name:    "print single id with single id",
			ids:     []string{"i-012345"},
			allFlag: false,
			expect:  "i-012345\n",
		},
		{
			name:    "print single id with multiple ids and disabled all flag",
			ids:     []string{"i-012345", "i-6789ab", "i-cdef01"},
			allFlag: false,
			expect:  "i-012345\n",
		},
		{
			name:    "print multiple ids with multiple ids and enabled all flag",
			ids:     []string{"i-012345", "i-6789ab", "i-cdef01"},
			allFlag: true,
			expect:  "i-012345\ni-6789ab\ni-cdef01\n",
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer

			printIds(&buf, tt.ids, tt.allFlag)

			if bufString := buf.String(); bufString != tt.expect {
				t.Errorf("expect %s, got id: %s", tt.expect, buf.String())
			}
		})
	}
}
