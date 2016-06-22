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
	OtherInstancePubIP string        `toml:"otherInstancePubIP"`
	HTTPPort           string        `toml:"httpport"`
	VpcID              string        `toml:"vpcID"`
	AwsRegion          string        `toml:"awsRegion"`
	RTCInterval        time.Duration `toml:"RouteTableCheckInterval"`
	MyRoutingTables    []string      `toml:"myRoutingTables"`
	Logfile            string        `toml:"logfile"`
}

var (
	config       natConfig
	pingschannel = make(chan bool)
)

func init() {
	//Parse config file.
	if _, err := toml.DecodeFile("natHealthConfig.conf", &config); err != nil {
		logging.Error.Println(err)
	}

	//Initalize logging.
	logging.Log(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr, config.Logfile)

	//Run up Ping and HttpdHandler.
	go httptools.HttpdHandler(config.HTTPPort)
	go func() { hostping.Ping(config.OtherInstancePubIP, pingschannel) }()

	//Process panic and error messages.
	go func() {
		for err := range errhandling.ErrorChannel {
			logging.Info.Print(err)
		}
	}()
}

func main() {
	//Get myInstanceID
	myInstanceID := awsapitools.MetadataInstanceID()
	//Check that my routes belong to me.
	go func() {
		for {
			session := awsapitools.AwsSessIon(config.AwsRegion)
			RTsInIDs := awsapitools.DescribeRouteTableIDNatInstanceID(session, config.VpcID)
			for _, routeTable := range config.MyRoutingTables {
				if RTsInIDs[routeTable] != myInstanceID {
					logging.Info.Print("Takeing back my route table:", routeTable)
					awsapitools.ReplaceRoute(session, routeTable, myInstanceID)
				}
			}
			time.Sleep(config.RTCInterval * time.Second)
		}
	}()

	//Check the other nat insance
	for ping := range pingschannel {
		if !ping {
			//Create session to aws api.
			session := awsapitools.AwsSessIon(config.AwsRegion)
			otherInstanceID := awsapitools.InstanceIDbyPublicIP(session, config.OtherInstancePubIP)
			logging.Error.Println("Nat instanceID:", otherInstanceID, "instanceIP:", config.OtherInstancePubIP, "is not pinging")
			//Check is the other nat instances http handler returns 200.
			respcode := httptools.RespCode("http://" + config.OtherInstancePubIP + ":" + config.HTTPPort)
			if respcode != 200 {
				logging.Error.Println("Nat instanceID:", otherInstanceID, "instanceIP:", config.OtherInstancePubIP, "is returning http response code:", respcode)
				//Return the other nat instance state.
				instanceState := awsapitools.InstanceStatebyInstancePubIP(session, config.OtherInstancePubIP)
				//If the other instance state is not pending.
				if instanceState != "pending" {
					RTsInIDs := awsapitools.DescribeRouteTableIDNatInstanceID(session, config.VpcID)
					//Check who owns the routes if not me take them.
					for routeTableID, instanceID := range RTsInIDs {
						if instanceID != myInstanceID {
							awsapitools.ReplaceRoute(session, routeTableID, myInstanceID)
						}
					}
				}
			}
		}
	}
}
