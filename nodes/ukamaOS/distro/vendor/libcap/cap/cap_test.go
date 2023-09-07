package cap

import (
	"fmt"
	"testing"
)

func TestAllMask(t *testing.T) {
	oldMask := maxValues
	oldWords := words
	defer func() {
		maxValues = oldMask
		words = oldWords
	}()

	maxValues = 35
	words = 3

	vs := []struct {
		val   Value
		index uint
		bit   uint32
		mask  uint32
	}{
		{val: CHOWN, index: 0, bit: 0x1, mask: ^uint32(0)},
		{val: 38, index: 1, bit: (1 << 6), mask: 0x7},
		{val: 34, index: 1, bit: (1 << 2), mask: 0x7},
		{val: 65, index: 2, bit: (1 << 1), mask: 0},
	}
	for i, v := range vs {
		index, bit, err := bitOf(Inheritable, v.val)
		if err != nil {
			t.Fatalf("[%d] %v(%d) - not bitOf: %v", i, v.val, v.val, err)
		} else if index != v.index {
			t.Errorf("[%d] %v(%d) - index: got=%d want=%d", i, v.val, v.val, index, v.index)
		}
		if bit != v.bit {
			t.Errorf("[%d] %v(%d) - bit: got=%b want=%b", i, v.val, v.val, bit, v.bit)
		}
		if mask := allMask(index); mask != v.mask {
			t.Errorf("[%d] %v(%d) - mask: got=%b want=%b", i, v.val, v.val, mask, v.mask)
		}
	}
}

func TestString(t *testing.T) {
	a := CHOWN
	if got, want := a.String(), "cap_chown"; got != want {
		t.Fatalf("pretty basic failure: got=%q, want=%q", got, want)
	}
}

func TestText(t *testing.T) {
	vs := []struct {
		from, to string
		err      error
	}{
		{"", "", ErrBadText},
		{"=", "=", nil},
		{"= cap_chown+iep cap_chown-i", "= cap_chown+ep", nil},
		{"= cap_setfcap,cap_chown+iep cap_chown-i", "= cap_setfcap+eip cap_chown+ep", nil},
	}
	for i, v := range vs {
		c, err := FromText(v.from)
		if err != v.err {
			t.Errorf("[%d] parsing %q failed: got=%v, want=%v", i, v.from, err, v.err)
			continue
		}
		if err != nil {
			continue
		}
		to := c.String()
		if to != v.to {
			t.Errorf("[%d] failed to stringify cap: %q -> got=%q, want=%q", i, v.from, to, v.to)
		}
		if d, err := FromText(to); err != nil {
			t.Errorf("[%d] failed to reparse %q: %v", i, to, err)
		} else if got := d.String(); got != to {
			t.Errorf("[%d] failed to stringify %q getting %q", i, to, got)
		}
	}
}

func same(a, b *Set) error {
	if (a == nil) != (b == nil) {
		return fmt.Errorf("nil-ness miscompare: %q vs %v", a, b)
	}
	if a == nil {
		return nil
	}
	if a.nsRoot != b.nsRoot {
		return fmt.Errorf("capabilities differ in nsRoot: a=%d b=%d", a.nsRoot, b.nsRoot)
	}
	for i, f := range a.flat {
		g := b.flat[i]
		for s := Effective; s <= Inheritable; s++ {
			if got, want := f[s], g[s]; got != want {
				return fmt.Errorf("capabilities differ: a[%d].flat[%v]=0x%08x b[%d].flat[%v]=0x%08x", i, s, got, i, s, want)
			}
		}
	}
	return nil
}

func TestImportExport(t *testing.T) {
	wantQ := "=ep cap_chown-e 63+ip"
	if q, err := FromText(wantQ); err != nil {
		t.Fatalf("failed to parse %q: %v", wantQ, err)
	} else if gotQ := q.String(); gotQ != wantQ {
		t.Fatalf("static test failed %q -> q -> %q", wantQ, gotQ)
	}

	// Sanity check empty import/export.
	c := NewSet()
	if ex, err := c.Export(); err != nil {
		t.Fatalf("failed to export empty set: %v", err)
	} else if len(ex) != 5 {
		t.Fatalf("wrong length: got=%d want=%d", len(ex), 5)
	} else if im, err := Import(ex); err != nil {
		t.Fatalf("failed to import empty set: %v", err)
	} else if got, want := im.String(), c.String(); got != want {
		t.Fatalf("import != export: got=%q want=%q", got, want)
	}
	// Now keep flipping bits on and off and validate that all
	// forms of import/export work.
	for i := uint(0); i < 7000; i += 13 {
		s := Flag(i % 3)
		v := Value(i % (maxValues + 3))
		c.SetFlag(s, i&17 < 8, v)
		if ex, err := c.Export(); err != nil {
			t.Fatalf("[%d] failed to export (%q): %v", i, c, err)
		} else if im, err := Import(ex); err != nil {
			t.Fatalf("[%d] failed to import (%q) set: %v", i, c, err)
		} else if got, want := im.String(), c.String(); got != want {
			t.Fatalf("[%d] import != export: got=%q want=%q [%02x]", i, got, want, ex)
		} else if parsed, err := FromText(got); err != nil {
			t.Fatalf("[%d] failed to parse %q: %v", i, got, err)
		} else if err := same(c, parsed); err != nil {
			t.Fatalf("[%d] miscompare (%q vs. %q): %v", i, got, parsed, err)
		}
	}
}

func TestIAB(t *testing.T) {
	vs := []struct {
		text string
		bad  bool
	}{
		{text: "cup_full", bad: true},
		{text: ""},
		{text: "!%cap_chown"},
		{text: "!cap_chown,^cap_setuid"},
		{text: "cap_chown,cap_setuid"},
		{text: "^cap_chown,cap_setuid"},
		{text: "^cap_chown,!cap_setuid"},
	}
	for i, v := range vs {
		want := v.text
		iab, err := IABFromText(want)
		if err != nil {
			if v.bad {
				continue
			}
			t.Errorf("[%d] want=%q, got=%q", i, want, iab)
			continue
		}
		if got := iab.String(); got != want {
			t.Errorf("[%d] got=%q want=%q", i, got, want)
		}
	}

	one, err := GetPID(1)
	if err != nil {
		t.Fatalf("failed to get init's capabilities: %v", err)
	}
	iab := IABInit()
	iab.Fill(Amb, one, Permitted)
	for i := 0; i < words; i++ {
		if iab.i[i] != iab.a[i] {
			t.Errorf("[%d] i=0x%08x != a=0x%08x", i, iab.i[i], iab.a[i])
		}
	}
	one.ClearFlag(Inheritable)
	iab.Fill(Inh, one, Inheritable)
	for i := 0; i < words; i++ {
		if iab.i[i] != iab.a[i] {
			t.Errorf("[%d] i=0x%08x != a=0x%08x", i, iab.i[i], iab.a[i])
		}
	}

	for n := uint(0); n < 1000; n += 13 {
		enabled := ((n % 5) & 2) != 0
		vec := Vector(n % 3)
		c := Value(n % maxValues)
		if err := iab.SetVector(vec, enabled, c); err != nil {
			t.Errorf("[%d] failed to set vec=%v enabled=%v %q in %q", n, vec, enabled, c, iab)
			continue
		}
		replay, err := IABFromText(iab.String())
		if err != nil {
			t.Errorf("failed to replay: %v", err)
			continue
		}
		for i := 0; i < words; i++ {
			if replay.i[i] != iab.i[i] || replay.a[i] != iab.a[i] || replay.nb[i] != iab.nb[i] {
				t.Errorf("[%d,%d] got=%q want=%q", n, i, replay, iab)
			}
		}
	}
}
