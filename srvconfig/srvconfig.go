package srvconfig

import (
	"github.com/notonthehighstreet/awsnathealth/errhandling"
	"github.com/notonthehighstreet/awsnathealth/logging"
	"github.com/notonthehighstreet/awsnathealth/othertools"
)

// Config is somethings
// type Config struct{ template, configFile string }

// ManageServiceConfig manages
func ManageServiceConfig() {
	//Catch and log panic events
	var err error
	defer errhandling.CatchPanic(&err, "ManageServiceConfig")

	serviceConfigfileTemplatefile := map[string]*struct{ template, configFile string }{
		"racoon": {"/etc/racoon/ipsec-tools.sh.tmpl", "/etc/racoon/ipsec-tools.sh"},
		"bgpd":   {"/etc/quagga/bgpd.conf.tmpl", "/etc/quagga/bgpd.conf"},
	}

	config := map[string]string{
		"privateIP": othertools.GetLocalIP(),
	}
	for service, configTemplateFile := range serviceConfigfileTemplatefile {

		othertools.TempalteParse(configTemplateFile.template, configTemplateFile.configFile+".tmp", config)
		currentConfigMd5 := othertools.HashFileMd5(configTemplateFile.configFile)
		tmpConfigMd5 := othertools.HashFileMd5(configTemplateFile.configFile + ".tmp")

		if currentConfigMd5 != tmpConfigMd5 {
			othertools.CmdExec("mv", []string{"-f", configTemplateFile.configFile + ".tmp", configTemplateFile.configFile})
			othertools.CmdExec("service", []string{service, "restart"})
			if service == "racoon" {
				othertools.CmdExec("chmod", []string{"+x", configTemplateFile.configFile})
				othertools.CmdExec(configTemplateFile.configFile, nil)
			}
			logging.Info.Println("Service: ", service, " has restarted")
		}
	}
}
