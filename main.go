package main

import (
	"aws_nat/awsapitools"
	"aws_nat/hostping"
	"aws_nat/httpops"
	"fmt"

	"github.com/BurntSushi/toml"
)

func main() {

	type natConfig struct {
		MyInstanceID    string
		OtherInstanceID string
		OtherInstanceIP string
		HTTPPort        string
		VpcID           string
		AwsRegion       string
	}

	var (
		config       natConfig
		pingschannel = make(chan bool)
	)

	// Parse config
	if _, err := toml.DecodeFile("natHealthConfig.conf", &config); err != nil {
		panic(err)
	}
	// Run up Ping and HttpdHandler
	go httpops.HttpdHandler(config.HTTPPort)
	go hostping.Ping(config.OtherInstanceIP, pingschannel)

	for ping := range pingschannel {
		if ping {
			fmt.Print("True\n")
		}
		if !ping {
			fmt.Print("False\n")
			respcode := httpops.RespCode("http://" + config.OtherInstanceIP + ":" + config.HTTPPort)
			fmt.Print(respcode)
			if respcode != 200 {

				fmt.Print("RespCode Not 200\n")
				session := awsapitools.AwsSessIon(config.AwsRegion)
				instanceState := awsapitools.InstanceState(session, config.OtherInstanceID)

				if instanceState == "running" {
					rt := awsapitools.DescribeRouteTableIDNatInstanceID(session, config.VpcID)
					for routeTableID, instanceID := range rt {
						fmt.Println(routeTableID, instanceID)
						if instanceID != config.MyInstanceID {
							awsapitools.ReplaceRoute(session, routeTableID, config.MyInstanceID)
						}
					}
				}
			}
		}
	}
}

// func main() {
//
// 	type natConfig struct {
// 		MyInstanceID    string
// 		OtherInstanceID string
// 		OtherInstanceIP string
// 		HTTPPort        string
// 		VpcID           string
// 		AwsRegion       string
// 	}
//
// 	var config natConfig
// 	if _, err := toml.DecodeFile("natConfig.conf", &config); err != nil {
// 		panic(err)
// 	}
// 	fmt.Printf("MyInstanceID: %s\n", config.MyInstanceID)
// 	fmt.Printf("OtherInstanceID: %s\n", config.OtherInstanceID)
// 	fmt.Printf("OtherInstanceIP: %s\n", config.OtherInstanceIP)
// 	fmt.Printf("HTTPPort: %s\n", config.HTTPPort)
// 	fmt.Printf("VpcID: %s\n", config.VpcID)
// 	fmt.Printf("AwsRegion: %s\n", config.AwsRegion)
// }

// func main() {
//
// 	session := awsapitools.AwsSessIon("eu-west-1")
//
// 	state := awsapitools.InstanceState(session, "i-08ece580")
// 	fmt.Print(state)
//
// }

// session := awsapitools.AwsSessIon("eu-west-1")
// rt := awsapitools.DescribeRouteTableIDNatInstanceID(session, "vpc-b6dd64d3")
//
// for routeTableID, instanceID := range rt {
// 	fmt.Println(routeTableID, instanceID)
// 	if instanceID != "i-09755883" {
// 		awsapitools.ReplaceRoute(session, routeTableID, "i-09755883")
// 	}
// }

// session := awsapitools.AwsSessIon("eu-west-1")
// rt := awsapitools.DescribeRouteTableIDNatInstanceID(session, "vpc-b6dd64d3")
//
// for routeTableID, instanceID := range rt {
// 	fmt.Println(routeTableID, instanceID)
// 	if instanceID != "i-08ece580" {
// 		awsapitools.ReplaceRoute(session, routeTableID, "i-08ece580")
// 	}
// }

// go func() {
// 	httpops.HttpdHandler("8001")
// }()
// for {
// 	code := httpops.RespCode("http://localhost:8001")
// 	// code := httpops.RespCode("http://google.com")
// 	fmt.Print(code)
// 	time.Sleep(2 * time.Second)
// }
