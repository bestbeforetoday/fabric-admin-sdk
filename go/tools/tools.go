/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package tools

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fabric-admin-sdk/internal/configtxgen/encoder"
	"fabric-admin-sdk/internal/configtxgen/genesisconfig"
	"fabric-admin-sdk/internal/pkg/identity"
	"fmt"
	"io/ioutil"

	"github.com/gogo/protobuf/proto"
	cb "github.com/hyperledger/fabric-protos-go/common"
	"github.com/hyperledger/fabric-protos-go/msp"
	"github.com/hyperledger/fabric/common/flogging"
)

var logger = flogging.MustGetLogger("common.tools.configtxgen")

// configtxGen
// base on Profile return block
func ConfigTxGen(config *genesisconfig.Profile, channelID string) (*cb.Block, error) {
	pgen, err := encoder.NewBootstrapper(config)
	if err != nil {
		return nil, err
	}
	genesisBlock := pgen.GenesisBlockForChannel(channelID)
	return genesisBlock, nil
}

// load profile
// file as file path
// profile_name name
func LoadProfile(configName, FABRIC_CFG_PATH string) (*genesisconfig.Profile, error) {
	return genesisconfig.Load(configName, FABRIC_CFG_PATH)
}

func CreateSigner(PrivKeyPath, SignCert, MSPID string) (*identity.CryptoImpl, error) {
	priv, err := GetPrivateKey(PrivKeyPath)
	if err != nil {
		return nil, err
	}

	cert, certBytes, err := GetCertificate(SignCert)
	if err != nil {
		return nil, err
	}

	id := &msp.SerializedIdentity{
		Mspid:   MSPID,
		IdBytes: certBytes,
	}

	name, err := proto.Marshal(id)
	if err != nil {
		return nil, err
	}

	//get signer
	CryptoImpl := &identity.CryptoImpl{
		Creator:  name,
		PrivKey:  priv,
		SignCert: cert,
	}

	return CryptoImpl, nil
}

func GetPrivateKey(f string) (*ecdsa.PrivateKey, error) {
	in, err := ioutil.ReadFile(f)
	if err != nil {
		return nil, err
	}

	k, err := PEMtoPrivateKey(in, []byte{})
	if err != nil {
		return nil, err
	}

	key, ok := k.(*ecdsa.PrivateKey)
	if !ok {
		return nil, errors.New("expecting ecdsa key")
	}

	return key, nil
}

func GetCertificate(f string) (*x509.Certificate, []byte, error) {
	in, err := ioutil.ReadFile(f)
	if err != nil {
		return nil, nil, err
	}

	block, _ := pem.Decode(in)

	c, err := x509.ParseCertificate(block.Bytes)
	return c, in, err
}

// PEMtoPrivateKey unmarshals a pem to private key
func PEMtoPrivateKey(raw []byte, pwd []byte) (interface{}, error) {
	if len(raw) == 0 {
		return nil, errors.New("Invalid PEM. It must be different from nil.")
	}
	block, _ := pem.Decode(raw)
	if block == nil {
		return nil, fmt.Errorf("Failed decoding PEM. Block must be different from nil. [% x]", raw)
	}

	// TODO: derive from header the type of the key

	if x509.IsEncryptedPEMBlock(block) {
		if len(pwd) == 0 {
			return nil, errors.New("Encrypted Key. Need a password")
		}

		decrypted, err := x509.DecryptPEMBlock(block, pwd)
		if err != nil {
			return nil, fmt.Errorf("Failed PEM decryption [%s]", err)
		}

		key, err := DERToPrivateKey(decrypted)
		if err != nil {
			return nil, err
		}
		return key, err
	}

	cert, err := DERToPrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return cert, err
}

// DERToPrivateKey unmarshals a der to private key
func DERToPrivateKey(der []byte) (key interface{}, err error) {

	if key, err = x509.ParsePKCS1PrivateKey(der); err == nil {
		return key, nil
	}

	if key, err = x509.ParsePKCS8PrivateKey(der); err == nil {
		switch key.(type) {
		case *ecdsa.PrivateKey:
			return
		default:
			return nil, errors.New("Found unknown private key type in PKCS#8 wrapping")
		}
	}

	if key, err = x509.ParseECPrivateKey(der); err == nil {
		return
	}

	return nil, errors.New("Invalid key type. The DER must contain an ecdsa.PrivateKey")
}
