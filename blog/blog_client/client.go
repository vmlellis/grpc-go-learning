package main

import (
	"context"
	"fmt"
	"log"

	"github.com/vmlellis/grpc-go-learning/blog/blogpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func main() {
	fmt.Println("Blog Client Started")

	tls := false
	opts := grpc.WithInsecure()

	if tls {
		certFile := "ssl/ca.crt" // Certificate Authority Trust certificate
		creds, sslErr := credentials.NewClientTLSFromFile(certFile, "")
		if sslErr != nil {
			log.Fatalf("Failed when loading CA trust certificate: %v", sslErr)
			return
		}

		opts = grpc.WithTransportCredentials(creds)
	}

	cc, err := grpc.Dial("localhost:50051", opts)
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer cc.Close()

	c := blogpb.NewBlogServiceClient(cc)

	// create Blog
	fmt.Println("Creating the blog")
	blog := &blogpb.Blog{
		AuthorId: "Victor",
		Title:    "My First Blog",
		Content:  "Content of the first blog",
	}
	createBlogRes, err := c.CreateBlog(context.Background(), &blogpb.CreateBlogRequest{Blog: blog})
	if err != nil {
		log.Fatalf("Unexpected error: %v", err)
	}
	fmt.Printf("Blog has been created: %v\n", createBlogRes)
	blogId := createBlogRes.Blog.GetId()

	// read Blog
	fmt.Println("Reading the blog")

	_, err = c.ReadBlog(context.Background(), &blogpb.ReadBlogRequest{BlogId: "213aad"})
	if err != nil {
		fmt.Printf("Error happened while reading: %v\n", err)
	}

	readBlogReq := &blogpb.ReadBlogRequest{BlogId: blogId}
	readBloqRes, err := c.ReadBlog(context.Background(), readBlogReq)
	if err != nil {
		fmt.Printf("Error happened while reading: %v\n", err)
	} else {
		fmt.Printf("Blog was read: %v\n", readBloqRes)
	}
}
