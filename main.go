package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/notonthehighstreet/awsnathealth/awsapitools"
	"github.com/notonthehighstreet/awsnathealth/errhandling"
	"github.com/notonthehighstreet/awsnathealth/hostping"
	"github.com/notonthehighstreet/awsnathealth/httptools"
	"github.com/notonthehighstreet/awsnathealth/logging"
	"github.com/notonthehighstreet/awsnathealth/othertools"

	"github.com/BurntSushi/toml"
	flag "github.com/docker/docker/pkg/mflag"
)

type natConfig struct {
	MyInstancePubIP            string        `toml:"myInstancePubIP"`
	OtherInstancePubIP         string        `toml:"otherInstancePubIP"`
	HTTPPort                   string        `toml:"httpport"`
	VpcID                      string        `toml:"vpcID"`
	AwsRegion                  string        `toml:"awsRegion"`
	SCInterval                 time.Duration `toml:"sessionCreateInterval"`
	PICInterval                time.Duration `toml:"publicIPCheckInterval"`
	RTCInterval                time.Duration `toml:"routeTableCheckInterval"`
	MyRoutingTables            []string      `toml:"myRoutingTables"`
	OtherInstanceRoutingTables []string      `toml:"otherInstanceRoutingTables"`
	Logfile                    string        `toml:"logfile"`
	Debug                      bool          `toml:"debug"`
}

var (
	config              natConfig
	pingschannel        = make(chan bool)
	version, configfile string
	ver                 bool
	session             *ec2.EC2
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

	// Get aws session.
	go func() {
		for {
			session = awsapitools.AwsSessIon(config.AwsRegion)
			time.Sleep(config.SCInterval * time.Second)
		}
	}()

	//Process panic and error messages.
	if config.Debug {
		go func() {
			for err := range errhandling.ErrorChannel {
				logging.Error.Print(err)
			}
		}()
	} else {
		go func() {
			for err := range errhandling.ErrorChannel {
				_ = err
			}
		}()
	}
}

func main() {
	//Get myInstanceID
	myInstanceID := awsapitools.MetadataInstanceID()
	//Check that my routes belongs to me.
	go func() {
		for {
			RTsInIDs := awsapitools.DescribeRouteTableIDNatInstanceID(session, config.VpcID)
			for _, routeTable := range config.MyRoutingTables {
				if RTsInIDs[routeTable] != myInstanceID {
					logging.Info.Print("Taking back my route table:", routeTable)
					awsapitools.ReplaceRoute(session, routeTable, myInstanceID)
				}
			}
			time.Sleep(config.RTCInterval * time.Second)
		}
	}()

	//Check that my ElasticIP belongs to me.
	go func() {
		for {
			if awsapitools.InstancePublicIP(session, myInstanceID) != config.MyInstancePubIP {
				awsapitools.AssociateElacticIP(session, config.MyInstancePubIP, myInstanceID)
				logging.Info.Print("Taking back my Elatic IP:", config.MyInstancePubIP)
			}
			time.Sleep(config.PICInterval * time.Second)
		}
	}()

	//Check the other nat insance
	for ping := range pingschannel {
		if !ping {
			otherInstanceID := awsapitools.InstanceIDbyPublicIP(session, config.OtherInstancePubIP)
			logging.Warning.Println("Nat instanceID:", otherInstanceID, "instanceIP:", config.OtherInstancePubIP, "is not pinging")
			//Check is the other nat instances http handler returns 200.
			respcode := httptools.RespCode("http://" + config.OtherInstancePubIP + ":" + config.HTTPPort)
			if respcode != 200 {
				logging.Warning.Println("Nat instanceID:", otherInstanceID, "instanceIP:", config.OtherInstancePubIP, "is returning http response code:", respcode)
				//Return the other nat instance state.
				instanceState := awsapitools.InstanceStatebyInstancePubIP(session, config.OtherInstancePubIP)
				//If the other instance state is not pending.
				if instanceState != "pending" {
					RTsInIDs := awsapitools.DescribeRouteTableIDNatInstanceID(session, config.VpcID)
					bothrtable := append(config.MyRoutingTables, config.OtherInstanceRoutingTables...)
					//Check who owns the routes if not me take them.
					for routeTableID, instanceID := range RTsInIDs {
						if othertools.StringInSlice(routeTableID, bothrtable) && instanceID != myInstanceID {
							logging.Info.Println("I've taken over Nat instanceID:", otherInstanceID, "instanceIP:", config.OtherInstancePubIP, "Route table:", routeTableID)
							awsapitools.ReplaceRoute(session, routeTableID, myInstanceID)
						} else {
							logging.Error.Println("Route table:", routeTableID, "does not belong to nat instance:", otherInstanceID)
						}
					}
				}
			}
		}
	}
}
