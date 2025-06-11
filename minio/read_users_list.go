package minio

import (
	"encoding/json"
	"io"
	"transactions-collector/tps"

	"github.com/minio/minio-go"
)

func ReadUsersList(accessKeyID, secretAccessKey string) (tps.UsersList, error) {
	useSSL := false
	endpoint := "minio-simple.lab.groupe-genes.fr"

	minioClient, err := minio.New(endpoint, accessKeyID, secretAccessKey, useSSL)
	if err != nil {
		return make([]tps.User, 0), err
	}

	usr, err := minioClient.GetObject("projet-datalab-group-jprat", "transactions-datasource/borrowers/borrowers_list.json", minio.GetObjectOptions{})
	if err != nil {
		return make([]tps.User, 0), err
	}

	usersBytes, err := io.ReadAll(usr)
	if err != nil {
		return make([]tps.User, 0), err
	}

	users := make([]tps.User, 0)
	err = json.Unmarshal(usersBytes, &users)
	if err != nil {
		return make([]tps.User, 0), err
	}
	return users, nil
}
