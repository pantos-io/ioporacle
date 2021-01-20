package iop

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	log "github.com/sirupsen/logrus"
	"go.dedis.ch/kyber/v3"
	"go.dedis.ch/kyber/v3/pairing"
	dkg "go.dedis.ch/kyber/v3/share/dkg/pedersen"
	vss "go.dedis.ch/kyber/v3/share/vss/pedersen"
	"go.dedis.ch/kyber/v3/sign/tbls"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"ioporaclenode/internal/pkg/kyber/pairing/bn256"
	"math/big"
	"net"
	"time"
)

type oracleNode struct {
	UnimplementedOracleNodeServer
	server           *grpc.Server
	txVerifier       TransactionVerifier
	aggregator       Aggregator
	ethClient        *ethclient.Client
	oracleContract   *OracleContract
	registryContract *RegistryContractWrapper
	dkg              *dkg.DistKeyGenerator
	ecdsaPrivateKey  *ecdsa.PrivateKey
	blsPrivateKey    kyber.Scalar
	blsPublicKey     kyber.Point
	account          common.Address
	suite            pairing.Suite
	distKey          *dkg.DistKeyShare
}

func NewOracleNode(
	ethClient *ethclient.Client,
	txVerifier TransactionVerifier,
	aggregator Aggregator,
	oracleContract *OracleContract,
	registryContract *RegistryContractWrapper,
	ecdsaPrivateKey *ecdsa.PrivateKey,
	blsPrivateKey kyber.Scalar,
	account common.Address,
	suite pairing.Suite,
) *oracleNode {
	grpcServer := grpc.NewServer()
	node := &oracleNode{
		server:           grpcServer,
		ethClient:        ethClient,
		txVerifier:       txVerifier,
		aggregator:       aggregator,
		oracleContract:   oracleContract,
		registryContract: registryContract,
		ecdsaPrivateKey:  ecdsaPrivateKey,
		blsPrivateKey:    blsPrivateKey,
		blsPublicKey:     suite.G2().Point().Mul(blsPrivateKey, nil),
		account:          account,
		suite:            suite,
	}
	RegisterOracleNodeServer(grpcServer, node)
	return node
}

func (n *oracleNode) Serve(lis net.Listener) error {
	go func() {
		err := n.watchDistributedKeyGenerationLog(context.Background())
		if err != nil {
			log.Errorf("watch distributed key generation log: %v", err)
		}
	}()
	go func() {
		err := n.watchVerifyTransactionLog(context.Background())
		if err != nil {
			log.Errorf("watch verify transaction log: %v", err)
		}
	}()
	err := n.register(lis.Addr().String())
	if err != nil {
		return fmt.Errorf("register: %v", err)
	}
	return n.server.Serve(lis)
}

func (n *oracleNode) watchDistributedKeyGenerationLog(ctx context.Context) error {
	sink := make(chan *RegistryContractDistributedKeyGenerationLog)
	defer close(sink)

	sub, err := n.registryContract.WatchDistributedKeyGenerationLog(
		&bind.WatchOpts{
			Context: context.Background(),
		},
		sink,
	)
	if err != nil {
		return err
	}
	defer sub.Unsubscribe()

	for {
		select {
		case event := <-sink:
			err = n.handleDistributedKeyGenerationLog(ctx, event)
			if err != nil {
				log.Errorf("handle distributed key generation log: %v", err)
			}
		case err = <-sub.Err():
			return err
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (n *oracleNode) handleDistributedKeyGenerationLog(ctx context.Context, event *RegistryContractDistributedKeyGenerationLog) error {
	pubKeys, err := n.registryContract.FindPubKeys()
	if err != nil {
		return fmt.Errorf("find public keys: %w", err)
	}

	distKeyGenerator, err := dkg.NewDistKeyGenerator(bn256.NewSuiteG2(), n.blsPrivateKey, pubKeys, int(event.Threshold.Int64()))
	if err != nil {
		return fmt.Errorf("dkg: %w", err)
	}
	n.dkg = distKeyGenerator

	err = n.sendDeals(ctx)
	if err != nil {
		return fmt.Errorf("send deals: %w", err)
	}

	timeout := time.After(30 * time.Second)

loop:
	for {
		select {
		case <-timeout:
			n.dkg.SetTimeout()
			break loop
		default:
			if n.dkg.Certified() {
				break loop
			}
			time.Sleep(50 * time.Millisecond)
		}
	}

	log.Infof("qualified shares: %v\n", n.dkg.QualifiedShares())
	log.Infof("QUAL: %v\n", n.dkg.QUAL())

	distrKey, err := n.dkg.DistKeyShare()
	if err != nil {
		return fmt.Errorf("distributed key share: %w", err)
	}

	n.distKey = distrKey

	log.Infof("Public Key: %v", n.distKey.Public())

	return nil
}

func (n *oracleNode) initOtherNodes() (map[int]OracleNodeClient, error) {
	nodes, err := n.registryContract.FindIopNodes()
	if err != nil {
		return nil, fmt.Errorf("find nodes: %w", err)
	}

	otherNodes := make(map[int]OracleNodeClient)
	for i := 0; i < len(nodes); i++ {
		if nodes[i].Addr == n.account {
			continue
		}
		conn, err := grpc.Dial(nodes[i].IpAddr, grpc.WithInsecure())
		if err != nil {
			log.Errorf("dial %s: %v", nodes[i].IpAddr, err)
			continue
		}
		log.Printf("Init node %s", nodes[i].Addr.String())
		otherNodes[i] = NewOracleNodeClient(conn)
	}
	return otherNodes, nil
}

func (n *oracleNode) sendDeals(ctx context.Context) error {
	deals, err := n.dkg.Deals()
	if err != nil {
		return fmt.Errorf("deals: %w", err)
	}

	otherNodes, err := n.initOtherNodes()
	if err != nil {
		return fmt.Errorf("init other nodes: %w", err)
	}

	for i, deal := range deals {
		ctxTimeout, cancel := context.WithTimeout(ctx, 10*time.Second)
		request := &ProcessDealRequest{
			Deal: dealToPb(deal),
		}
		log.Infof("Sending deal to node %d", i)
		_, err := otherNodes[i].ProcessDeal(ctxTimeout, request)
		if err != nil {
			log.Errorf("process deal: %v", err)
		}
		cancel()
	}
	return nil
}

func (n *oracleNode) ProcessDeal(ctx context.Context, request *ProcessDealRequest) (*ProcessDealResponse, error) {
	log.Infof("Process deal from node %d", request.Deal.Index)

	response, err := n.dkg.ProcessDeal(pbToDeal(request.Deal))
	if err != nil {
		return nil, fmt.Errorf("process deal: %w", err)
	}

	err = n.sendResponse(context.Background(), response)
	if err != nil {
		log.Errorf("send response: %v", err)
	}

	return &ProcessDealResponse{
		Response: responseToPb(response),
	}, nil
}

func (n *oracleNode) sendResponse(ctx context.Context, response *dkg.Response) error {
	otherNodes, err := n.initOtherNodes()
	if err != nil {
		return fmt.Errorf("init other nodes: %w", err)
	}
	for i, otherNode := range otherNodes {
		otherNode := otherNode
		i := i
		go func() {
			ctxTimeout, cancel := context.WithTimeout(ctx, 10*time.Second)
			request := &ProcessResponseRequest{
				Response: responseToPb(response),
			}
			log.Infof("Sending response to dealer %d with verifier %d to other node %d", response.Index, response.Response.Index, i)
			_, err := otherNode.ProcessResponse(ctxTimeout, request)
			if err != nil {
				log.Errorf("process response: %v", err)
			}
			log.Infof("Sent response to dealer %d with verifier %d to other node %d", response.Index, response.Response.Index, i)
			cancel()
		}()
	}

	return nil
}

func (n *oracleNode) ProcessResponse(ctx context.Context, request *ProcessResponseRequest) (*ProcessResponseResponse, error) {
	log.Infof("Process response to dealer %d with verifier %d", request.Response.Index, request.Response.Response.Index)

	_, err := n.dkg.ProcessResponse(pbToResponse(request.Response))
	for errors.Is(err, vss.ErrNoDealBeforeResponse) {
		if err != nil {
			log.Infof("No deal response to dealer %d with verifier %d. Retry later.", request.Response.Index, request.Response.Response.Index)
		}
		_, err = n.dkg.ProcessResponse(pbToResponse(request.Response))
		time.Sleep(1 * time.Second)
	}
	log.Infof("PROCESSED response to dealer %d with verifier %d", request.Response.Index, request.Response.Response.Index)

	return &ProcessResponseResponse{
		Justification: nil,
	}, nil
}

func (n *oracleNode) watchVerifyTransactionLog(ctx context.Context) error {
	sink := make(chan *OracleContractVerifyTransactionLog)
	defer close(sink)

	sub, err := n.oracleContract.WatchVerifyTransactionLog(
		&bind.WatchOpts{
			Context: context.Background(),
		},
		sink,
		nil,
		nil,
	)
	if err != nil {
		return err
	}
	defer sub.Unsubscribe()

	for {
		select {
		case event := <-sink:
			isLeader, err := n.registryContract.IsLeader(nil, n.account)
			if err != nil {
				log.Errorf("is leader: %v", err)
				continue
			}
			if !isLeader {
				continue
			}
			err = n.handleVerifyTransactionLog(ctx, event)
			if err != nil {
				log.Errorf("handle verify transaction log: %v", err)
			}
		case err = <-sub.Err():
			return err
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (n *oracleNode) handleVerifyTransactionLog(ctx context.Context, event *OracleContractVerifyTransactionLog) error {
	result, _, err := n.aggregator.Aggregate(ctx, event.Id, common.BytesToHash(event.Hash[:]), event.Confirmations.Uint64())
	if err != nil {
		return fmt.Errorf("verify transaction remote: %w", err)
	}
	auth := bind.NewKeyedTransactor(n.ecdsaPrivateKey)
	_, err = n.oracleContract.SubmitVerification(auth, event.Id, result, [2]*big.Int{big.NewInt(0), big.NewInt(0)})
	if err != nil {
		return fmt.Errorf("verify transaction result: %v", err)
	}
	return nil
}

func (n *oracleNode) VerifyTransaction(ctx context.Context, request *VerifyTransactionRequest) (*VerifyTransactionResponse, error) {
	result, err := n.txVerifier.VerifyTransaction(ctx, common.HexToHash(request.Tx), request.Confirmations)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "verify transaction: %v", err)
	}

	uint256Ty, _ := abi.NewType("uint256", "", nil)
	boolTy, _ := abi.NewType("bool", "", nil)

	arguments := abi.Arguments{
		{
			Type: uint256Ty,
		},
		{
			Type: boolTy,
		},
	}

	message, _ := arguments.Pack(
		big.NewInt(request.Id),
		result,
	)

	sig, err := tbls.Sign(n.suite, n.distKey.PriShare(), message)
	if err != nil {
		return nil, fmt.Errorf("tbls sign: %v", err)
	}

	return &VerifyTransactionResponse{
		Id:        request.Id,
		Result:    result,
		Signature: sig,
	}, nil
}

func (n *oracleNode) register(ipAddr string) error {
	isRegistered, err := n.registryContract.OracleNodeIsRegistered(nil, n.account)
	if err != nil {
		return fmt.Errorf("is registered: %w", err)
	}

	b, err := n.blsPublicKey.MarshalBinary()
	if err != nil {
		return fmt.Errorf("marshal bls public key: %v", err)
	}

	auth := bind.NewKeyedTransactor(n.ecdsaPrivateKey)
	if !isRegistered {
		_, err = n.registryContract.RegisterOracleNode(auth, ipAddr, b)
		if err != nil {
			return fmt.Errorf("register iop node: %w", err)
		}
	}
	return nil
}

func (n *oracleNode) Stop() {
	n.server.Stop()
}

func (n *oracleNode) GracefulStop() {
	n.server.GracefulStop()
}
