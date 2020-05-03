package server

import (
	"context"
	"fmt"
	"github.com/yosef32/go-file-system-server-service/config"
	"log"
	"net"

	fileModel "github.com/yosef32/go-file-system-server-service/model"
	"github.com/yosef32/go-file-system-server-service/mongo"
	files "github.com/yosef32/go-file-system-server-service/proto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc"
)

type Server struct{}

var mongoManager mongo.Controller

func NewServer() *Server {
	config.InitMongo()
	mongoManager = mongo.NewController(config.MongoEndPoint())
	return &Server{}
}

func (s *Server) Serve() {

	config.InitGrpc()

	srv := grpc.NewServer()
	files.RegisterFileServiceServer(srv, s)

	listener, err := net.Listen(config.GrpcNetwork(), config.GrpcPort())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("gRpc run on " + config.GrpcPort())
	if e := srv.Serve(listener); e != nil {
		log.Fatal("gRpc not run: ", err.Error())
	}
}

func (s *Server) Create(ctx context.Context, req *files.CreateReq) (*files.FileRes, error) {
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

	protoFile := &files.File{
		Id:       file.ID.Hex(),
		Owner:    file.Owner,
		Path:     file.Path,
		Name:     file.Name,
		IsFolder: file.IsFolder,
	}

	result := &files.FileRes{
		File: protoFile,
	}

	return result, nil
}

func (s *Server) Read(ctx context.Context, req *files.ReadReq) (*files.FileRes, error) {
	owner := req.GetOwner()
	oid, err := primitive.ObjectIDFromHex(req.GetId())
	if err != nil {
		return nil, fmt.Errorf("InvalidArgument %s", fmt.Sprintf("Could not convert to ObjectId: %v", err))
	}

	file, err := fileModel.Fined(ctx, mongoManager.Client, bson.M{"_id": oid})
	if err != nil {
		return nil, err
	}

	if file.Owner != owner {
		return nil, fmt.Errorf("permittion %s", fmt.Sprintf("permittion denied %s: %v", req.GetOwner(), err))
	}

	protoFile := &files.File{
		Id:       file.ID.Hex(),
		Owner:    file.Owner,
		Path:     file.Path,
		Name:     file.Name,
		IsFolder: file.IsFolder,
	}

	result := &files.FileRes{
		File: protoFile,
	}

	return result, nil
}

func (s *Server) Update(ctx context.Context, req *files.UpdateReq) (*files.FileRes, error) {
	oid, err := primitive.ObjectIDFromHex(req.GetFile().Id)
	if err != nil {
		return nil, fmt.Errorf("InvalidArgument %s", fmt.Sprintf("Could not convert to ObjectId: %v", err))
	}

	owner := req.GetFile().Owner
	name := req.GetFile().Name
	path := req.GetFile().Path
	isFolder := req.GetFile().IsFolder

	newFile := fileModel.New(owner, name, path, isFolder)
	newFile.ID = oid

	updatedFile, err := newFile.Update(ctx, mongoManager.Client, bson.M{"_id": oid})
	if err != nil {
		return nil, err
	}

	protoFile := &files.File{
		Id:       updatedFile.ID.Hex(),
		Owner:    updatedFile.Owner,
		Path:     updatedFile.Path,
		Name:     updatedFile.Name,
		IsFolder: updatedFile.IsFolder,
	}

	result := &files.FileRes{
		File: protoFile,
	}

	return result, nil
}

func (s *Server) Delete(ctx context.Context, req *files.DeleteReq) (*files.DeleteRes, error) {
	owner := req.GetOwner()
	oid, err := primitive.ObjectIDFromHex(req.GetId())
	if err != nil {
		return nil, fmt.Errorf("InvalidArgument %s", fmt.Sprintf("Could not convert to ObjectId: %v", err))
	}

	err = fileModel.Delete(ctx, mongoManager.Client, bson.M{"_id": oid, "owner": owner})
	if err != nil {
		return nil, err
	}

	return &files.DeleteRes{Success: true}, nil
}

func (s *Server) ListFiles(ctx context.Context, req *files.ListFilesReq) (*files.ListFilesRes, error) {

	owner := req.GetOwner()
	filesList, err := fileModel.FindAll(ctx, mongoManager.Client, bson.M{"owner": owner})
	if err != nil {
		return nil, err
	}

	list := []*files.File{}

	for i := 0; i < len(filesList); i++ {
		list = append(list, &files.File{
			Id:       filesList[i].ID.Hex(),
			Owner:    filesList[i].Owner,
			Path:     filesList[i].Path,
			Name:     filesList[i].Name,
			IsFolder: filesList[i].IsFolder,
		})
	}

	return &files.ListFilesRes{Files: list}, nil
}
