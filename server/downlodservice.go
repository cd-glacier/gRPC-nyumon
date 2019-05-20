package main

import (
	"io"
	"os"
	"path/filepath"

	"github.com/g-hyoga/gRPC-nyumon/downloader"
)

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
