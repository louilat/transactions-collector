package minio

import (
	"fmt"
	"os"
	"transactions-collector/tps"

	"github.com/minio/minio-go"
	"github.com/xitongsys/parquet-go-source/local"
	"github.com/xitongsys/parquet-go/parquet"
	"github.com/xitongsys/parquet-go/writer"
)

func SaveTxRecordsToParquet(endpoint, bucket, key, accessKeyID, secretAccessKey string, rec []tps.Transaction) error {
	useSSL := false
	minioClient, err := minio.New(endpoint, accessKeyID, secretAccessKey, useSSL)
	if err != nil {
		fmt.Println("one")
		panic(err)
	}

	// var b bytes.Buffer
	fw, err := local.NewLocalFileWriter("output.parquet")
	if err != nil {
		return err
	}
	pw, err := writer.NewParquetWriterFromWriter(fw, new(tps.Transaction), int64(len(rec)))
	if err != nil {
		fmt.Println("Could not create parquet writer")
		panic(err)
	}

	pw.CompressionType = parquet.CompressionCodec_GZIP

	for _, d := range rec {
		err = pw.Write(d)
		if err != nil {
			fmt.Println("two")
			return err
		}
	}
	err = pw.WriteStop()
	if err != nil {
		fmt.Println("three")
		return err
	}
	fw.Close()

	// r := bytes.NewReader(&b)
	_, err = minioClient.FPutObject(bucket, key, "output.parquet", minio.PutObjectOptions{})

	if err != nil {
		fmt.Println("four")
		return err
	}

	err = os.Remove("output.parquet")
	if err != nil {
		return err
	}
	return nil
}
