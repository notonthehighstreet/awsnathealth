package main

import (
	"aws_nat/awsapitools"
	"aws_nat/errhandling"
	"aws_nat/hostping"
	"aws_nat/httptools"
	"aws_nat/logging"
	"io/ioutil"
	"os"
	"time"

	"github.com/BurntSushi/toml"
)

type natConfig struct {
	MyInstanceID    string
	OtherInstanceID string
	OtherInstanceIP string
	HTTPPort        string
	VpcID           string
	AwsRegion       string
	RTOCInterval    time.Duration
	MyRoutingTables []string
	Logfile         string
}

var (
	config       natConfig
	pingschannel = make(chan bool)
)

func init() {
	// Parse config
	if _, err := toml.DecodeFile("natHealthConfig.conf", &config); err != nil {
		logging.Error.Println(err)
	}

	//Initalize logging
	logging.Log(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr, config.Logfile)

	// Run up Ping and HttpdHandler
	go httptools.HttpdHandler(config.HTTPPort)
	go hostping.Ping(config.OtherInstanceIP, pingschannel)

	//Process panic and error messages
	go func() {
		for err := range errhandling.ErrorChannel {
			logging.Info.Print(err)
		}
	}()
}

func main() {

	// Check that my routes belong to me.
	go func() {
		for {
			session := awsapitools.AwsSessIon(config.AwsRegion)
			RTsInIDs := awsapitools.DescribeRouteTableIDNatInstanceID(session, config.VpcID)
			for _, routeTable := range config.MyRoutingTables {
				if RTsInIDs[routeTable] != config.MyInstanceID {
					awsapitools.ReplaceRoute(session, routeTable, config.MyInstanceID)
				}
			}
			time.Sleep(config.RTOCInterval * time.Second)
			logging.Info.Print("Route Ownership check is sleeping 5 second\n")
		}
	}()

	for ping := range pingschannel {
		if !ping {
			logging.Error.Println("Nat instanceID:", config.OtherInstanceID, "instanceIP:", config.OtherInstanceIP, "is not pinging")
			respcode := httptools.RespCode("http://" + config.OtherInstanceIP + ":" + config.HTTPPort)
			if respcode != 200 {
				logging.Error.Println("Nat instanceID:", config.OtherInstanceID, "instanceIP:", config.OtherInstanceIP, "is returning http response code:", respcode)
				session := awsapitools.AwsSessIon(config.AwsRegion)
				instanceState := awsapitools.InstanceState(session, config.OtherInstanceID)

				if instanceState == "running" {
					RTsInIDs := awsapitools.DescribeRouteTableIDNatInstanceID(session, config.VpcID)
					for routeTableID, instanceID := range RTsInIDs {
						if instanceID != config.MyInstanceID {
							awsapitools.ReplaceRoute(session, routeTableID, config.MyInstanceID)
						}
					}
				}
			}
		}
	}
	// defer func() {
	// 	if e := recover(); e != nil {
	// 		errhandling.ErrorChannel <- errhandling.Error{fmt.Sprint(e)}
	// 	}
	// }()
}
