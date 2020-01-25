// Copyright 2018 fatedier, fatedier@gmail.com
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package frpc

import (
	"context"
	"math/rand"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/charlesbao/frpc/client"
	"github.com/charlesbao/frpc/models/config"
	"github.com/charlesbao/frpc/utils/log"
	"github.com/fatedier/golib/crypto"
)

var (
	cfgFile     string
	showVersion bool

	serverAddr      string
	user            string
	protocol        string
	token           string
	logLevel        string
	logFile         string
	logMaxDays      int
	disableLogColor bool

	proxyName         string
	localIp           string
	localPort         int
	remotePort        int
	useEncryption     bool
	useCompression    bool
	customDomains     string
	subDomain         string
	httpUser          string
	httpPwd           string
	locations         string
	hostHeaderRewrite string
	role              string
	sk                string
	serverName        string
	bindAddr          string
	bindPort          int

	kcpDoneCh chan struct{}
)

func init() {
	crypto.DefaultSalt = "frp"
	rand.Seed(time.Now().UnixNano())
	kcpDoneCh = make(chan struct{})
}

func handleSignal(svr *client.Service) {
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
	svr.Close()
	time.Sleep(250 * time.Millisecond)
	close(kcpDoneCh)
}

func parseClientCommonCfgFromIni(content string) (config.ClientCommonConf, error) {
	cfg, err := config.UnmarshalClientConfFromIni(content)
	if err != nil {
		return config.ClientCommonConf{}, err
	}
	return cfg, err
}

func Run(cfgFilePath string) (err error) {
	var content string
	content, err = config.GetRenderedConfFromFile(cfgFilePath)
	if err != nil {
		return
	}

	cfg, err := parseClientCommonCfgFromIni(content)
	if err != nil {
		return
	}

	pxyCfgs, visitorCfgs, err := config.LoadAllConfFromIni(cfg.User, content, cfg.Start)
	if err != nil {
		return err
	}

	err = startService(cfg, pxyCfgs, visitorCfgs, cfgFilePath)
	return
}

func startService(cfg config.ClientCommonConf, pxyCfgs map[string]config.ProxyConf, visitorCfgs map[string]config.VisitorConf, cfgFile string) (err error) {
	log.InitLog(cfg.LogWay, cfg.LogFile, cfg.LogLevel,
		cfg.LogMaxDays, cfg.DisableLogColor)

	if cfg.DnsServer != "" {
		s := cfg.DnsServer
		if !strings.Contains(s, ":") {
			s += ":53"
		}
		// Change default dns server for frpc
		net.DefaultResolver = &net.Resolver{
			PreferGo: true,
			Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
				return net.Dial("udp", s)
			},
		}
	}
	svr, errRet := client.NewService(cfg, pxyCfgs, visitorCfgs, cfgFile)
	if errRet != nil {
		err = errRet
		return
	}

	// Capture the exit signal if we use kcp.
	if cfg.Protocol == "kcp" {
		go handleSignal(svr)
	}

	err = svr.Run()
	if cfg.Protocol == "kcp" {
		<-kcpDoneCh
	}
	return
}
