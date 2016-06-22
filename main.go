package main

import (
	"aws_nat/awsapitools"
	"aws_nat/errhandling"
	"aws_nat/hostping"
	"aws_nat/httptools"
	"aws_nat/logging"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/BurntSushi/toml"
	flag "github.com/docker/docker/pkg/mflag"
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
	config              natConfig
	pingschannel        = make(chan bool)
	version, configfile string
	ver                 bool
)

func init() {
	//Menu
	flag.StringVar(&configfile, []string{"c", "-config-file"}, "/etc/awsnathealth.conf", "Config file. Default is /etc/awsnathealth.conf.")
	flag.BoolVar(&ver, []string{"v", "-version"}, false, "awsnathealth Version.")
	flag.Parse()

	// Display app version
	if ver == true {
		fmt.Printf("Awsnathealth Version: %s\n", version)
		os.Exit(1)
	}

	//Check config file exist
	if _, err := os.Stat(configfile); err != nil {
		fmt.Printf("Config file: %s does not exist!\n", configfile)
		logging.Error.Printf("Config file: %s does not exist!\n", configfile)
		os.Exit(1)
	}

	//Parse config file.
	if _, err := toml.DecodeFile(configfile, &config); err != nil {
		logging.Error.Println(err)
	}

	//Initalize logging.
	logging.Log(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr, config.Logfile)

	//Run up Ping and HttpdHandler.
	go httptools.HttpdHandler(config.HTTPPort)
	go hostping.Ping(config.OtherInstancePubIP, pingschannel)

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
