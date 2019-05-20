package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"time"

	cpb "github.com/g-hyoga/gRPC-nyumon/chat"
	"github.com/g-hyoga/gRPC-nyumon/downloader"
	epb "github.com/g-hyoga/gRPC-nyumon/echo"
	"github.com/g-hyoga/gRPC-nyumon/uploader"
	"google.golang.org/grpc"
)

func init() {
	log.SetFlags(0)
	log.SetPrefix("[gRPC client] ")
}

func main() {
	target := "localhost:50051"

	mode := flag.String("mode", "download", "echo,downoload,upload")
	message := flag.String("message", "hello~", "hello gRPC world")
	name := flag.String("name", "resource.txt", "--mode=resource.txt")
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
		download(conn, *name)
	case "upload":
		upload(conn, *name)
	case "chat":
		chat(conn, *name)
	}
}

func echo(conn *grpc.ClientConn, message string) {
	client := epb.NewEchoServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := client.Echo(ctx, &epb.EchoRequest{Message: message})
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

func chat(conn *grpc.ClientConn, name string) {
	c := cpb.NewChatServiceClient(conn)

	stream, err := c.Connect(context.Background())
	if err != nil {
		log.Fatalf("could not connect: %s", err)
	}

	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				return
			}
			if err != nil {
				log.Fatalf("Failed to recv: %v", err)
			}
			fmt.Println(fmt.Sprintf("[%s] %s",
				res.GetName(), res.GetMessage()))
		}
	}()

	for {
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		msg := scanner.Text()
		if msg == ":quit" {
			stream.CloseSend()
			return
		}
		stream.Send(&cpb.Post{Name: name, Message: msg})
	}
}
