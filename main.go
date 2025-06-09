package main

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"os"
	"time"
	"transactions-collector/block"
	"transactions-collector/minio"
	"transactions-collector/tps"

	"transactions-collector/etlutils"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	provider := os.Getenv("PROVIDER")
	endpoint := "minio-simple.lab.groupe-genes.fr"
	bucket := "projet-datalab-group-jprat"
	accessKeyId := os.Getenv("ACCESS_KEY_ID")
	secretAccessKey := os.Getenv("SECRET_ACCESS_KEY")
	start := os.Getenv("START_DATE")
	stop := os.Getenv("STOP_DATE")

	fmt.Println("Starting Job...")

	// Parse dates
	startDate, err := time.Parse("2006-01-02", start)
	if err != nil {
		panic(err)
	}
	stopDate, err := time.Parse("2006-01-02", stop)
	if err != nil {
		panic(err)
	}

	// Create povider client
	client, err := ethclient.Dial(provider)
	if err != nil {
		panic(err)
	}

	references := etlutils.GetBlockReferences()

	var day time.Time
	for day = startDate; day.Before(stopDate); day = day.AddDate(0, 0, 1) {
		ref := references[time.Date(day.Year(), day.Month(), 1, 0, 0, 0, 0, time.UTC).Unix()]

		// Query transactions
		transactions, err := QueryTxPerDay(client, day, ref, 7000)
		if err != nil {
			panic(err)
		}

		// Save transactions to minio
		key := "transactions-datasource/raw-transactions/transactions_snapshot_date=" + fmt.Sprint(day)[:10] + "/raw_transactions.json"
		err = minio.SaveTxRecordsToParquet(endpoint, bucket, key, accessKeyId, secretAccessKey, transactions)
		if err != nil {
			panic(err)
		}
	}
}

func QueryTxPerDay(backend *ethclient.Client, day time.Time, blockref *big.Int, finderstep int64) ([]tps.Transaction, error) {
	start := day.Unix()
	end := day.AddDate(0, 0, 1).Unix()

	fmt.Printf("*****Extracting price for day: %v*****\n", day)

	fmt.Printf("   Start timestamp: %v\n", start)
	fmt.Printf("   End timestamp:   %v\n", end)

	startblock, _, err := block.FindClosestBlocks(backend, uint64(start), blockref, finderstep)
	if err != nil {
		return make([]tps.Transaction, 0), err
	}

	_, endblock, err := block.FindClosestBlocks(backend, uint64(end), blockref, finderstep)
	if err != nil {
		return make([]tps.Transaction, 0), err
	}

	blockNumber := new(big.Int)
	txRecords := make([]tps.Transaction, 0)
	for blockNumber.Set(startblock); blockNumber.Cmp(endblock) <= 0; blockNumber.Add(blockNumber, big.NewInt(1)) {
		// Query the block corresponding to block number
		fmt.Printf("Treating block %v / %v\n", blockNumber, endblock)
		block, err := backend.BlockByNumber(context.Background(), blockNumber)
		if err != nil {
			return make([]tps.Transaction, 0), err
		}

		// Add each transaction of the block
		for _, tx := range block.Transactions() {
			var signer types.Signer
			switch {
			case tx.Type() == types.AccessListTxType:
				signer = types.NewEIP2930Signer(tx.ChainId())
			case tx.Type() == types.DynamicFeeTxType:
				signer = types.NewLondonSigner(tx.ChainId())
			case tx.Type() == types.BlobTxType:
				signer = types.NewCancunSigner(tx.ChainId())
			case tx.Type() == types.SetCodeTxType:
				signer = types.NewPragueSigner(tx.ChainId())
			default:
				signer = types.NewEIP155Signer(tx.ChainId())
			}
			sender, err := types.Sender(signer, tx)
			if err != nil {
				return make([]tps.Transaction, 0), err
			}

			receiver := ""
			if tx.To() != nil {
				receiver = tx.To().Hex()
			}

			txRecords = append(txRecords, tps.Transaction{
				BlockNumber: block.Number().Int64(),
				BlockTime:   int64(block.Time()),
				BlockHash:   block.Hash().Hex(),
				BlockNbTx:   len(block.Transactions()),
				TxHash:      tx.Hash().Hex(),
				From:        sender.Hex(),
				To:          receiver,
				TxGas:       int64(tx.Gas()),
				TxGasPrice:  tx.GasPrice().Int64(),
				TxValue:     tx.Value().String(),
				TxNonce:     int64(tx.Nonce()),
				TxData:      hex.EncodeToString(tx.Data()),
			})
		}
	}
	return txRecords, nil
}
