# About
This app is meant to monitor the aws nat box health.

# AWS auth/premission
The nat instance should have an ami role attached to it which allows read-only access for ec2 instance information retrieval and rw access for the vpc routing table management.

# Usage
### Cli help
```
Usage of ./aws_nat:

  -c, --config-file=/etc/awsnathealth.conf    Config file. Default is /etc/awsnathealth.conf.
  -v, --version                               awsnathealth Version.
```

# Config file example
```
# Nat Health Config
otherInstancePubIP = "52.45.65.23"
httpport = "8001"
vpcID = "vpc-b6dd64d3"
awsRegion = "eu-west-1"
RouteTableCheckInterval = 10
myRoutingTables = [ "rtb-7d5dde19", "rtb-1f73f07b"]
logfile = "awsnathealth.log"
```

# Application work flow
![alt tag](workflow.png)


# Build the tool
You need to have go installed on the system where you wish to compile it.
For more information in regards go installation please check https://golang.org/doc/install

```
go get github.com/notonthehighstreet/awsnathealth
go install github.com/notonthehighstreet/awsnathealth

```

If the git repo requires ssh key auth you might want to set the global git config to over write the https protocol with the git one. If this is the case please run the following line.

```
git config --global url."git@github.com:".insteadOf "https://github.com/"
```

After running go install if your $GOPATH is set correctly you should find the binary in your $GOPATH/bin folder.


#Versioning

Versioning is baked into the build process.

If you pass in the **main.version** variable during the go install process the tool will return the set value if you run **awsnathealth -v**

```
go install --ldflags "-X main.version='app version' -extldflags -static -s" -v github.com/notonthehighstreet/awsnathealth
```

Have fun...