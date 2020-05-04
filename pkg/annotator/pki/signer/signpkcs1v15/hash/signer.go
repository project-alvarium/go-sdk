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
		name = "MD4"
	case crypto.MD5:
		name = "MD5"
	case crypto.SHA1:
		name = "SHA1"
	case crypto.SHA224:
		name = "SHA224"
	case crypto.SHA256:
		name = "SHA256"
	case crypto.SHA384:
		name = "SHA384"
	case crypto.SHA512:
		name = "SHA512"
	case crypto.MD5SHA1:
		name = "MD5SHA1"
	case crypto.RIPEMD160:
		name = "RIPEMD160"
	case crypto.SHA3_224:
		name = "SHA3_224"
	case crypto.SHA3_256:
		name = "SHA3_256"
	case crypto.SHA3_384:
		name = "SHA3_384"
	case crypto.SHA3_512:
		name = "SHA3_512"
	case crypto.SHA512_224:
		name = "SHA512_224"
	case crypto.SHA512_256:
		name = "SHA512_256"
	case crypto.BLAKE2s_256:
		name = "BLAKE2s_256"
	case crypto.BLAKE2b_256:
		name = "BLAKE2b_256"
	case crypto.BLAKE2b_384:
		name = "BLAKE2b_384"
	case crypto.BLAKE2b_512:
		name = "BLAKE2b_512"
	default:
		name = "unknown"
	}
	return
}

// ToSigner converts a string representation to a crypto.Hash (or nil if unable to do so).
func ToSigner(name string) (hash crypto.Hash) {
	switch name {
	case "MD4":
		hash = crypto.MD4
	case "MD5":
		hash = crypto.MD5
	case "SHA1":
		hash = crypto.SHA1
	case "SHA224":
		hash = crypto.SHA224
	case "SHA256":
		hash = crypto.SHA256
	case "SHA384":
		hash = crypto.SHA384
	case "SHA512":
		hash = crypto.SHA512
	case "MD5SHA1":
		hash = crypto.MD5SHA1
	case "RIPEMD160":
		hash = crypto.RIPEMD160
	case "SHA3_224":
		hash = crypto.SHA3_224
	case "SHA3_256":
		hash = crypto.SHA3_256
	case "SHA3_384":
		hash = crypto.SHA3_384
	case "SHA3_512":
		hash = crypto.SHA3_512
	case "SHA512_224":
		hash = crypto.SHA512_224
	case "SHA512_256":
		hash = crypto.SHA512_256
	case "BLAKE2s_256":
		hash = crypto.BLAKE2s_256
	case "BLAKE2b_256":
		hash = crypto.BLAKE2b_256
	case "BLAKE2b_384":
		hash = crypto.BLAKE2b_384
	case "BLAKE2b_512":
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
