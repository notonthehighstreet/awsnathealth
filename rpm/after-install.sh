echo "0 2147483647" > /proc/sys/net/ipv4/ping_group_range
grep -q net.ipv4.ping_group_range /etc/sysctl.conf || echo "net.ipv4.ping_group_range = 0 2147483647" >> /etc/sysctl.conf
chkconfig awsnh on --level 345
