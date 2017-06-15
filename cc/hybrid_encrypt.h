// Copyright 2017 Google Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
///////////////////////////////////////////////////////////////////////////////

#ifndef TINK_HYBRID_ENCRYPT_H_
#define TINK_HYBRID_ENCRYPT_H_

#include "cc/util/statusor.h"
#include "google/protobuf/stubs/stringpiece.h"

namespace crypto {
namespace tink {

///////////////////////////////////////////////////////////////////////////////
// The interface for hybrid encryption.
//
// Implementations of this interface are secure against adaptive
// chosen ciphertext attacks.  In addition to 'plaintext' the
// encryption takes an extra parameter 'context_info', which usually
// is public data implicit from the context, but should be bound to
// the resulting ciphertext: upon decryption the ciphertext allows for
// checking the integrity of 'context_info' (but there are no
// guarantees wrt. to secrecy or authenticity of 'context_info').
//
// 'context_info' can be empty or null, but to ensure the correct
// decryption of the resulting ciphertext the same value must be
// provided for decryption operation (cf. HybridDecrypt-interface).
//
// A concrete instantiation of this interface can implement the
// binding of 'context_info' to the ciphertext in various ways, for
// example:
//
// - use 'context_info' as "associated data"-input for the employed
//   AEAD symmetric encryption (cf. https://tools.ietf.org/html/rfc5116).
// - use 'context_info' as "CtxInfo"-input for HKDF (if the implementation uses
//   HKDF as key derivation function, cf. https://tools.ietf.org/html/rfc5869).
class HybridEncrypt {
 public:
  // Encrypts 'plaintext' binding 'context_info' to the resulting ciphertext.
  virtual util::StatusOr<std::string> Encrypt(
      google::protobuf::StringPiece plaintext,
      google::protobuf::StringPiece context_info) const = 0;

  virtual ~HybridEncrypt() {}
};

}  // namespace tink
}  // namespace crypto

#endif  // TINK_HYBRID_ENCRYPT_H_