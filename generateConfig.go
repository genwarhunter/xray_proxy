package main

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"sync"
)

func GenerateConfig(link string, wg *sync.WaitGroup) {
	defer wg.Done()
	parseUrl, err := url.Parse(link)
	if err != nil {
		return
	}
	var rc rawConfig

	switch parseUrl.Scheme {
	case "vless":
		parseVless(parseUrl, &rc)
	case "vmess":
		parseVmess(parseUrl, &rc)
	case "trojan":
		parseTrojan(parseUrl, &rc)
	case "ss":
		parseSS(parseUrl, &rc)
	case "ssr":
		parseSSR(parseUrl, &rc)
	default:
		return
	}

	jsonConf := buildConfig(&rc)
	hash := md5.Sum([]byte(jsonConf))
	hashString := hex.EncodeToString(hash[:])
	HashProtocolMap.Store(hashString, rc.protocol)
	createConfigFile(hashString, jsonConf)
	return
}

func parseVless(link *url.URL, rc *rawConfig) {
	rc.id = link.User.Username()
	var query = link.Query()
	if path, ok := query["path"]; ok {
		rc.path = path[0]
	}
	if enc, ok := query["encryption"]; ok {
		rc.enc = enc[0]
	} else {
		rc.enc = "none"
	}
	if security, ok := query["security"]; ok {
		rc.security = security[0]
	}
	if tls, ok := query["security"]; ok {
		rc.tls = tls[0]
	}
	rc.ip = strings.Split(link.Host, ":")[0]
	if host, ok := query["host"]; ok {
		rc.host = host[0]
	}
	if Type, ok := query["type"]; ok {
		rc.Type = Type[0]
		rc.net = Type[0]
	}
	port, _ := strconv.Atoi(strings.Split(link.Host, ":")[1])
	rc.port = uint16(port)
	rc.protocol = link.Scheme
	if serviceName, ok := query["serviceName"]; ok {
		rc.serviceName = serviceName[0]
	}
	if spiderX, ok := query["spx"]; ok {
		rc.spiderX = spiderX[0]
	}
	if shortId, ok := query["sid"]; ok {
		rc.shortId = shortId[0]
	}
	if publicKey, ok := query["pbk"]; ok {
		rc.publicKey = publicKey[0]
	}
	if serverName, ok := query["sni"]; ok {
		rc.serverName = serverName[0]
	}
	if fingerprint, ok := query["fp"]; ok {
		rc.fingerprint = fingerprint[0]
	}
	if flow, ok := query["flow"]; ok {
		rc.flow = flow[0]
	}

}

func parseVmess(link *url.URL, rc *rawConfig) {
	decoded, err := base64.StdEncoding.DecodeString(link.Host)
	if err != nil {
		return
	}
	var c VmessConfig
	err = json.Unmarshal(decoded, &c)
	rc.protocol = link.Scheme
	rc.ip = c.Add
	rc.id = c.ID
	rc.host = c.Host
	rc.port = c.Port
	rc.path = c.Path
	rc.net = c.Net
	rc.Type = c.Type
	rc.fp = c.Fp
	rc.security = c.Scy
	rc.serverName = c.Sni
	rc.tls = c.Tls
}

func parseTrojan(link *url.URL, rc *rawConfig) {
	rc.id = link.User.Username()
	var query = link.Query()
	if path, ok := query["path"]; ok {
		rc.path = path[0]
	}
	if enc, ok := query["encryption"]; ok {
		rc.enc = enc[0]
	} else {
		rc.enc = "none"
	}
	if security, ok := query["security"]; ok {
		rc.security = security[0]
	}
	if tls, ok := query["security"]; ok {
		rc.tls = tls[0]
	}
	rc.ip = strings.Split(link.Host, ":")[0]
	if host, ok := query["host"]; ok {
		rc.host = host[0]
	}
	if Type, ok := query["type"]; ok {
		rc.Type = Type[0]
		rc.net = Type[0]
	}
	port, _ := strconv.Atoi(strings.Split(link.Host, ":")[1])
	rc.port = uint16(port)
	rc.protocol = link.Scheme
	if serviceName, ok := query["serviceName"]; ok {
		rc.serviceName = serviceName[0]
	}
	if spiderX, ok := query["spx"]; ok {
		rc.spiderX = spiderX[0]
	}
	if shortId, ok := query["sid"]; ok {
		rc.shortId = shortId[0]
	}
	if publicKey, ok := query["pbk"]; ok {
		rc.publicKey = publicKey[0]
	}
	if serverName, ok := query["sni"]; ok {
		rc.serverName = serverName[0]
	}
	if fingerprint, ok := query["fp"]; ok {
		rc.fingerprint = fingerprint[0]
	}
	if flow, ok := query["flow"]; ok {
		rc.flow = flow[0]
	}
	if alpn, ok := query["alpn"]; ok {
		rc.alpn = alpn[0]
	}
}

func parseSS(link *url.URL, rc *rawConfig) {
	// Implementation here
}

func parseSSR(link *url.URL, rc *rawConfig) {
	// Implementation here
}

func buildConfig(rc *rawConfig) string {
	conf := Config{
		Log: &LOG{Loglevel: "warning"},
		Dns: &DNS{
			DisableFallback: true,
			Servers: []Server{
				{Address: "https://8.8.8.8/dns-query"},
				{Address: "localhost", Domains: []string{"full:" + rc.host}},
			},
			Tag: "dns",
		},
		Stats: &STATS{},
		Routing: &Routing{
			DomainStrategy: "AsIs",
			Rules: []Rule{
				{InboundTag: []string{"socks-in"}, OutboundTag: "dns-out", Port: "53", Type: "field"},
				{OutboundTag: "proxy", Port: "0-65535", Type: "field"},
			},
		},
		Policy: &Policy{
			Levels: map[int]level{1: {30}},
			System: SYSTEM{
				StatsOutboundDownlink: true,
				StatsOutboundUplink:   true,
			},
		},
		Inbounds: &[]Inbound{
			{
				Listen:   AppConfig.Ip,
				Port:     0,
				Protocol: "socks",
				Settings: map[string]bool{"udp": true},
				Tag:      "socks-in",
				Sniffing: sniffing{
					Enabled:      true,
					MetadataOnly: false,
					RouteOnly:    true,
				},
			},
		},
		Outbounds: buildOutbounds(rc),
	}

	confJson, _ := json.Marshal(conf)
	//log.Println(string(confJson))
	return string(confJson)
}

func buildOutbounds(rc *rawConfig) *[]Outbound {
	var VNex vnext
	var Serv Server
	if rc.protocol != "trojan" {
		VNex = vnext{
			Address: rc.ip,
			Port:    getPort(rc.port),
			Users: []USER{
				{Id: rc.id, AlterId: 0, Security: "auto", Encryption: rc.enc, Flow: rc.flow},
			},
		}
	} else {
		Serv.Address = rc.ip
		Serv.Port = rc.port
		Serv.Password = rc.id
	}
	outbounds := []Outbound{
		{
			DomainStrategy: "AsIs",
			Protocol:       rc.protocol,
			Settings: &OUTBOUNDSETTING{
				Vnext: []vnext{
					VNex,
				},
				Servers: []Server{
					Serv,
				},
			},
			StreamSettings: buildStreamSettings(rc),
			Tag:            "proxy",
		},
		{Protocol: "freedom", Tag: "direct"},
		{Protocol: "freedom", Tag: "bypass"},
		{Protocol: "blackhole", Tag: "block"},
		{
			Protocol:      "dns",
			ProxySettings: &PROXYSETTINGS{Tag: "proxy", TransportLayer: true},
			Settings:      &OUTBOUNDSETTING{Address: "8.8.8.8", Network: "tcp", Port: 53, UserLevel: 1},
			Tag:           "dns-out",
		},
	}
	return &outbounds
}

func buildStreamSettings(rc *rawConfig) *StreamSettingsObject {
	ss := &StreamSettingsObject{Network: rc.net, Security: rc.security}
	switch rc.net {
	case "ws":
		ss.WsSettings = &wsSettings{
			Path: rc.path,
			Host: rc.serverName,
			Headers: map[string]string{
				"path": rc.path,
				"host": rc.host,
			},
		}
	case "grpc":
		ss.GrpcSettings = &grpcSettings{}
		ss.RealitySettings = &REALITYSETTINGS{
			Fingerprint: rc.fingerprint,
			PublicKey:   rc.publicKey,
			ServerName:  rc.serverName,
			ShortId:     rc.shortId,
			SpiderX:     rc.spiderX,
		}
	case "tcp":
		ss.RealitySettings = &REALITYSETTINGS{
			Fingerprint: rc.fingerprint,
			PublicKey:   rc.publicKey,
			ServerName:  rc.serverName,
			ShortId:     rc.shortId,
			SpiderX:     rc.spiderX,
		}
	case "http":
		ss.RealitySettings = &REALITYSETTINGS{
			Fingerprint: rc.fingerprint,
			PublicKey:   rc.publicKey,
			ServerName:  rc.serverName,
			ShortId:     rc.shortId,
			SpiderX:     rc.spiderX,
		}
	case "httpupgrade":
		return ss
	default:
		ss.RealitySettings = &REALITYSETTINGS{
			Fingerprint: rc.fingerprint,
			PublicKey:   rc.publicKey,
			ServerName:  rc.serverName,
			ShortId:     rc.shortId,
			SpiderX:     rc.spiderX,
		}
	}
	getSecurity(rc, ss)
	return ss
}

func getPort(port uint16) uint16 {
	if port == 0 {
		return 443
	}
	return port
}

func isZeroValue(x interface{}) bool {
	return reflect.DeepEqual(x, reflect.Zero(reflect.TypeOf(x)).Interface())
}

func getSecurity(rc *rawConfig, s *StreamSettingsObject) {
	if rc.tls == "none" || rc.tls == "" {
		return
	}
	s.Security = rc.tls
	if !isZeroValue(s) {
		s.TlsSettings = &tlsSettings{
			ServerName: rc.serverName,
			Alpn:       []string{rc.alpn},
		}
	}
}

func createConfigFile(filename string, content string) {
	dir := "configs"

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.Mkdir(dir, 0755)
		if err != nil {
			return
		}
	}

	filePath := filepath.Join(dir, filepath.Base(filename))

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		err = os.WriteFile(filePath, []byte(content), 0644)
		if err != nil {
			return
		}
	} else {
		return
	}
}
