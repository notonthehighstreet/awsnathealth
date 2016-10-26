# About
This app is meant to monitor the aws nat box health.In the case of the nat box failure it will take over the routing table of the other nat box. The required number of nat boxes for the HA setup is 2.

# AWS auth/premission
The nat instance should have an ami role attached to it which allows read-only access for ec2 instance information retrieval and rw access for the vpc routing table management.

# Usage
### Cli help
```
Usage of ./github.com/notonthehighstreet/awsnathealth:

  -c, --config-file=/etc/awsnathealth.conf    Config file. Default is /etc/awsnathealth.conf.
  -v, --version                               awsnathealth Version.
```

# Config file example
```
# Nat Health Config
myInstancePubIP = "127.0.0.1"
otherInstancePubIP = "127.0.0.1"
httpport = "8001"
vpcID = "vpc-a6bc64d3"
awsRegion = "eu-west-1"
sessionCreateInterval = 3540
publicIPCheckInterval = 300
routeTableCheckInterval = 300
serviceCheckInterval = 300
pingTimeout = 30
myRoutingTables = [ 'rtb-7d5d5e19', 'rtb-1g73f07b' ]
otherInstanceRoutingTables = [ 'rtb-a31a98c7' ]
peerPubIPS = ['127.0.0.1', '127.0.0.1']
logfile = "/var/log/awsnathealth.log"
managedSecurityGroups = true
manageRacoonBgpd = true
standAlone = false
awsnathealthDisabled = false
debug = false
```

# Config parameters
**myInstancePubIP** - the nat instance elastic public ip.<br/>
**otherInstancePubIP** - the pair nat instance elastic public ip.<br/>
**httpport** - tcp port on which the http handeler is listening on. <br/>
**vpcID** - aws vpc id where the nat box is located.<br/>
**awsRegion** - aws region where the nat box is located.<br/>
**sessionCreateInterval** - the interval in seconds how often is the aws api session created.<br/>
**publicIPCheckInterval** - the interval in seconds how often is the elastic ip association checked.<br/>
**routeTableCheckInterval** - the interval in seconds how often is the routeTable association checked.<br/>
**serviceCheckInterval** - the interval in seconds how often is the service Racoon and Bgpd config is checked, it is used with **manageRacoonBgpd**.<br/>
**pingTimeout** - the number of failed pings before the route table is being taken over from the other instance.<br/>
**myRoutingTables** - the route tables associated with the nat instance.<br/>
**otherInstanceRoutingTables** - the route tables associated with the other nat instance.<br/>
**peerPubIPS** - the public addresses of the other vpn endpoints.<br/>
**logfile** - log file path.<br/>
**managedSecurityGroups** - if it set to true it'll create all the security rules to enable the IPSEC-GRE tunnel between the VPN peers and the natbox ping and http checks beetween the natbox pair.<br/>
This option requires the default security group id which is passed with the UserData with cloud formation when the stack is created.<br/>
**manageRacoonBgpd** - if it set to true it will managed the bgpd and racoon configs in case the natbox private ip changes.<br/>
**standAlone** - if it set to true the it wont check the other natbox health.<br/>
**awsnathealthDisabled** - if it set to true the awsnathealth app wont start up.<br/>
**debug** - if it set to true it enables verbose logging.


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

# Note

Please make sure that the user who runs the awsnathealt binary can create ICMP sockets on your linux distro.

You can check the following kernel parameter:

```
# cat /proc/sys/net/ipv4/ping_group_range
0	2147483647
```
You can set it in sysctl.conf permanatly so any user can create a ICMP socket with the below value.

```
net.ipv4.ping_group_range = 0 2147483647
```

Or you can set it to specific group id which the user who runs the awsnathealt needs to be member of. So if gid id 500 then the settings would be:

```
net.ipv4.ping_group_range = 500 500
```

#Create an RPM
Install fpm and rpm

```
gem install --no-ri --no-rdoc fpm
```
CD into the rpm folder and run

```
cd rpm/
fpm -s dir -t rpm -n "awsnathealth" -v 'app version'  --rpm-os linux --after-install scripts/after-install.sh --before-remove scripts/before-remove.sh --before-install scripts/before-install.sh etc usr
```

#Versioning

Versioning is baked into the build process.

If you pass in the **main.version** variable during the go install process the tool will return the set value if you run **awsnathealth -v**

```
cd awsnathealth/
GOOS=linux GOARCH=amd64 go build --ldflags "-X main.version='app version' -extldflags -static -s"
```

or

```
go install --ldflags "-X main.version='app version' -extldflags -static -s" -v github.com/notonthehighstreet/awsnathealth
```

Have fun...
