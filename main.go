package main

import (
	"aws_nat/awsroutingtable"
	"aws_nat/hostping"
	"aws_nat/httpops"
	"fmt"
)

func main() {
	var pingschannel = make(chan bool)
	// session := awsroutingtable.AwsSessIon("eu-west-1")
	// rt := awsroutingtable.DescribeRouteTableIDNatInstanceID(session, "vpc-b6dd64d3")
	//
	// for routeTableID, instanceID := range rt {
	// 	fmt.Println(routeTableID, instanceID)
	// 	if instanceID != "i-09755883" {
	// 		awsroutingtable.ReplaceRoute(session, routeTableID, "i-09755883")
	// 	}
	// }

	go httpops.HttpdHandler("8001")
	go hostping.Ping("8.8.8.9", pingschannel)

	for ping := range pingschannel {
		if ping {
			fmt.Print("True\n")
		}
		if !ping {
			fmt.Print("False\n")
			session := awsroutingtable.AwsSessIon("eu-west-1")
			rt := awsroutingtable.DescribeRouteTableIDNatInstanceID(session, "vpc-b6dd64d3")

			for routeTableID, instanceID := range rt {
				fmt.Println(routeTableID, instanceID)
				if instanceID != "i-08ece580" {
					awsroutingtable.ReplaceRoute(session, routeTableID, "i-08ece580")
				}
			}

		}
	}

	// go func() {
	// 	httpops.HttpdHandler("8001")
	// }()
	// for {
	// 	code := httpops.RespCode("http://localhost:8001")
	// 	// code := httpops.RespCode("http://google.com")
	// 	fmt.Print(code)
	// 	time.Sleep(2 * time.Second)
	// }
}
