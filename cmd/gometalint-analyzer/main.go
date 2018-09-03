package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/src-d/lookout-gometalint-analyzer"

	"github.com/kelseyhightower/envconfig"
	"github.com/sanity-io/litter"
	"github.com/src-d/lookout"
	"github.com/src-d/lookout/util/grpchelper"
	"google.golang.org/grpc"
	log "gopkg.in/src-d/go-log.v1"
)

var usageMessage = fmt.Sprintf(`usage: %s [-version] [OPTIONS]

%s is a lookout analyzer implementation, based on https://github.com/alecthomas/gometalinter.

`, name, name)

var (
	name        = "gometalint-analyzer"
	version     string
	build       string
	versionFlag = flag.Bool("version", false, "show version")
)

type config struct {
	Host           string `envconfig:"HOST" default:"0.0.0.0"`
	Port           int    `envconfig:"PORT" default:"2001"`
	DataServiceURL string `envconfig:"DATA_SERVICE_URL" default:"ipv4://localhost:10301"`
	LogLevel       string `envconfig:"LOG_LEVEL" default:"info" description:"Logging level (info, debug, warning or error)"`
}

func main() {
	litter.Config.Compact = true
	flag.Usage = func() {
		fmt.Printf(usageMessage)
		flag.PrintDefaults()
	}
	flag.Parse()

	if *versionFlag {
		fmt.Printf("%s %s built on %s\n", name, version, build)
		return
	}

	var conf config
	envconfig.MustProcess("GOMETALINT", &conf)
	log.Infof("Starting %s, %s", name, litter.Sdump(conf))

	grpcAddr, err := grpchelper.ToGoGrpcAddress(conf.DataServiceURL)
	if err != nil {
		log.Errorf(err, "failed to parse DataService addres %s", conf.DataServiceURL)
		return
	}

	conn, err := grpchelper.DialContext(
		context.Background(),
		grpcAddr,
		grpc.WithInsecure(),
		grpc.WithDefaultCallOptions(grpc.FailFast(false)),
	)
	if err != nil {
		log.Errorf(err, "cannot create connection to DataService %s", grpcAddr)
		return
	}

	analyzer := &gometalint.Analyzer{
		Version:    version,
		DataClient: lookout.NewDataClient(conn),
		Args:       append([]string(nil), os.Args[1:]...),
	}

	server := grpchelper.NewServer()
	lookout.RegisterAnalyzerServer(server, analyzer)

	analyzerURL := fmt.Sprintf("ipv4://%s:%d", conf.Host, conf.Port)
	lis, err := grpchelper.Listen(analyzerURL)
	if err != nil {
		log.Errorf(err, "failed to start analyzer gRPC server on %s", analyzerURL)
		return
	}

	log.Infof("server has started on '%s'", analyzerURL)
	err = server.Serve(lis)
	if err != nil {
		log.Errorf(err, "gRPC server failed listening on %v", lis)
	}
	return
}
