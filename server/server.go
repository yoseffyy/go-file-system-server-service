package server

import (
	"context"
	"fmt"
	"log"
	"net"

	fileModel "file-service.com/model"
	"file-service.com/mongo"
	"file-service.com/proto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc"
	// "google.golang.org/grpc/reflection"
)

type Server struct {
	port string
}

const (
	mongoID string = "_id"
)

var mongoManager mongo.Controller

func New(port string) *Server {
	s := Server{}
	s.port = port
	mongoManager = mongo.NewController("mongodb://mongo:27017")
	return &s
}

func (s *Server) Listen() {
	srv := grpc.NewServer()
	proto.RegisterAddServiceServer(srv, s)
	// reflection.Register(srv)

	listener, err := net.Listen("tcp", s.port)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("GRpc run on" + s.port)
	if e := srv.Serve(listener); e != nil {
		log.Fatal("GRpc not run: ", err.Error())
	}
}

func (s *Server) Create(ctx context.Context, req *proto.CreateReq) (*proto.FileRes, error) {
	if req.GetFile() == nil {
		log.Fatal("The file is empty")
	}

	owner := req.GetFile().Owner
	name := req.GetFile().Name
	path := req.GetFile().Path
	isFolder := req.GetFile().IsFolder

	newFile := fileModel.New(owner, name, path, isFolder)
	file, err := newFile.Insert(ctx, mongoManager.Client)

	fmt.Println(file.ID.Hex())

	if err != nil {
		return nil, err
	}

	protoFile := &proto.File{
		Id:       file.ID.Hex(),
		Owner:    file.Owner,
		Path:     file.Path,
		Name:     file.Name,
		IsFolder: file.IsFolder,
	}

	result := &proto.FileRes{
		File: protoFile,
	}

	return result, nil
}

func (s *Server) Read(ctx context.Context, req *proto.ReadReq) (*proto.FileRes, error) {
	owner := req.GetOwner()
	oid, err := primitive.ObjectIDFromHex(req.GetId())
	if err != nil {
		return nil, fmt.Errorf("InvalidArgument %s", fmt.Sprintf("Could not convert to ObjectId: %v", err))
	}

	file, err := fileModel.Fined(ctx, mongoManager.Client, bson.M{mongoID: oid})
	if err != nil {
		return nil, err
	}

	if file.Owner != owner {
		return nil, fmt.Errorf("permittion %s", fmt.Sprintf("permittion denied %s: %v", req.GetOwner(), err))
	}

	protoFile := &proto.File{
		Id:       file.ID.Hex(),
		Owner:    file.Owner,
		Path:     file.Path,
		Name:     file.Name,
		IsFolder: file.IsFolder,
	}

	result := &proto.FileRes{
		File: protoFile,
	}

	return result, nil
}

func (s *Server) Update(ctx context.Context, req *proto.UpdateReq) (*proto.FileRes, error) {
	oid, err := primitive.ObjectIDFromHex(req.GetFile().Id)
	if err != nil {
		return nil, fmt.Errorf("InvalidArgument %s", fmt.Sprintf("Could not convert to ObjectId: %v", err))
	}

	owner := req.GetFile().Owner
	name := req.GetFile().Name
	path := req.GetFile().Path
	isFolder := req.GetFile().IsFolder

	newFile := fileModel.New(owner, name, path, isFolder)

	updatedFile, err := newFile.Update(ctx, mongoManager.Client, bson.M{mongoID: oid})
	if err != nil {
		return nil, err
	}

	protoFile := &proto.File{
		Id:       updatedFile.ID.Hex(),
		Owner:    updatedFile.Owner,
		Path:     updatedFile.Path,
		Name:     updatedFile.Name,
		IsFolder: updatedFile.IsFolder,
	}

	result := &proto.FileRes{
		File: protoFile,
	}

	return result, nil
}

func (s *Server) Delete(ctx context.Context, req *proto.DeleteReq) (*proto.BoolRes, error) {
	owner := req.GetOwner()
	oid, err := primitive.ObjectIDFromHex(req.GetId())
	if err != nil {
		return nil, fmt.Errorf("InvalidArgument %s", fmt.Sprintf("Could not convert to ObjectId: %v", err))
	}

	file, err := fileModel.Fined(ctx, mongoManager.Client, bson.M{mongoID: oid})
	if err != nil {
		return nil, err
	}

	if file.Owner != owner {
		return nil, fmt.Errorf("permittion %s", fmt.Sprintf("permittion denied %s: %v", req.GetOwner(), err))
	}

	return &proto.BoolRes{Success: true}, nil
}

func (s *Server) ListFiles(ctx context.Context, req *proto.ListFilesReq) (*proto.ListFilesRes, error) {

	owner := req.GetOwner()
	files, err := fileModel.FinedAll(ctx, mongoManager.Client, bson.M{"owner": owner})
	if err != nil {
		return nil, err
	}

	list := []*proto.File{}

	for i := 0; i < len(files); i++ {
		list = append(list, &proto.File{
			Id:       files[i].ID.Hex(),
			Owner:    files[i].Owner,
			Path:     files[i].Path,
			Name:     files[i].Name,
			IsFolder: files[i].IsFolder,
		})
	}

	return &proto.ListFilesRes{Files: list}, nil
}
