package main

import (
	"context"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/g-hyoga/gRPC-nyumon/downloader"
	"github.com/g-hyoga/gRPC-nyumon/echo"
	"github.com/g-hyoga/gRPC-nyumon/uploader"
)

type echoService struct{}

func (s *echoService) Echo(ctx context.Context, req *echo.EchoRequest) (*echo.EchoResponse, error) {
	return &echo.EchoResponse{Message: req.GetMessage()}, nil
}

type downloadService struct{}

func (s *downloadService) Download(req *downloader.FileRequest, stream downloader.FileService_DownloadServer) error {
	fp := filepath.Join("./resource", req.GetName())

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

	fp := filepath.Join("./resource", name)
	ioutil.WriteFile(fp, blob, 0644)
	stream.SendAndClose(&uploader.FileResponse{Size: int64(len(blob))})

	return nil
}
