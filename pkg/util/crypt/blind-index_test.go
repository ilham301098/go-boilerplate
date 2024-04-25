package crypt

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tink-crypto/tink-go/v2/keyset"
	"github.com/tink-crypto/tink-go/v2/mac"
)

func TestBlindIndexes(t *testing.T) {
	mgr := keyset.NewManager()
	hid, err := mgr.Add(mac.HMACSHA256Tag256KeyTemplate())
	require.NoError(t, err, "should add mac handle")
	err = mgr.SetPrimary(hid)
	require.NoError(t, err, "should set primary handle")
	handle, err := mgr.Handle()
	require.NoError(t, err, "should obtain mac handle")
	m, err := mac.New(handle)
	require.NoError(t, err, "should create mac primitive")

	data := []byte("asdasjdiu9lksdlfkjasopfijaposdpasi09ie283u023hj02i0t83089tu045jt054050j")
	v, err := m.ComputeMAC(data[:])
	require.NoError(t, err, "should compute mac")

	bidx, err := NewBIDX(handle, 0)
	require.NoError(t, err)
	v1, err := bidx.ComputePrimary(data)
	require.NoError(t, err)
	assert.Equal(t, v, v1)

	hid, err = mgr.Add(mac.HMACSHA256Tag128KeyTemplate())
	require.NoError(t, err, "should add new mac handle")
	err = mgr.SetPrimary(hid)
	require.NoError(t, err, "should set new primary handle")
	handle, err = mgr.Handle()
	require.NoError(t, err, "should obtain new mac handle")

	for i := 0; i < 10; i++ {
		vs, err := GetBlindIdxs(handle, data[:], len(v))
		require.NoError(t, err, "should compute multiple mac")
		vs1, err := bidx.ComputeAll(data)
		require.NoError(t, err, "should compute multiple mac")

		assert.Contains(t, vs, v, "should contain previous mac")
		assert.Contains(t, vs1, v, "should contain previous mac")
	}
}

func TestCopyBIDXWithLen(t *testing.T) {
	h, err := keyset.NewHandle(mac.HMACSHA256Tag256KeyTemplate())
	require.NoError(t, err, "should return new handle")
	bidx, err := NewPrimitiveBIDX(h)
	require.NoError(t, err, "should return bidx handle")
	v, err := bidx.ComputePrimary([]byte("data"))
	require.NoError(t, err, "should produce bidx")

	var bidxN BIDX = bidx
	for i := len(v); i > 0; i-- {
		bidxN, err = CopyBIDXWithLen(bidxN, i)
		require.NoError(t, err, "should return bidx handle")
		vn, err := bidxN.ComputePrimary([]byte("data"))
		require.NoError(t, err, "should produce bidx")
		if i == 0 {
			assert.Len(t, vn, len(v))
		} else {
			assert.Len(t, vn, i)
		}
	}

	v1, err := bidx.ComputePrimary([]byte("data"))
	require.NoError(t, err, "should produce bidx")
	assert.Equal(t, v, v1, "should produce the same bidx")
}
