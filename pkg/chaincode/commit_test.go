/*
Copyright IBM Corp. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package chaincode_test

import (
	"errors"
	"fmt"

	"github.com/golang/mock/gomock"
	"github.com/hyperledger/fabric-admin-sdk/pkg/chaincode"
	"github.com/hyperledger/fabric-admin-sdk/pkg/chaincode/mock"
	"github.com/hyperledger/fabric-protos-go-apiv2/peer"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Commit", func() {

	var chaincodeDefinition chaincode.Definition
	var endorsementClients []peer.EndorserClient
	BeforeEach(func() {
		chaincodeDefinition = chaincode.Definition{
			Name:      "CHAINCODE",
			Version:   "1.0",
			Sequence:  1,
			ChannelID: "CHANNEL",
		}
	})

	Context("CreateCommitProposal", func() {
		It("Should work for function CreateCommitProposal", func() {
			controller := gomock.NewController(GinkgoT())
			defer controller.Finish()

			mockSigner := NewMockSigner(controller, "", nil, nil)
			_, err := chaincode.CreateCommitProposal(chaincodeDefinition, mockSigner)
			Expect(err).ShouldNot(HaveOccurred())
		})

		It("when the channel name is not provided", func() {
			errorData := chaincode.Definition{
				Name:      "CHAINCODE",
				Version:   "1.0",
				Sequence:  1,
				ChannelID: "",
			}
			_, err := chaincode.CreateCommitProposal(errorData, nil)
			Expect(err).Should(HaveOccurred())
		})

		It("when the chaincode name is not provided", func() {
			errorData := chaincode.Definition{
				Name:      "",
				Version:   "1.0",
				Sequence:  1,
				ChannelID: "CHANNEL",
			}
			_, err := chaincode.CreateCommitProposal(errorData, nil)
			Expect(err).Should(HaveOccurred())
		})

		It("when the chaincode version is not provided", func() {
			errorData := chaincode.Definition{
				Name:      "CHAINCODE",
				Version:   "",
				Sequence:  1,
				ChannelID: "CHANNEL",
			}
			_, err := chaincode.CreateCommitProposal(errorData, nil)
			Expect(err).Should(HaveOccurred())
		})

		It("when the sequence is not provided", func() {
			errorData := chaincode.Definition{
				Name:      "CHAINCODE",
				Version:   "1.0",
				Sequence:  0,
				ChannelID: "CHANNEL",
			}
			_, err := chaincode.CreateCommitProposal(errorData, nil)
			Expect(err).Should(HaveOccurred())
		})
	})

	Context("Commit", func() {
		It("Should handle Sign error when Commit", func() {
			controller := gomock.NewController(GinkgoT())
			defer controller.Finish()

			mockSigner := NewMockSigningIdentity(controller)
			mockSigner.EXPECT().MspID().Return("").AnyTimes()
			mockSigner.EXPECT().Credentials().Return(nil).AnyTimes()
			mockSigner.EXPECT().Sign(gomock.Any()).Return(nil, fmt.Errorf("tea"))

			err := chaincode.Commit(chaincodeDefinition, mockSigner, endorsementClients, nil)
			Expect(err).Should(HaveOccurred())
		})

		It("Should handle Endorsement error when Commit", func() {
			controller := gomock.NewController(GinkgoT())
			defer controller.Finish()

			mockSigner := NewMockSigner(controller, "", nil, nil)
			mockEndorserClient := &mock.EndorserClient{}
			endorsementClients = make([]peer.EndorserClient, 0)
			endorsementClients = append(endorsementClients, mockEndorserClient)
			mockEndorserClient.ProcessProposalReturns(nil, errors.New("latte"))
			err := chaincode.Commit(chaincodeDefinition, mockSigner, endorsementClients, nil)
			Expect(err).Should(HaveOccurred())
		})

		It("Should handle BroadcastClient error when Commit", func() {
			controller := gomock.NewController(GinkgoT())
			defer controller.Finish()

			mockSigner := NewMockSigner(controller, "", nil, nil)
			mockEndorserClient := &mock.EndorserClient{}
			endorsementClients = make([]peer.EndorserClient, 0)
			endorsementClients = append(endorsementClients, mockEndorserClient)
			mockProposalResponse := &peer.ProposalResponse{
				Response: &peer.Response{
					Status: 200,
				},
				Endorsement: &peer.Endorsement{},
			}
			mockEndorserClient.ProcessProposalReturns(mockProposalResponse, nil)
			mockBroadcastClient := &mock.AtomicBroadcast_BroadcastClient{}
			mockBroadcastClient.SendReturns(errors.New("coffee"))
			err := chaincode.Commit(chaincodeDefinition, mockSigner, endorsementClients, mockBroadcastClient)
			Expect(err).Should(HaveOccurred())
		})

		It("Should works when Commit", func() {
			controller := gomock.NewController(GinkgoT())
			defer controller.Finish()

			mockSigner := NewMockSigner(controller, "", nil, nil)
			mockEndorserClient := &mock.EndorserClient{}
			endorsementClients = make([]peer.EndorserClient, 0)
			endorsementClients = append(endorsementClients, mockEndorserClient)
			mockProposalResponse := &peer.ProposalResponse{
				Response: &peer.Response{
					Status: 200,
				},
				Endorsement: &peer.Endorsement{},
			}
			mockEndorserClient.ProcessProposalReturns(mockProposalResponse, nil)
			mockBroadcastClient := &mock.AtomicBroadcast_BroadcastClient{}
			err := chaincode.Commit(chaincodeDefinition, mockSigner, endorsementClients, mockBroadcastClient)
			Expect(err).ShouldNot(HaveOccurred())
		})
	})
})
