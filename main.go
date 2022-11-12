package main

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"regexp"
	"strconv"
	"syscall"
	"time"

	lvn "github.com/Lavina-Tech-LLC/lavinagopackage/v2"
)

type (
	IpConf struct {
		Ip        string
		Interface string
	}
)

var (
	currentlyOn bool
)

func main() {
	go macBookPower()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	fmt.Println("Shuting down Server ...")
}

func macBookPower() {
	currentlyOn = false
	fmt.Println("Controlling Mason's macBook power")

	for {
		batt, err := getBatteryPercentage()
		if err != nil {
			lvn.Logger.Error("Cannot get battery info: " + err.Error())
			time.Sleep(time.Minute * 10)
		}

		if batt < 30 {
			http.Get("https://as-apia.coolkit.cc/v2/smartscene2/webhooks/execute?id=8fb139b880ca42f5bdf01976811024ec")
			if !currentlyOn {
				lvn.Logger.Infof("Turning on charger. battery is %v", batt)
			}
			currentlyOn = true
		}

		if batt > 90 {
			http.Get("https://as-apia.coolkit.cc/v2/smartscene2/webhooks/execute?id=a11b479d457544a79db7a07f4b275fd6")
			if currentlyOn {
				lvn.Logger.Infof("Turning off charger. battery is %v", batt)
			}
			currentlyOn = false
		}
		time.Sleep(time.Second * 10)
	}

}

func getBatteryPercentage() (int, error) {

	cmd := exec.Command(
		"pmset",
		"-g", "batt",
	)

	cmd.SysProcAttr = &syscall.SysProcAttr{
		Foreground: false,
	}
	out, err := cmd.Output()
	if err != nil {
		return 0, err
	}
	return parseCmd(string(out))
}

func parseCmd(s string) (int, error) {
	re := regexp.MustCompile(`[0-9]{1,3}[%]`)
	match := re.FindStringSubmatch(s)
	if len(match) < 1 {
		return 100, nil
	}
	str := (match[0])[:len(match[0])-1]

	return strconv.Atoi(str)
}

func sendData(powerNeeded int) {

}
