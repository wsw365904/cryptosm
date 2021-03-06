// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package cryptosm collects common cryptographic constants.
package cryptosm

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"github.com/wsw365904/cryptosm/sm3"
	"golang.org/x/crypto/blake2b"
	"golang.org/x/crypto/blake2s"
	"golang.org/x/crypto/ripemd160"
	"golang.org/x/crypto/sha3"
	"hash"
	"strconv"
)

func init() {
	registerHash(MD4, nil)
	registerHash(MD5, md5.New)
	registerHash(SHA1, sha1.New)
	registerHash(SHA224, sha256.New224)
	registerHash(SHA256, sha256.New)
	registerHash(SHA384, sha512.New384)
	registerHash(SHA512, sha512.New)
	registerHash(MD5SHA1, nil)
	registerHash(RIPEMD160, ripemd160.New)
	registerHash(SHA3_224, sha3.New224)
	registerHash(SHA3_256, sha3.New256)
	registerHash(SHA3_384, sha3.New384)
	registerHash(SHA3_512, sha3.New512)
	registerHash(SHA512_224, sha512.New512_224)
	registerHash(SHA512_256, sha512.New512_256)
	newHash256 := func() hash.Hash {
		h, _ := blake2s.New256(nil)
		return h
	}
	registerHash(BLAKE2s_256, newHash256)

	newHash256 = func() hash.Hash {
		h, _ := blake2b.New256(nil)
		return h
	}
	newHash384 := func() hash.Hash {
		h, _ := blake2b.New384(nil)
		return h
	}

	newHash512 := func() hash.Hash {
		h, _ := blake2b.New512(nil)
		return h
	}

	registerHash(BLAKE2b_256, newHash256)
	registerHash(BLAKE2b_384, newHash384)
	registerHash(BLAKE2b_512, newHash512)

	registerHash(SM3, sm3.New)
}

// Hash identifies a cryptographic hash function that is implemented in another
// package.
type Hash uint

// HashFunc simply returns the value of h so that Hash implements SignerOpts.
func (h Hash) HashFunc() Hash {
	return h
}

func (h Hash) String() string {
	switch h {
	case MD4:
		return "MD4"
	case MD5:
		return "MD5"
	case SHA1:
		return "SHA-1"
	case SHA224:
		return "SHA-224"
	case SHA256:
		return "SHA-256"
	case SHA384:
		return "SHA-384"
	case SHA512:
		return "SHA-512"
	case MD5SHA1:
		return "MD5+SHA1"
	case RIPEMD160:
		return "RIPEMD-160"
	case SHA3_224:
		return "SHA3-224"
	case SHA3_256:
		return "SHA3-256"
	case SHA3_384:
		return "SHA3-384"
	case SHA3_512:
		return "SHA3-512"
	case SHA512_224:
		return "SHA-512/224"
	case SHA512_256:
		return "SHA-512/256"
	case BLAKE2s_256:
		return "BLAKE2s-256"
	case BLAKE2b_256:
		return "BLAKE2b-256"
	case BLAKE2b_384:
		return "BLAKE2b-384"
	case BLAKE2b_512:
		return "BLAKE2b-512"
	case SM3:
		return "SM3"
	default:
		return "unknown hash value " + strconv.Itoa(int(h))
	}
}

const (
	MD4         Hash = 1 + iota // import golang.org/x/crypto/md4
	MD5                         // import crypto/md5
	SHA1                        // import crypto/sha1
	SHA224                      // import crypto/sha256
	SHA256                      // import crypto/sha256
	SHA384                      // import crypto/sha512
	SHA512                      // import crypto/sha512
	MD5SHA1                     // no implementation; MD5+SHA1 used for TLS RSA
	RIPEMD160                   // import golang.org/x/crypto/ripemd160
	SHA3_224                    // import golang.org/x/crypto/sha3
	SHA3_256                    // import golang.org/x/crypto/sha3
	SHA3_384                    // import golang.org/x/crypto/sha3
	SHA3_512                    // import golang.org/x/crypto/sha3
	SHA512_224                  // import crypto/sha512
	SHA512_256                  // import crypto/sha512
	BLAKE2s_256                 // import golang.org/x/crypto/blake2s
	BLAKE2b_256                 // import golang.org/x/crypto/blake2b
	BLAKE2b_384                 // import golang.org/x/crypto/blake2b
	BLAKE2b_512                 // import golang.org/x/crypto/blake2b
	SM3
	maxHash
)

var digestSizes = []uint8{
	MD4:         16,
	MD5:         16,
	SHA1:        20,
	SHA224:      28,
	SHA256:      32,
	SHA384:      48,
	SHA512:      64,
	SHA512_224:  28,
	SHA512_256:  32,
	SHA3_224:    28,
	SHA3_256:    32,
	SHA3_384:    48,
	SHA3_512:    64,
	MD5SHA1:     36,
	RIPEMD160:   20,
	BLAKE2s_256: 32,
	BLAKE2b_256: 32,
	BLAKE2b_384: 48,
	BLAKE2b_512: 64,
	SM3:         32,
}

// Size returns the length, in bytes, of a digest resulting from the given hash
// function. It doesn't require that the hash function in question be linked
// into the program.
func (h Hash) Size() int {
	if h > 0 && h < maxHash {
		return int(digestSizes[h])
	}
	panic("crypto: Size of unknown hash function")
}

var hashes = make([]func() hash.Hash, maxHash)

// New returns a new hash.Hash calculating the given hash function. New panics
// if the hash function is not linked into the binary.
func (h Hash) New() hash.Hash {
	if h > 0 && h < maxHash {
		f := hashes[h]
		if f != nil {
			return f()
		}
	}
	panic("crypto: requested hash function #" + strconv.Itoa(int(h)) + " is unavailable")
}

// Available reports whether the given hash function is linked into the binary.
func (h Hash) Available() bool {
	return h < maxHash && hashes[h] != nil
}

// RegisterHash registers a function that returns a new instance of the given
// hash function. This is intended to be called from the init function in
// packages that implement hash functions.
func registerHash(h Hash, f func() hash.Hash) {
	if h >= maxHash {
		panic("cryptosm: RegisterHash of unknown hash function")
	}
	hashes[h] = f
}
