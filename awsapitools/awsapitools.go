package awsapitools

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

// AwsSessIon Returns AWS api session.
func AwsSessIon(region string) *ec2.EC2 {
	session := ec2.New(session.New(), &aws.Config{Region: aws.String(region)})
	return session
}

// DescribeRouteTableIDNatInstanceID Returns a map with RouteTableId InstanceId.
func DescribeRouteTableIDNatInstanceID(session *ec2.EC2, vpcid string) map[string]string {
	var rtIDInstID = make(map[string]string)

	params := &ec2.DescribeRouteTablesInput{
		DryRun: aws.Bool(false),
		Filters: []*ec2.Filter{
			{Name: aws.String("vpc-id"),
				Values: []*string{
					aws.String(vpcid),
				},
			},
		},
	}

	resp, err := session.DescribeRouteTables(params)

	if err != nil {
		panic(err)
	}
	for _, r := range resp.RouteTables {
		if r.Routes[1].InstanceId != nil {
			rtIDInstID[*r.Associations[0].RouteTableId] = *r.Routes[1].InstanceId
		}
	}
	return rtIDInstID
}

// ReplaceRoute the routing table route entry.
func ReplaceRoute(session *ec2.EC2, routeTableID, instanceID string) {
	params := &ec2.ReplaceRouteInput{
		DestinationCidrBlock: aws.String("0.0.0.0/0"),  // Required
		RouteTableId:         aws.String(routeTableID), // Required
		DryRun:               aws.Bool(false),
		InstanceId:           aws.String(instanceID),
	}

	resp, err := session.ReplaceRoute(params)

	if err != nil {
		panic(err)
	}
	fmt.Println(resp)
}

// InstanceState returns a sting with the instance state.
func InstanceState(session *ec2.EC2, instanceID string) string {
	params := &ec2.DescribeInstancesInput{
		InstanceIds: []*string{
			aws.String(instanceID),
		},
	}

	resp, err := session.DescribeInstances(params)

	if err != nil {
		panic(err)
	}

	instanceState := *resp.Reservations[0].Instances[0].State.Name
	return instanceState
}
