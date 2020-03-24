package server

import (
	"fmt"
	"net"
	"context"
	"strings"
	"romano.com/proto"
	"romano.com/mongo"
	"romano.com/model/fileModel"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Server struct {
	port string
}

const (
	mongoID string = "_id";
)

var mongoManager mongo.Controller;

func New(port string) Server {
	var s Server;
	s.port = port;
	mongoManager = mongo.NewController("mongodb://mongo:27017");
  	return s;
}

func (s *Server) CreateListener() {

	listener, err := net.Listen("tcp", s.port);
	if err != nil {
		panic(err)
	}

	srv := grpc.NewServer();
	proto.RegisterAddServiceServer(srv, s);
	reflection.Register(srv);

	fmt.Println("GRpc run on" + s.port);
	if e := srv.Serve(listener); e != nil {
		fmt.Println("GRpc not run: ", err.Error());
	}
}

func (s *Server) Create(ctx context.Context, req *proto.CreateReq) (*proto.FileRes, error) {
	owner, name, path, isFolder := req.GetFile().Owner, req.GetFile().Name, req.GetFile().Path, req.GetFile().IsFolder;
	file := fileModel.New(owner, name, path, isFolder);
	decoded, err := file.Insert(mongoManager.Client)

	fmt.Println(decoded.ID.Hex());

	if err != nil {
		return nil, err
	}

	return &proto.FileRes{
		File: &proto.File{
			Id:       decoded.ID.Hex(),
			Owner:    decoded.Owner,
			Path:     decoded.Path,
			Name:     decoded.Name,
			IsFolder: decoded.IsFolder,
		},
	}, nil
}

func (s *Server) Read(ctx context.Context, req *proto.ReadReq) (*proto.FileRes, error) {
	owner := req.GetOwner();
	oid, err := primitive.ObjectIDFromHex(req.GetId())
	if err != nil {
		return nil, fmt.Errorf("InvalidArgument", fmt.Sprintf("Could not convert to ObjectId: %v", err))
	}

	decoded, err := fileModel.FinedOne(mongoManager.Client, bson.M {mongoID : oid} )
	
	if strings.Compare(decoded.Owner, owner) == 0 {
		return nil, fmt.Errorf("Permition", fmt.Sprintf("Permition denaid %s: %v", req.GetOwner(), err));
	}

	return &proto.FileRes{
		File: &proto.File{
			Id:       decoded.ID.Hex(),
			Owner:    decoded.Owner,
			Path:     decoded.Path,
			Name:     decoded.Name,
			IsFolder: decoded.IsFolder,
		},
	}, nil
}


func (s *Server) Update(ctx context.Context, req *proto.UpdateReq) (*proto.FileRes, error) {
	owner := req.GetOwner();
	oid, err := primitive.ObjectIDFromHex(req.GetFile().Id)
	if err != nil {
		return nil, fmt.Errorf("InvalidArgument", fmt.Sprintf("Could not convert to ObjectId: %v", err))
	}

	owner, name, path, isFolder := req.GetFile().Owner, req.GetFile().Name, req.GetFile().Path, req.GetFile().IsFolder;
	file := fileModel.New(owner, name, path, isFolder);

	decoded, err := file.Update(mongoManager.Client, bson.M{mongoID : oid})
	if err != nil {
		return nil, err;
	}

	return &proto.FileRes{
		File: &proto.File{
			Id:       decoded.ID.Hex(),
			Owner:    decoded.Owner,
			Path:     decoded.Path,
			Name:     decoded.Name,
			IsFolder: decoded.IsFolder,
		},
	}, nil
}

func (s *Server) Delete(ctx context.Context, req *proto.DeleteReq) (*proto.BoolRes, error) {
	owner := req.GetOwner();
	oid, err := primitive.ObjectIDFromHex(req.GetId())
	if err != nil {
		return nil, fmt.Errorf("InvalidArgument", fmt.Sprintf("Could not convert to ObjectId: %v", err))
	}

	decoded, err := fileModel.FinedOne(mongoManager.Client, bson.M {mongoID : oid} )
	
	if strings.Compare(decoded.Owner, owner) == 1 {
		return nil, fmt.Errorf("Permition", fmt.Sprintf("Permition denaide %s: %v", req.GetOwner(), err));
	}

	return &proto.BoolRes{Success: true}, nil;
}

func (s *Server) ListFiles(ctx context.Context, req *proto.ListFilesReq) (*proto.ListFilesRes, error) {
	
	owner := req.GetOwner();
	decoded, err := fileModel.FinedAll(mongoManager.Client, bson.M {});
	if err != nil {
		return nil, err;
	}

	list := []*proto.File{};

	for i := 0; i < len(decoded); i++ {

		if strings.Compare(decoded[i].Owner, owner) == 0 {
			list = append(list, &proto.File{
				Id:       decoded[i].ID.Hex(),
				Owner:    decoded[i].Owner,
				Path:     decoded[i].Path,
				Name:     decoded[i].Name,
				IsFolder: decoded[i].IsFolder,
			});
		}

	}

	fmt.Println(list);
	return &proto.ListFilesRes{Files: list}, nil;
}


