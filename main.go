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
	fmt.Println("Controlling Mason's macBook power")

	for {
		batt, chargerOn, err := batteryStatus()

		if err != nil {
			lvn.Logger.Error("Cannot get battery info: " + err.Error())
			time.Sleep(time.Minute * 5)
			continue
		}

		if batt < 30 && !chargerOn {
			lvn.Logger.Infof("Turning on charger. battery is %v", batt)
			http.Get("https://as-apia.coolkit.cc/v2/smartscene2/webhooks/execute?id=0e88474b75f04f17874ec186bc1b1a33")
		}

		if batt > 90 && chargerOn {
			lvn.Logger.Infof("Turning off charger. battery is %v", batt)
			http.Get("https://as-apia.coolkit.cc/v2/smartscene2/webhooks/execute?id=cfa90d7765054895ae2cc9aa126dc4c7")
		}
		time.Sleep(time.Second * 10)
	}

}

func batteryStatus() (int, bool, error) {
	batt, err := getBatteryPercentage()
	if err != nil {
		return 0, false, err
	}
	time.Sleep(1 * time.Minute)
	batt2, err := getBatteryPercentage()
	if err != nil {
		return 0, false, err
	}

	return batt2, batt2 > batt || batt2 == 100, nil
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
