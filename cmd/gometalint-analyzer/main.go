package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/bzz/lookout-gometalint-analyzer"
	//TODO: extract to golang sdk
	"github.com/bzz/lookout-gometalint-analyzer/util/grpchelper"

	"github.com/kelseyhightower/envconfig"
	"github.com/sanity-io/litter"
	"github.com/src-d/lookout"
	"google.golang.org/grpc"
	log "gopkg.in/src-d/go-log.v1"
)

var usageMessage = fmt.Sprintf(`usage: %s [-version] [OPTIONS]

%s is a lookout analyzer implementation, based on https://github.com/alecthomas/gometalinter.

OPTIONS - any of the supported by gometalinter.
`, name, name)

var (
	name        = "gometalint-analyzer"
	version     string
	build       string
	versionFlag = flag.Bool("version", false, "show version")
)

type config struct {
	Host          string `envconfig:"HOST" default:"0.0.0.0"`
	Port          int    `envconfig:"PORT" default:"2001"`
	DataServerURL string `envconfig:"DATA_SERVER_URL" default:"ipv4://localhost:10301"`
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

	conn, err := grpchelper.DialContext(
		context.Background(),
		conf.DataServerURL,
		grpc.WithInsecure(),
		grpc.WithDefaultCallOptions(grpc.FailFast(false)),
	)
	if err != nil {
		log.Errorf(err, "cannot create connection to DataServer %s", conf.DataServerURL)
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
