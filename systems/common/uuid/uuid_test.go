package uuid

import (
	"fmt"
	"reflect"
	"testing"
)

func TestEqual(t *testing.T) {
	u1 := UUID{0x6b, 0xa7, 0xb8, 0x10, 0x9d, 0xad, 0x11, 0xd1, 0x80, 0xb4, 0x00, 0xc0, 0x4f, 0xd4, 0x30, 0xc8}
	u2 := UUID{0x6b, 0xa7, 0xb8, 0x10, 0x9d, 0xad, 0x11, 0xd1, 0x80, 0xb4, 0x00, 0xc0, 0x4f, 0xd4, 0x30, 0xc8}

	if !Equal(u1, u2) {
		t.Errorf("expected %v == %v but did not", u1, u2)
	}
}

func TestFromBytes(t *testing.T) {
	in := []byte{0x99, 0xa7, 0x08, 0x50, 0xad, 0x11, 0x11, 0x2f, 0xa0, 0x74, 0x80, 0xc0, 0x3d, 0x55, 0x4e, 0xcc}
	want := UUID{0x99, 0xa7, 0x08, 0x50, 0xad, 0x11, 0x11, 0x2f, 0xa0, 0x74, 0x80, 0xc0, 0x3d, 0x55, 0x4e, 0xcc}

	got, err := FromBytes(in)

	if err != nil {
		t.Error("expected no errors but got one")
	}

	if want != got {
		t.Errorf("expected %v but got %v", want, got)
	}
}

func TestFromBytesOrNil(t *testing.T) {
	t.Run("InvalidUUIDBytes", func(t *testing.T) {
		in := []byte{0x99, 0xa7, 0x08, 0x50}

		got := FromBytesOrNil(in)

		if got != Nil {
			t.Errorf("expected %v but got %v", Nil, got)
		}
	})

	t.Run("ValidUUIDBytes", func(t *testing.T) {
		in := []byte{0x91, 0x77, 0xb8, 0x59, 0xad, 0x13, 0xab, 0xe1, 0xa0, 0x04, 0x80, 0x20, 0x3d, 0xff, 0x6e, 0xcc}
		want := UUID{0x91, 0x77, 0xb8, 0x59, 0xad, 0x13, 0xab, 0xe1, 0xa0, 0x04, 0x80, 0x20, 0x3d, 0xff, 0x6e, 0xcc}

		got := FromBytesOrNil(in)

		if got != want {
			t.Errorf("expected %v but got %v", want, got)
		}
	})
}

func TestFromString(t *testing.T) {
	in := "6ba7b810-9dad-11d1-80b4-00c04fd430c8"
	want := UUID{0x6b, 0xa7, 0xb8, 0x10, 0x9d, 0xad, 0x11, 0xd1, 0x80, 0xb4, 0x00, 0xc0, 0x4f, 0xd4, 0x30, 0xc8}

	got, err := FromString(in)

	if err != nil {
		t.Error("expected no errors but got one")
	}

	if want != got {
		t.Errorf("expected %v but got %v", want, got)
	}
}

func TestFromStringOrNil(t *testing.T) {
	t.Run("InvalidUUIDString", func(t *testing.T) {
		in := "ukama"

		got := FromStringOrNil(in)

		if got != Nil {
			t.Errorf("expected %v but got %v", Nil, got)
		}
	})

	t.Run("ValidUUIDString", func(t *testing.T) {
		in := "6b77b859-9dad-13d1-b404-c0203dff6ec1"
		want := UUID{0x6b, 0x77, 0xb8, 0x59, 0x9d, 0xad, 0x13, 0xd1, 0xb4, 0x04, 0xc0, 0x20, 0x3d, 0xff, 0x6e, 0xc1}

		got := FromStringOrNil(in)

		if got != want {
			t.Errorf("expected %v but got %v", want, got)
		}
	})
}

func TestMust(t *testing.T) {
	t.Run("PanicHelper", func(t *testing.T) {
		defer func() {
			v := recover()

			if v == nil {
				t.Error("expected panic but did not")
			}
		}()

		Must(func() (UUID, error) {
			return Nil, fmt.Errorf("uuid: expected error")
		}())
	})

	t.Run("NonPanicHelper", func(t *testing.T) {
		defer func() {
			v := recover()

			if v != nil {
				t.Error("expected no panic  but did")
			}
		}()

		Must(func() (UUID, error) {
			return Nil, nil
		}())
	})
}

func TestNewV1(t *testing.T) {
	got := NewV1().Version()
	if got != 1 {
		t.Errorf("expected version %d gut got %d", 1, got)
	}
}

func TestNewV2(t *testing.T) {
	got := NewV2(1).Version()
	if got != 2 {
		t.Errorf("expected version %d gut got %d", 2, got)
	}
}

func TestNewV3(t *testing.T) {
	defer func() {
		v := recover()

		if v != nil {
			t.Error("expected no panic  but did")
		}
	}()

	// Default URL Namespace
	NamespaceURL := Must(FromString("6ba7b811-9dad-11d1-80b4-00c04fd430c8"))

	got := NewV3(NamespaceURL, "MD5-nonce").Version()
	if got != 3 {
		t.Errorf("expected version %d gut got %d", 3, got)
	}
}

func TestNewV4(t *testing.T) {
	got := NewV4().Version()
	if got != 4 {
		t.Errorf("expected version %d gut got %d", 4, got)
	}
}

func TestNewV5(t *testing.T) {
	defer func() {
		v := recover()

		if v != nil {
			t.Error("expected no panic  but did")
		}
	}()

	// Default DNS Namespace
	NamespaceDNS := Must(FromString("6ba7b810-9dad-11d1-80b4-00c04fd430c8"))

	got := NewV5(NamespaceDNS, "SHA1-nonce").Version()
	if got != 5 {
		t.Errorf("expected version %d gut got %d", 5, got)
	}
}

func TestUUID_Bytes(t *testing.T) {
	in := UUID{0x99, 0xa7, 0x08, 0x50, 0xad, 0x11, 0x11, 0x2f, 0xa0, 0x74, 0x80, 0xc0, 0x3d, 0x55, 0x4e, 0xcc}
	want := []byte{0x99, 0xa7, 0x08, 0x50, 0xad, 0x11, 0x11, 0x2f, 0xa0, 0x74, 0x80, 0xc0, 0x3d, 0x55, 0x4e, 0xcc}

	got := in.Bytes()

	if !reflect.DeepEqual(want, got) {
		t.Errorf("expected %v but got %v", want, got)
	}
}

func TestUUID_MarshalBinary(t *testing.T) {
	in := UUID{0x99, 0xa7, 0x08, 0x50, 0xad, 0x11, 0x11, 0x2f, 0xa0, 0x74, 0x80, 0xc0, 0x3d, 0x55, 0x4e, 0xcc}
	want := []byte{0x99, 0xa7, 0x08, 0x50, 0xad, 0x11, 0x11, 0x2f, 0xa0, 0x74, 0x80, 0xc0, 0x3d, 0x55, 0x4e, 0xcc}

	got, err := in.MarshalBinary()

	if err != nil {
		t.Error("expected no errors but got one")
	}

	if !reflect.DeepEqual(want, got) {
		t.Errorf("expected %v but got %v", want, got)
	}
}

func TestUUID_UnmarshalBinary(t *testing.T) {
	t.Run("ValidByte", func(t *testing.T) {
		b := []byte{0x99, 0xa7, 0x08, 0x50, 0xad, 0x11, 0x11, 0x2f, 0xa0, 0x74, 0x80, 0xc0, 0x3d, 0x55, 0x4e, 0xcc}
		want := UUID{0x99, 0xa7, 0x08, 0x50, 0xad, 0x11, 0x11, 0x2f, 0xa0, 0x74, 0x80, 0xc0, 0x3d, 0x55, 0x4e, 0xcc}
		u := UUID{}

		err := u.UnmarshalBinary(b)

		if err != nil {
			t.Error("expected no errors but got one")
		}

		if !Equal(u, want) {
			t.Errorf("expected %v but got %v", want, u)
		}
	})

	t.Run("InvalidValidByte", func(t *testing.T) {
		b := []byte{}
		u := UUID{}

		err := u.UnmarshalBinary(b)

		if err == nil {
			t.Error("expected an error but did not get one")
		}
	})
}

func TestUUID_MarshalText(t *testing.T) {
	in := UUID{0x6b, 0xa7, 0xb8, 0x10, 0x9d, 0xad, 0x11, 0xd1, 0x80, 0xb4, 0x00, 0xc0, 0x4f, 0xd4, 0x30, 0xc8}
	want := []byte("6ba7b810-9dad-11d1-80b4-00c04fd430c8")

	got, err := in.MarshalText()

	if err != nil {
		t.Error("expected no errors but got one")
	}

	if !reflect.DeepEqual(want, got) {
		t.Errorf("expected %v but got %v", want, got)
	}
}

func TestUUID_UnmarshalText(t *testing.T) {
	t.Run("ValidStringByte", func(t *testing.T) {
		b := []byte("6ba7b810-9dad-11d1-80b4-00c04fd430c8")
		want := UUID{0x6b, 0xa7, 0xb8, 0x10, 0x9d, 0xad, 0x11, 0xd1, 0x80, 0xb4, 0x00, 0xc0, 0x4f, 0xd4, 0x30, 0xc8}
		u := UUID{}

		err := u.UnmarshalText(b)

		if err != nil {
			t.Error("expected no errors but got one")
		}

		if !Equal(u, want) {
			t.Errorf("expected %v but got %v", want, u)
		}
	})

	t.Run("InvalidValidStringByte", func(t *testing.T) {
		b := []byte("")
		u := UUID{}

		err := u.UnmarshalText(b)

		if err == nil {
			t.Error("expected an error but did not get one")
		}
	})
}

func TestUUID_String(t *testing.T) {
	in := UUID{0x6b, 0x77, 0xb8, 0x59, 0x9d, 0xad, 0x13, 0xd1, 0xb4, 0x04, 0xc0, 0x20, 0x3d, 0xff, 0x6e, 0xc1}
	want := "6b77b859-9dad-13d1-b404-c0203dff6ec1"

	got := in.String()

	if got != want {
		t.Errorf("expected %v but got %v", want, got)
	}
}

func TestUUID_SetAndGetVariant(t *testing.T) {
	// The Nil UUID always matches all variants.
	u := UUID{}

	for i := 0; i < 4; i++ {
		u.SetVariant(byte(i))
		v := u.Variant()

		if v != byte(i) {
			t.Errorf("expected %d but got %d", i, v)
		}
	}
}

func TestUUID_SetAndGetVersion(t *testing.T) {
	// The Nil UUID always matches all variants.
	u := UUID{}

	for i := 1; i < 6; i++ {
		u.SetVersion(byte(i))
		v := u.Version()

		if v != byte(i) {
			t.Errorf("expected %d but got %d", i, v)
		}
	}
}

func TestNullUUID_Value(t *testing.T) {
	t.Run("Null", func(t *testing.T) {
		u := NullUUID{}

		val, err := u.Value()

		if err != nil {
			t.Error("expected no error but got one")
		}

		if val != nil {
			t.Errorf("expected nil value but got %v", val)
		}
	})

	t.Run("NonNull", func(t *testing.T) {
		u := NullUUID{UUID{0x6b, 0xa7, 0xb8, 0x10, 0x9d, 0xad, 0x11, 0xd1, 0x80, 0xb4, 0x00, 0xc0, 0x4f, 0xd4, 0x30, 0xc8}, true}

		val, err := u.Value()

		if err != nil {
			t.Error("expected no error but got one")
		}

		if val == nil {
			t.Errorf("expected non nil value but got %v", val)
		}

	})
}

func TestNullUUID_Scan(t *testing.T) {
	t.Run("ScanNil", func(t *testing.T) {
		u := NullUUID{UUID{0x6b, 0xa7, 0xb8, 0x10, 0x9d, 0xad, 0x11, 0xd1, 0x80, 0xb4, 0x00, 0xc0, 0x4f, 0xd4, 0x30, 0xc8}, true}

		err := u.Scan(nil)

		if err != nil {
			t.Error("expected no error but got one")
		}

		if u.Valid {
			t.Errorf("expected u.Valid to be %t but got %t", false, u.Valid)
		}

		if u.UUID != Nil {
			t.Errorf("expected UUID Nil value but got %v", u.UUID)
		}
	})

	t.Run("ScanNonNil", func(t *testing.T) {
		want := UUID{0x6b, 0xa7, 0xb8, 0x10, 0x9d, 0xad, 0x11, 0xd1, 0x80, 0xb4, 0x00, 0xc0, 0x4f, 0xd4, 0x30, 0xc8}
		s := "6ba7b810-9dad-11d1-80b4-00c04fd430c8"

		u := NullUUID{}

		err := u.Scan(s)

		if err != nil {
			t.Error("expected no error but got one")
		}

		if !u.Valid {
			t.Errorf("expected u.Valid to be %t but got %t", true, u.Valid)
		}

		if u.UUID != want {
			t.Errorf("expected UUID Nil value but got %v", u.UUID)
		}
	})
}

func TestUUID_Value(t *testing.T) {
	t.Run("UUIDIsNotNil", func(t *testing.T) {
		u, err := FromString("6ba7b810-9dad-11d1-80b4-00c04fd430c8")

		if err != nil {
			t.Error("expected no error but got one")
		}

		got, err := u.Value()
		if err != nil {
			t.Error("expected no error but got one")
		}

		if u.String() != got {
			t.Errorf("expected %q but got %q", u.String(), got)
		}
	})

	t.Run("UUIDIsNil", func(t *testing.T) {
		u := UUID{}
		got, err := u.Value()

		if err != nil {
			t.Error("expected no error but got one")
		}

		if Nil.String() != got {
			t.Errorf("expected %q but got %q", u.String(), got)
		}
	})
}

func TestUUID_Scan(t *testing.T) {
	t.Run("ScanValidByte", func(t *testing.T) {
		b := []byte{0x6b, 0xa7, 0xb8, 0x10, 0x9d, 0xad, 0x11, 0xd1, 0x80, 0xb4, 0x00, 0xc0, 0x4f, 0xd4, 0x30, 0xc8}
		want := UUID{0x6b, 0xa7, 0xb8, 0x10, 0x9d, 0xad, 0x11, 0xd1, 0x80, 0xb4, 0x00, 0xc0, 0x4f, 0xd4, 0x30, 0xc8}
		u := UUID{}

		err := u.Scan(b)

		if err != nil {
			t.Error("expected no error but got one")
		}

		if !Equal(u, want) {
			t.Errorf("expected %v but got %v", want, u)
		}
	})

	t.Run("ScanInvalidValidByte", func(t *testing.T) {
		b := []byte{}
		u := UUID{}

		err := u.Scan(b)

		if err == nil {
			t.Error("expected an error but did not get one")
		}
	})

	t.Run("ScanValidString", func(t *testing.T) {
		s := "6ba7b810-9dad-11d1-80b4-00c04fd430c8"
		want := UUID{0x6b, 0xa7, 0xb8, 0x10, 0x9d, 0xad, 0x11, 0xd1, 0x80, 0xb4, 0x00, 0xc0, 0x4f, 0xd4, 0x30, 0xc8}
		u := UUID{}

		err := u.Scan(s)

		if err != nil {
			t.Error("expected no error but got one")
		}

		if !Equal(u, want) {
			t.Errorf("expected %v but got %v", want, u)
		}
	})

	t.Run("ScanInvalidValidString", func(t *testing.T) {
		s := ""
		u := UUID{}

		err := u.Scan(s)

		if err == nil {
			t.Error("expected an error but did not get one")
		}
	})

	t.Run("ScanValidText", func(t *testing.T) {
		b := []byte("6ba7b810-9dad-11d1-80b4-00c04fd430c8")
		want := UUID{0x6b, 0xa7, 0xb8, 0x10, 0x9d, 0xad, 0x11, 0xd1, 0x80, 0xb4, 0x00, 0xc0, 0x4f, 0xd4, 0x30, 0xc8}
		u := UUID{}

		err := u.Scan(b)

		if err != nil {
			t.Error("expected no error but got one")
		}

		if !Equal(u, want) {
			t.Errorf("expected %v but got %v", want, u)
		}
	})

	t.Run("ScanInvalidValidText", func(t *testing.T) {
		b := []byte("")
		u := UUID{}

		err := u.Scan(b)

		if err == nil {
			t.Error("expected an error but did not get one")
		}
	})

	t.Run("ScanUnsupportedValue", func(t *testing.T) {
		u := UUID{}

		err := u.Scan(true)

		if err == nil {
			t.Error("expected an error but did not get one")
		}
	})

	t.Run("ScanNil", func(t *testing.T) {
		u := UUID{}

		err := u.Scan(nil)

		if err == nil {
			t.Error("expected an error but did not get one")
		}
	})
}
