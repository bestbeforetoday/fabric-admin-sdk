/*
Copyright IBM Corp. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package chaincode

import "github.com/hyperledger/fabric-protos-go-apiv2/peer"

const (
	lifecycleChaincodeName              = "_lifecycle"
	approveTransactionName              = "ApproveChaincodeDefinitionForMyOrg"
	commitTransactionName               = "CommitChaincodeDefinition"
	queryInstalledTransactionName       = "QueryInstalledChaincodes"
	checkCommitReadinessTransactionName = "CheckCommitReadiness"
	installTransactionName              = "InstallChaincode"
	// MetadataFile is the expected location of the metadata json document
	// in the top level of the chaincode package.
	MetadataFile = "metadata.json"

	// CodePackageFile is the expected location of the code package in the
	// top level of the chaincode package
	CodePackageFile = "code.tar.gz"
)

// Chaincode Define
type Definition struct {
	ChannelID                string
	InputTxID                string
	PackageID                string
	Name                     string
	Version                  string
	EndorsementPlugin        string
	EndorsementPolicy        string
	ValidationPlugin         string
	Sequence                 int64
	ValidationParameterBytes []byte
	InitRequired             bool
	CollectionConfigPackage  *peer.CollectionConfigPackage
}
