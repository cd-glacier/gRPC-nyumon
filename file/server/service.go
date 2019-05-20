package main

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/g-hyoga/gRPC-nyumon/file/downloader"
	"github.com/g-hyoga/gRPC-nyumon/file/uploader"
)

type downloadService struct{}

func (s *downloadService) Download(req *downloader.FileRequest, stream downloader.FileService_DownloadServer) error {
	fp := filepath.Join("./file/resource", req.GetName())

	fs, err := os.Open(fp)
	if err != nil {
		return err
	}
	defer fs.Close()

	buf := make([]byte, 1000*1024)
	for {
		n, err := fs.Read(buf)
		if err != nil {
			if err != io.EOF {
				return err
			}
			break
		}
		err = stream.Send(&downloader.FileResponse{
			Data: buf[:n],
		})
		if err != nil {
			return err
		}
	}

	return nil
}

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

	fp := filepath.Join("./file/resource", name)
	ioutil.WriteFile(fp, blob, 0644)
	stream.SendAndClose(&uploader.FileResponse{Size: int64(len(blob))})

	return nil
}
