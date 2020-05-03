package config

import (
	"github.com/spf13/viper"
)

const (
	config_db_name         = "db_name"
	config_collection_name = "collection_name"
	config_db_host         = "db_host"
	config_db_port         = "db_port"
	config_db_end_point    = "db_end_point"

	config_grpc_network    = "grpc_network"
	config_grpc_port    = "grpc_port"
)

func InitMongo() {

	viper.SetDefault(config_db_name, "file_service")
	viper.SetDefault(config_collection_name, "file")
	viper.SetDefault(config_db_host, "mongodb://localhost")
	viper.SetDefault(config_db_port, "27017")
	viper.AutomaticEnv()
	viper.SetDefault(config_db_end_point, viper.GetString(config_db_host)+":"+viper.GetString(config_db_port)+"/"+viper.GetString(config_db_name))
	viper.AutomaticEnv()

}

func MongoEndPoint() string {
	return viper.GetString(config_db_end_point)
}

func DBName() string {
	return viper.GetString(config_db_name)
}

func CollectionName() string {
	return viper.GetString(config_collection_name)
}

func InitGrpc()  {
	viper.SetDefault(config_grpc_network, "tcp")
	viper.SetDefault(config_grpc_port, ":4040")
	viper.AutomaticEnv()
}

func GrpcNetwork() string  {
	return viper.GetString(config_grpc_network)
}

func GrpcPort() string  {
	return viper.GetString(config_grpc_port)
}