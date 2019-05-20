package main

import (
	"context"
	"flag"
	"io"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/g-hyoga/gRPC-nyumon/downloader"
	e "github.com/g-hyoga/gRPC-nyumon/echo"
	"github.com/g-hyoga/gRPC-nyumon/uploader"
	"google.golang.org/grpc"
)

func init() {
	log.SetFlags(0)
	log.SetPrefix("[gRPC] ")
}

func main() {
	target := "localhost:50051"

	mode := flag.String("mode", "download", "echo,downoload,upload")
	message := flag.String("message", "hello~", "hello gRPC world")
	filename := flag.String("filename", "resource.txt", "--mode=resource.txt")
	flag.Parse()

	conn, err := grpc.Dial(target, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %s\n", err)
	}
	defer conn.Close()

	switch *mode {
	case "echo":
		echo(conn, *message)
	case "download":
		download(conn, *filename)
	case "upload":
		upload(conn, *filename)
	}
}

func echo(conn *grpc.ClientConn, message string) {
	client := e.NewEchoServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := client.Echo(ctx, &e.EchoRequest{Message: message})
	if err != nil {
		log.Println(err)
	}

	log.Println(r.GetMessage())
}

func download(conn *grpc.ClientConn, name string) {
	c := downloader.NewFileServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stream, err := c.Download(ctx, &downloader.FileRequest{Name: name})
	if err != nil {
		log.Fatalf("could not downlaod: %s\n", err)
	}

	var blob []byte
	for {
		c, err := stream.Recv()
		if err == io.EOF {
			log.Printf("done %d bytes\n", len(blob))
			break
		}
		if err != nil {
			log.Fatalln(err)
		}
		blob = append(blob, c.GetData()...)
	}

	ioutil.WriteFile(name, blob, 0644)
}

func upload(conn *grpc.ClientConn, name string) {
	c := uploader.NewFileServiceClient(conn)

	fs, err := os.Open(name)
	if err != nil {
		log.Fatalf("could not open file %s\n", err)
	}
	defer fs.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stream, err := c.Upload(ctx)
	if err != nil {
		log.Fatalf("could not upload file: %s\n", err)
	}

	buf := make([]byte, 1000*1024)
	for {
		n, err := fs.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("could not read file %s\n", err)
		}
		stream.Send(&uploader.FileRequest{
			Name: name,
			Data: buf[:n],
		})
	}

	res, err := stream.CloseAndRecv()
	log.Printf("done %d bytes\n", res.GetSize())
}
