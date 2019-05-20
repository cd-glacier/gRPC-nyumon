package main

import (
	"io"
	"io/ioutil"
	"log"
	"path/filepath"

	"github.com/g-hyoga/gRPC-nyumon/uploader"
)

type uploadService struct{}

func (s *uploadService) Upload(stream uploader.FileService_UploadServer) error {
	var blob []byte
	var name string

	for {
		c, err := stream.Recv()
		if err == io.EOF {
			log.Printf("done %d bytes\n", len(blob))
			break
		}
		if err != nil {
			panic(err)
		}

		name = c.GetName()
		blob = append(blob, c.GetData()...)
	}

	fp := filepath.Join("./resource", name)
	ioutil.WriteFile(fp, blob, 0644)
	stream.SendAndClose(&uploader.FileResponse{Size: int64(len(blob))})

	return nil
}
