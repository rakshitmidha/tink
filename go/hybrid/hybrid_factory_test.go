// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
////////////////////////////////////////////////////////////////////////////////

package hybrid

import (
	"bytes"
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/google/tink/go/aead"
	"github.com/google/tink/go/subtle/random"
	"github.com/google/tink/go/testkeyset"
	"github.com/google/tink/go/testutil"
	commonpb "github.com/google/tink/proto/common_go_proto"
	tinkpb "github.com/google/tink/proto/tink_go_proto"
)

func TestHybridFactoryTest(t *testing.T) {
	c := commonpb.EllipticCurveType_NIST_P256
	ht := commonpb.HashType_SHA256
	primaryPtFmt := commonpb.EcPointFormat_UNCOMPRESSED
	rawPtFmt := commonpb.EcPointFormat_COMPRESSED
	primaryDek := aead.AES128CTRHMACSHA256KeyTemplate()
	rawDek := aead.AES128CTRHMACSHA256KeyTemplate()
	primarySalt := []byte("some salt")
	rawSalt := []byte("other salt")

	primaryPrivProto, err := testutil.GenerateECIESAEADHKDFPrivateKey(c, ht, primaryPtFmt, primaryDek, primarySalt)
	if err != nil {
		t.Error(err)
	}
	sPrimaryPriv, err := proto.Marshal(primaryPrivProto)
	if err != nil {
		t.Error(err)
	}
	sPrimaryPub, err := proto.Marshal(primaryPrivProto.PublicKey)
	if err != nil {
		t.Error(err)
	}

	primaryPrivKey := testutil.NewKey(
		testutil.NewKeyData(eciesAEADHKDFPrivateKeyTypeURL, sPrimaryPriv, tinkpb.KeyData_ASYMMETRIC_PRIVATE),
		tinkpb.KeyStatusType_ENABLED, 8, tinkpb.OutputPrefixType_RAW)
	primaryPubKey := testutil.NewKey(
		testutil.NewKeyData(eciesAEADHKDFPublicKeyTypeURL, sPrimaryPub, tinkpb.KeyData_ASYMMETRIC_PUBLIC),
		tinkpb.KeyStatusType_ENABLED, 42, tinkpb.OutputPrefixType_RAW)

	rawPrivProto, err := testutil.GenerateECIESAEADHKDFPrivateKey(c, ht, rawPtFmt, rawDek, rawSalt)
	if err != nil {
		t.Error(err)
	}
	sRawPriv, err := proto.Marshal(rawPrivProto)
	if err != nil {
		t.Error(err)
	}
	sRawPub, err := proto.Marshal(rawPrivProto.PublicKey)
	if err != nil {
		t.Error(err)
	}
	rawPrivKey := testutil.NewKey(
		testutil.NewKeyData(eciesAEADHKDFPrivateKeyTypeURL, sRawPriv, tinkpb.KeyData_ASYMMETRIC_PRIVATE),
		tinkpb.KeyStatusType_ENABLED, 11, tinkpb.OutputPrefixType_RAW)
	rawPubKey := testutil.NewKey(
		testutil.NewKeyData(eciesAEADHKDFPublicKeyTypeURL, sRawPub, tinkpb.KeyData_ASYMMETRIC_PUBLIC),
		tinkpb.KeyStatusType_ENABLED, 43, tinkpb.OutputPrefixType_RAW)

	pubKeys := []*tinkpb.Keyset_Key{primaryPubKey, rawPubKey}
	pubKeyset := testutil.NewKeyset(pubKeys[0].KeyId, pubKeys)
	khPub, err := testkeyset.NewHandle(pubKeyset)
	if err != nil {
		t.Error(err)
	}

	privKeys := []*tinkpb.Keyset_Key{primaryPrivKey, rawPrivKey}
	privKeyset := testutil.NewKeyset(privKeys[0].KeyId, privKeys)
	khPriv, err := testkeyset.NewHandle(privKeyset)
	if err != nil {
		t.Error(err)
	}
	e, err := NewHybridEncrypt(khPub)
	if err != nil {
		t.Error(err)
	}
	d, err := NewHybridDecrypt(khPriv)
	if err != nil {
		t.Error(err)
	}

	for i := 0; i < 1000; i++ {
		pt := random.GetRandomBytes(20)
		ci := random.GetRandomBytes(20)
		ct, err := e.Encrypt(pt, ci)
		if err != nil {
			t.Error(err)
		}
		gotpt, err := d.Decrypt(ct, ci)
		if err != nil {
			t.Error(err)
		}
		if !bytes.Equal(pt, gotpt) {
			t.Error(err)
		}
	}
}
