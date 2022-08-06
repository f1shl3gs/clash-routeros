package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"path"
	"syscall"
	"time"

	C "github.com/Dreamacro/clash/constant"
	"github.com/Dreamacro/clash/hub/executor"
	"github.com/Dreamacro/clash/hub/route"
	"github.com/Dreamacro/clash/log"

	"go.uber.org/automaxprocs/maxprocs"
)

func main() {
	subscription, interval := subscriptionFromEnv()

	maxprocs.Set(maxprocs.Logger(func(string, ...any) {}))
	home, err := os.Getwd()
	if err != nil {
		log.Fatalln("get current directory failed, %s", err)
	}

	// Make sure Country.mmdb and ui is under this directory
	C.SetHomeDir(home)
	route.SetUIPath(path.Join(home, "ui"))

	if err := loadConfigFromRemote(subscription, true); err != nil {
		log.Fatalln("first configuration init failed, %s", err)
	}

	// Start ExternalController
	go route.Start(":9090", "")

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
		case <-sigCh:
			return
		}

		if err := loadConfigFromRemote(subscription, false); err != nil {
			log.Warnln("update config failed, %s", err)
		}
	}
}

func subscriptionFromEnv() (string, time.Duration) {
	subscription := os.Getenv("SUBSCRIPTION")
	if subscription == "" {
		log.Fatalln("SUBSCRIPTION env must be provided")
	}

	value := os.Getenv("SUBSCRIPTION_UPDATE_INTERVAL")
	if value == "" {
		log.Infoln("SUBSCRIPTION_UPDATE_INTERVAL not provided use default 24h")
		return subscription, 24 * time.Hour
	}

	interval, err := time.ParseDuration(value)
	if err != nil {
		log.Fatalln("SUBSCRIPTION_UPDATE_INTERVAL must an valid duration string")
	}

	return subscription, interval
}

func loadConfigFromRemote(subscription string, force bool) error {
	data, err := fetch(subscription)
	if err != nil {
		return err
	}

	cfg, err := executor.ParseWithBytes(data)
	if err != nil {
		return err
	}

	executor.ApplyConfig(cfg, force)

	return nil
}

func fetch(endpoint string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code %d %s", resp.StatusCode, resp.Status)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response failed, err: %s", err)
	}

	return data, nil
}
