package mptrie

import (
	"context"
	"testing"
	"time"

	"github.com/iotexproject/go-pkgs/hash"
	"github.com/iotexproject/iotex-core/config"
	"github.com/iotexproject/iotex-core/db"
	"github.com/iotexproject/iotex-core/db/trie"
	"github.com/iotexproject/iotex-core/testutil"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

/* new version
func BenchmarkTrie(b *testing.B) {
	if testing.Short() {
		b.Skip("Skipping TestPressure in short mode.")
	}

	require := require.New(b)

	testPath, err := testutil.PathOfTempFile("test-kv-store.bolt")
	require.NoError(err)
	cfg := config.Default.DB
	cfg.DbPath = testPath
	kv, err := trie.NewKVStore("test", db.NewBoltDB(cfg))
	require.NoError(kv.Start(context.Background()))

	require.NoError(err)
	tr, err := New(KVStoreOption(kv), KeyLengthOption(20), AsyncOption())
	require.NoError(err)
	require.NoError(tr.Start(context.Background()))
	root, err := tr.RootHash()
	require.NoError(err)
	seed := time.Now().Nanosecond()

	b.ResetTimer()
	// insert 128k entries
	var k [32]byte
	k[0] = byte(seed)
	c := 0
	b.Logf("iter %d", b.N)
	for c = 0; c < b.N+3; c++ {
		k = hash.Hash256b(k[:])
		v := testV[k[0]&7]
		//if _, err := tr.Get(k[:20]); err == nil {
		//b.Logf("Warning: collision on k %x", k[:20])
		//break
		//}
		require.NoError(tr.Upsert(k[:20], v))
		newRoot, err := tr.RootHash()
		require.NoError(err)
		require.False(tr.IsEmpty())
		require.NotEqual(newRoot, root)
		root = newRoot
		vb, err := tr.Get(k[:20])
		require.NoError(err)
		require.Equal(v, vb)
	}
	b.Logf("iter+3 %d", c)
	// delete 128k entries
	var d [32]byte
	d[0] = byte(seed)
	// save the first 3, delete them last
	d1 := hash.Hash256b(d[:])
	d2 := hash.Hash256b(d1[:])
	d3 := hash.Hash256b(d2[:])
	d = d3
	for i := 0; i < c-3; i++ {
		d = hash.Hash256b(d[:])
		require.NoError(tr.Delete(d[:20]))
		newRoot, err := tr.RootHash()
		require.NoError(err)
		require.False(tr.IsEmpty())
		require.NotEqual(newRoot, root)
		root = newRoot
		_, err = tr.Get(d[:20])
		require.Equal(trie.ErrNotExist, errors.Cause(err))
	}
	require.NoError(tr.Delete(d1[:20]))
	require.NoError(tr.Delete(d2[:20]))
	require.NoError(tr.Delete(d3[:20]))

	b.StopTimer()
	// trie should fallback to empty
	//require.True(tr.IsEmpty())
	require.NoError(tr.Stop(context.Background()))
	require.NoError(kv.Stop(context.Background()))
	b.Logf("Warning: test %d entries", c)
}
*/

func BenchmarkTrie(b *testing.B) {
	if testing.Short() {
		b.Skip("Skipping TestPressure in short mode.")
	}

	require := require.New(b)

	testPath, err := testutil.PathOfTempFile("test-kv-store.bolt")
	require.NoError(err)
	cfg := config.Default.DB
	cfg.DbPath = testPath
	kv, err := trie.NewKVStore("test", db.NewBoltDB(cfg))
	require.NoError(kv.Start(context.Background()))

	require.NoError(err)
	tr, err := New(KVStoreOption(kv), KeyLengthOption(20))
	require.NoError(err)
	require.NoError(tr.Start(context.Background()))
	root := tr.RootHash()
	seed := time.Now().Nanosecond()

	b.ResetTimer()
	// insert 128k entries
	var k [32]byte
	k[0] = byte(seed)
	c := 0
	b.Logf("iter %d", b.N)
	for c = 0; c < b.N+3; c++ {
		k = hash.Hash256b(k[:])
		v := testV[k[0]&7]
		//if _, err := tr.Get(k[:20]); err == nil {
		//b.Logf("Warning: collision on k %x", k[:20])
		//break
		//}
		require.NoError(tr.Upsert(k[:20], v))
		newRoot := tr.RootHash()
		require.False(tr.IsEmpty())
		require.NotEqual(newRoot, root)
		root = newRoot
		vb, err := tr.Get(k[:20])
		require.NoError(err)
		require.Equal(v, vb)
	}
	b.Logf("iter+3 %d", c)
	// delete 128k entries
	var d [32]byte
	d[0] = byte(seed)
	// save the first 3, delete them last
	d1 := hash.Hash256b(d[:])
	d2 := hash.Hash256b(d1[:])
	d3 := hash.Hash256b(d2[:])
	d = d3
	for i := 0; i < c-3; i++ {
		d = hash.Hash256b(d[:])
		require.NoError(tr.Delete(d[:20]))
		newRoot := tr.RootHash()
		require.False(tr.IsEmpty())
		require.NotEqual(newRoot, root)
		root = newRoot
		_, err = tr.Get(d[:20])
		require.Equal(trie.ErrNotExist, errors.Cause(err))
	}
	require.NoError(tr.Delete(d1[:20]))
	require.NoError(tr.Delete(d2[:20]))
	require.NoError(tr.Delete(d3[:20]))

	b.StopTimer()
	// trie should fallback to empty
	//require.True(tr.IsEmpty())
	require.NoError(tr.Stop(context.Background()))
	require.NoError(kv.Stop(context.Background()))
	b.Logf("Warning: test %d entries", c)
}
