/*
Copyright IBM Corp. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package mock

import (
	"context"
	"sync"

	"github.com/hyperledger/fabric-protos-go-apiv2/common"
	"github.com/hyperledger/fabric-protos-go-apiv2/orderer"
	"github.com/hyperledger/fabric-protos-go-apiv2/peer"
	"google.golang.org/grpc"
)

type EndorserClient struct {
	ProcessProposalStub        func(context.Context, *peer.SignedProposal, ...grpc.CallOption) (*peer.ProposalResponse, error)
	processProposalMutex       sync.RWMutex
	processProposalArgsForCall []struct {
		arg1 context.Context
		arg2 *peer.SignedProposal
		arg3 []grpc.CallOption
	}
	processProposalReturns struct {
		result1 *peer.ProposalResponse
		result2 error
	}
	processProposalReturnsOnCall map[int]struct {
		result1 *peer.ProposalResponse
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *EndorserClient) ProcessProposalReturns(result1 *peer.ProposalResponse, result2 error) {
	fake.ProcessProposalStub = nil
	fake.processProposalReturns = struct {
		result1 *peer.ProposalResponse
		result2 error
	}{result1, result2}
}

func (fake *EndorserClient) ProcessProposal(arg1 context.Context, arg2 *peer.SignedProposal, arg3 ...grpc.CallOption) (*peer.ProposalResponse, error) {
	fake.processProposalMutex.Lock()
	ret, specificReturn := fake.processProposalReturnsOnCall[len(fake.processProposalArgsForCall)]
	fake.processProposalArgsForCall = append(fake.processProposalArgsForCall, struct {
		arg1 context.Context
		arg2 *peer.SignedProposal
		arg3 []grpc.CallOption
	}{arg1, arg2, arg3})
	fake.recordInvocation("ProcessProposal", []interface{}{arg1, arg2, arg3})
	fake.processProposalMutex.Unlock()
	if fake.ProcessProposalStub != nil {
		return fake.ProcessProposalStub(arg1, arg2, arg3...)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.processProposalReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *EndorserClient) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}

type AtomicBroadcast_BroadcastClient struct {
	error_whenSend error
	grpc.ClientStream
}

func (client *AtomicBroadcast_BroadcastClient) Send(*common.Envelope) error {
	return client.error_whenSend
}
func (client *AtomicBroadcast_BroadcastClient) Recv() (*orderer.BroadcastResponse, error) {
	return nil, nil
}

func (client *AtomicBroadcast_BroadcastClient) SendReturns(err error) {
	client.error_whenSend = err
}
