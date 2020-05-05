/*******************************************************************************
 * Copyright 2020 Dell Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software distributed under the License
 * is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express
 * or implied. See the License for the specific language governing permissions and limitations under
 * the License.
 *******************************************************************************/

package hash

import "crypto"

// FromSigner converts a crypto.Hash value to a string representation.
func FromSigner(hash crypto.Hash) (name string) {
	switch hash {
	case crypto.MD4:
		name = "md4"
	case crypto.MD5:
		name = "md5"
	case crypto.SHA1:
		name = "sha1"
	case crypto.SHA224:
		name = "sha224"
	case crypto.SHA256:
		name = "sha256"
	case crypto.SHA384:
		name = "sha384"
	case crypto.SHA512:
		name = "sha512"
	case crypto.MD5SHA1:
		name = "md5sha1"
	case crypto.RIPEMD160:
		name = "ripemd160"
	case crypto.SHA3_224:
		name = "sha3_224"
	case crypto.SHA3_256:
		name = "sha3_256"
	case crypto.SHA3_384:
		name = "sha3_384"
	case crypto.SHA3_512:
		name = "sha3_512"
	case crypto.SHA512_224:
		name = "sha512_224"
	case crypto.SHA512_256:
		name = "sha512_256"
	case crypto.BLAKE2s_256:
		name = "blake2s_256"
	case crypto.BLAKE2b_256:
		name = "blake2b_256"
	case crypto.BLAKE2b_384:
		name = "blake2b_384"
	case crypto.BLAKE2b_512:
		name = "blake2b_512"
	default:
		name = "unknown"
	}
	return
}

// ToSigner converts a string representation to a crypto.Hash (or nil if unable to do so).
func ToSigner(name string) (hash crypto.Hash) {
	switch name {
	case "md4":
		hash = crypto.MD4
	case "md5":
		hash = crypto.MD5
	case "sha1":
		hash = crypto.SHA1
	case "sha224":
		hash = crypto.SHA224
	case "sha256":
		hash = crypto.SHA256
	case "sha384":
		hash = crypto.SHA384
	case "sha512":
		hash = crypto.SHA512
	case "md5sha1":
		hash = crypto.MD5SHA1
	case "ripemd160":
		hash = crypto.RIPEMD160
	case "sha3_224":
		hash = crypto.SHA3_224
	case "sha3_256":
		hash = crypto.SHA3_256
	case "sha3_384":
		hash = crypto.SHA3_384
	case "sha3_512":
		hash = crypto.SHA3_512
	case "sha512_224":
		hash = crypto.SHA512_224
	case "sha512_256":
		hash = crypto.SHA512_256
	case "blake2s_256":
		hash = crypto.BLAKE2s_256
	case "blake2b_256":
		hash = crypto.BLAKE2b_256
	case "blake2b_384":
		hash = crypto.BLAKE2b_384
	case "blake2b_512":
		hash = crypto.BLAKE2b_512
	default:
		hash = 0
	}
	return
}

// SupportedSigner returns a slice of supported crypto.Hash.
func SupportedSigner() []crypto.Hash {
	return []crypto.Hash{
		crypto.MD4,
		crypto.MD5,
		crypto.SHA1,
		crypto.SHA224,
		crypto.SHA256,
		crypto.SHA384,
		crypto.SHA512,
		crypto.MD5SHA1,
		crypto.RIPEMD160,
		crypto.SHA3_224,
		crypto.SHA3_256,
		crypto.SHA3_384,
		crypto.SHA3_512,
		crypto.SHA512_224,
		crypto.SHA512_256,
		crypto.BLAKE2s_256,
		crypto.BLAKE2b_256,
		crypto.BLAKE2b_384,
		crypto.BLAKE2b_512,
	}
}
