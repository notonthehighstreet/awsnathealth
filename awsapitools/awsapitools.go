package awsapitools

import (
	"fmt"

	"github.com/notonthehighstreet/awsnathealth/errhandling"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

// AwsSessIon Returns AWS api session.
func AwsSessIon(region string) *ec2.EC2 {
	session := ec2.New(session.New(), &aws.Config{Region: aws.String(region)})
	return session
}

// DescribeRouteTableIDNatInstanceID Returns a map with RouteTableId  with the associated Nat InstanceId.
func DescribeRouteTableIDNatInstanceID(session *ec2.EC2, vpcid string) map[string]string {
	//Catch and log panic events
	var err error
	defer errhandling.CatchPanic(&err, "DescribeRouteTableIDNatInstanceID")

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
		for _, rt := range r.Routes {
			if rt.InstanceId != nil {
				rtIDInstID[*r.Associations[0].RouteTableId] = *rt.InstanceId
			} else {
				rtIDInstID[*r.Associations[0].RouteTableId] = "not_assigned"
			}
		}
	}
	return rtIDInstID
}

// ReplaceRoute replaces the routing table route instance entry.
func ReplaceRoute(session *ec2.EC2, routeTableID, instanceID string) {
	//Catch and log panic events
	var err error
	defer errhandling.CatchPanic(&err, "ReplaceRoute")

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
	if resp == nil {
		fmt.Println(resp)
	}
}

// InstanceStatebyInstanceID returns a sting with the instance state.
func InstanceStatebyInstanceID(session *ec2.EC2, instanceID string) string {
	//Catch and log panic events
	var err error
	defer errhandling.CatchPanic(&err, "InstanceStatebyInstanceID")

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

// InstanceStatebyInstancePubIP returns a sting with the instance state.
func InstanceStatebyInstancePubIP(session *ec2.EC2, instancePublicIP string) string {
	//Catch and log panic events
	var err error
	defer errhandling.CatchPanic(&err, "InstanceStatebyInstancePubIP")

	params := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			&ec2.Filter{
				Name: aws.String("ip-address"),
				Values: []*string{
					aws.String(instancePublicIP),
				},
			},
		},
	}

	resp, err := session.DescribeInstances(params)
	if err != nil {
		panic(err)
	}
	instanceState := *resp.Reservations[0].Instances[0].State.Name
	return instanceState
}

//InstanceIDbyPublicIP returns a sting with the instanceID.
func InstanceIDbyPublicIP(session *ec2.EC2, instancePublicIP string) string {
	//Catch and log panic events
	var err error
	defer errhandling.CatchPanic(&err, "InstanceIDbyPublicIP")

	params := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			&ec2.Filter{
				Name: aws.String("ip-address"),
				Values: []*string{
					aws.String(instancePublicIP),
				},
			},
		},
	}
	resp, err := session.DescribeInstances(params)
	if err != nil {
		panic(err)
	}
	instanceID := *resp.Reservations[0].Instances[0].InstanceId
	return instanceID
}

// MetadataInstanceID returns instanceID.
func MetadataInstanceID() string {
	//Catch and log panic events
	var err error
	defer errhandling.CatchPanic(&err, "MetadataInstanceID")

	session := ec2metadata.New(session.New(), &aws.Config{Endpoint: aws.String("http://169.254.169.254/latest")})
	resp, err := session.GetInstanceIdentityDocument()
	if err != nil {
		panic(err)
	}
	return resp.InstanceID
}

// AssociateElacticIP function associate Elatic IP to an instance.
func AssociateElacticIP(session *ec2.EC2, elaticIP, instanceID string) {
	//Catch and log panic events
	var err error
	defer errhandling.CatchPanic(&err, "AssociateElacticIP")

	params := &ec2.AssociateAddressInput{
		AllowReassociation: aws.Bool(true),
		DryRun:             aws.Bool(false),
		InstanceId:         aws.String(instanceID),
		PublicIp:           aws.String(elaticIP),
	}
	resp, err := session.AssociateAddress(params)
	if err != nil {
		panic(err)
	}
	if resp == nil {
		fmt.Println(resp)
	}
}

// InstancePublicIP returns a sting with the instance Elactic IP.
func InstancePublicIP(session *ec2.EC2, instanceID string) string {
	//Catch and log panic events
	var (
		err      error
		publicIP string
	)
	defer errhandling.CatchPanic(&err, "InstancePublicIP")

	params := &ec2.DescribeInstancesInput{
		InstanceIds: []*string{
			aws.String(instanceID),
		},
	}

	resp, err := session.DescribeInstances(params)
	if err != nil {
		panic(err)
	}

	if *resp.Reservations[0].Instances[0].PublicDnsName != "" {
		publicIP = *resp.Reservations[0].Instances[0].NetworkInterfaces[0].Association.PublicIp
	} else {
		publicIP = "has_no_PublicIP"
	}
	return publicIP

}
