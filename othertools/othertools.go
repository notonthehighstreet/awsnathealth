package othertools

import (
	"crypto/md5"
	"encoding/hex"
	"io/ioutil"
	"net"
	"os"
	"os/exec"

	"github.com/notonthehighstreet/awsnathealth/errhandling"

	"text/template"
)

// StringInSlice function checks if the slice contains the given string, return bool.
func StringInSlice(str string, list []string) bool {
	for _, v := range list {
		if v == str {
			return true
		}
	}
	return false
}

// GetLocalIP returns the non loopback local IP of the host
func GetLocalIP() string {
	//Catch and log panic events
	var err error
	defer errhandling.CatchPanic(&err, "GetLocalIP")

	addrs, err := net.InterfaceAddrs()
	if err != nil {
		panic(err)
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}

// TempalteParse parses go templates and drops the files to the file system.
func TempalteParse(templateFilePath, filePath string, config map[string]string) {
	//Catch and log panic events
	var err error
	defer errhandling.CatchPanic(&err, "TempalteParse")

	t, err := template.ParseFiles(templateFilePath)
	if err != nil {
		panic(err)
	}

	f, err := os.Create(filePath)
	if err != nil {
		panic(err)
	}

	err = t.Execute(f, config)
	if err != nil {
		panic(err)
	}
	f.Close()
}

// CmdExec can exec cli commands
func CmdExec(cmd string, args []string) {
	//Catch and log panic events
	var err error
	defer errhandling.CatchPanic(&err, "CmdExec")

	if err := exec.Command(cmd, args...).Run(); err != nil {
		panic(err)
	}
}

// HashFileMd5 returns the file md5 has
func HashFileMd5(filePath string) string {
	//Catch and log panic events
	var err error
	defer errhandling.CatchPanic(&err, "HashFileMd5")

	hasher := md5.New()
	s, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic(err)
	}
	hasher.Write(s)
	if err != nil {
		panic(err)
	}
	return hex.EncodeToString(hasher.Sum(nil))
}

// // ManageServiceConfig manages
// func ManageServiceConfig(serviceConfigfileTemplatefile map[string]*struct{ template, configFile string }) {
// 	//Catch and log panic events
// 	var err error
// 	defer errhandling.CatchPanic(&err, "ManageServiceConfig")
//
// 	config := map[string]string{
// 		"privateIP": GetLocalIP(),
// 	}
// 	for service, configTemplateFile := range serviceConfigfileTemplatefile {
//
// 		TempalteParse(configTemplateFile.template, "/tmp/"+configTemplateFile.configFile+".tmp", config)
// 		currentConfigMd5 := HashFileMd5(configTemplateFile.configFile)
// 		tmpConfigMd5 := HashFileMd5("/tmp/" + configTemplateFile.configFile + ".tmp")
//
// 		if currentConfigMd5 != tmpConfigMd5 {
// 			CmdExec("cp", []string{"/tmp/" + configTemplateFile.configFile + ".tmp", configTemplateFile.configFile})
// 			CmdExec("rm -rf", []string{"/tmp/" + configTemplateFile.configFile + ".tmp"})
// 			CmdExec("service restart", []string{service})
// 		}
// 	}
// }
