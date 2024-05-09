package gopeg_test

import (
	"gopeg"
	"testing"
)

func TestAdd(t *testing.T) {
	p := gopeg.P(2)
	m := p.Match("test", 0)
	if m == nil {
		t.Errorf("Match expected for P(2)")
	} else {
		val := m.GetValue()
		if val != "te" {
			t.Errorf("For P(2).match('test', 0) 'te' was expected. Got '%s'", val)
		}
	}
	//expected := 5
	//if result != expected {
	//	t.Errorf("Add(2, 3) returned %d, expected %d", result, expected)
	//}
}
