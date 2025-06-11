package tps

import "sync"

type Transaction struct {
	BlockNumber int64  `json:"blockNumber" parquet:"name=blockNumber, type=INT64"`
	BlockTime   int64  `json:"blockTime" parquet:"name=blockTime, type=INT64"`
	BlockHash   string `json:"blockHash" parquet:"name=blockHash, type=BYTE_ARRAY, convertedtype=UTF8"`
	BlockNbTx   int    `json:"blockNbTx" parquet:"name=blockNbTx, type=INT64"`
	TxHash      string `json:"txHash" parquet:"name=txHash, type=BYTE_ARRAY, convertedtype=UTF8"`
	From        string `json:"from" parquet:"name=from, type=BYTE_ARRAY, convertedtype=UTF8"`
	To          string `json:"to" parquet:"name=to, type=BYTE_ARRAY, convertedtype=UTF8"`
	TxGas       int64  `json:"txGas" parquet:"name=txGas, type=INT64"`
	TxGasPrice  int64  `json:"txGasPrice" parquet:"name=txGasPrice, type=INT64"`
	TxValue     string `json:"txValue" parquet:"name=txValue, type=BYTE_ARRAY, convertedtype=UTF8"`
	TxNonce     int64  `json:"txNonce" parquet:"name=txNonce, type=INT64"`
	TxData      string `json:"txData" parquet:"name=txData, type=BYTE_ARRAY, convertedtype=UTF8"`
	Error       string `json:"error" parquet:"name=error, type=BYTE_ARRAY, convertedtype=UTF8"`
}

type TransactionRecord struct {
	Records []Transaction
	mu      sync.Mutex
}

func (tr *TransactionRecord) Append(t Transaction) {
	tr.mu.Lock()
	defer tr.mu.Unlock()
	tr.Records = append(tr.Records, t)
}
