package main

import (
	"flag"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	log "github.com/sirupsen/logrus"
	"ioporaclenode/internal/pkg/kyber/pairing/bn256"
	"ioporaclenode/pkg/iop"
	"net"
	"os"
	"os/signal"
)

var (
	addrFlag             = flag.String("address", "127.0.0.1:25565", "server address")
	ethFlag              = flag.String("eth", "ws://127.0.0.1:7545", "eth node address")
	oracleContractFlag   = flag.String("oracleContract", "0x94367f58cC0296196543E378cA81Fec50c929b2C", "oracle contract address")
	registryContractFlag = flag.String("registryContract", "0x045d04c9C372Ea303A6A5F055dead7aF848bCFf6", "registry contract address")
	ecdsaPrivateKeyFlag  = flag.String("ecdsaPrivateKey", "0xe63ff25be694842b3d25f3c8981dbe44b36b23a6effdbe04f9ee11e7965c922b", "private key")
	blsPrivateKeyFlag    = flag.String("blsPrivateKey", "0x2e931ebbc908ec1993a789166f5690ee2ea34830df69a0fd0fc6a456b4aa8a46", "value of the private share")
)

func main() {
	flag.Parse()
	ethClient, err := ethclient.Dial(*ethFlag)
	if err != nil {
		log.Fatalf("dial eth client: %v", err)
	}

	registryContract, err := iop.NewRegistryContract(common.HexToAddress(*registryContractFlag), ethClient)
	if err != nil {
		log.Fatalf("registry contract: %v", err)
	}

	oracleContract, err := iop.NewOracleContract(common.HexToAddress(*oracleContractFlag), ethClient)
	if err != nil {
		log.Fatalf("oracle contract: %v", err)
	}

	ecdsaPrivateKey, err := crypto.HexToECDSA(*ecdsaPrivateKeyFlag)
	if err != nil {
		log.Fatalf("hex to ecdsa: %v", err)
	}

	suite := bn256.NewSuite()
	blsPrivateKey, err := iop.HexToScalar(suite.G2(), *blsPrivateKeyFlag)

	hexAddress, err := iop.AddressFromPrivateKey(ecdsaPrivateKey)
	if err != nil {
		log.Fatalf("address from private key: %v", err)
	}
	account := common.HexToAddress(hexAddress)

	lis, err := net.Listen("tcp", *addrFlag)
	if err != nil {
		log.Fatalf("listen on %s: %v", *addrFlag, err)
	}

	aggregator := iop.NewAggregator(ethClient, registryContract, account)
	txVerifier := iop.NewTransactionVerifier(ethClient)
	oracleNode := iop.NewOracleNode(ethClient, txVerifier, aggregator, oracleContract, registryContract, ecdsaPrivateKey, blsPrivateKey, account, suite)

	go func() {
		if err := oracleNode.Serve(lis); err != nil {
			log.Fatalf("serve %s: %v", lis.Addr(), err)
		}
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	<-sig

	oracleNode.GracefulStop()
	ethClient.Close()
}
