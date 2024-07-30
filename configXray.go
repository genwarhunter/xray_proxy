package main

type Config struct {
	Log       *LOG        `json:"log,omitempty"` //+
	Api       *API        `json:"api,omitempty"`
	Dns       *DNS        `json:"dns,omitempty"`     //+
	Routing   *Routing    `json:"routing,omitempty"` //+
	Policy    *Policy     `json:"policy,omitempty"`  //+
	Inbounds  *[]Inbound  `json:"inbounds,omitempty"`
	Outbounds *[]Outbound `json:"outbounds,omitempty"`
	Stats     *STATS      `json:"stats,omitempty"`   //+
	Reverse   *REVERSE    `json:"reverse,omitempty"` //+
	Fakedns   *FAKEDNS    `json:"fakedns,omitempty"` //+
	Metrics   *METRICS    `json:"metrics,omitempty"` //+
}

type API struct {
}

type LOG struct {
	Access   string `json:"access,omitempty"`
	Error    string `json:"error,omitempty"`
	Loglevel string `json:"loglevel,omitempty"`
	DnsLog   bool   `json:"dnsLog,omitempty"`
}

type STATS struct {
}
type REVERSE struct {
}
type FAKEDNS struct {
}
type METRICS struct {
}

type DNS struct {
	DisableFallback bool     `json:"disableFallback"`
	Servers         []Server `json:"servers"`
	Tag             string   `json:"tag"`
}

type Server struct {
	Address       string   `json:"address,omitempty"`
	Method        string   `json:"method,omitempty"`
	Domains       []string `json:"domains,omitempty"`
	Password      string   `json:"password,omitempty"`
	Port          uint16   `json:"port,omitempty"`
	Uot           bool     `json:"uot,omitempty"`
	QueryStrategy string   `json:"queryStrategy,omitempty"`
}

type Inbound struct {
	Listen   string          `json:"listen"`
	Port     uint16          `json:"port"`
	Protocol string          `json:"protocol"`
	Settings map[string]bool `json:"settings"`
	Tag      string          `json:"tag"`
	Sniffing sniffing        `json:"sniffing"`
}

type sniffing struct {
	Enabled      bool     `json:"enabled"`
	MetadataOnly bool     `json:"metadataOnly"`
	RouteOnly    bool     `json:"routeOnly"`
	DestOverride []string `json:"destOverride"`
}

type Routing struct {
	DomainStrategy string `json:"domainStrategy"`
	Rules          []Rule `json:"rules"`
}

type Rule struct {
	InboundTag  []string `json:"inboundTag,omitempty"`
	OutboundTag string   `json:"outboundTag"`
	Port        string   `json:"port"`
	Type        string   `json:"type"`
}

type Policy struct {
	Levels map[int]level `json:"levels"`
	System SYSTEM        `json:"system"`
}

type SYSTEM struct {
	StatsOutboundDownlink bool `json:"statsOutboundDownlink"`
	StatsOutboundUplink   bool `json:"statsOutboundUplink"`
}

type level struct {
	ConnIdle int `json:"connIdle"`
}

type Outbound struct {
	DomainStrategy string                `json:"domainStrategy"`
	Protocol       string                `json:"protocol,omitempty"`
	Settings       *OUTBOUNDSETTING      `json:"settings,omitempty"`
	StreamSettings *StreamSettingsObject `json:"streamSettings,omitempty"`
	Tag            string                `json:"tag,omitempty"`
	ProxySettings  *PROXYSETTINGS        `json:"proxySettings,omitempty"`
}

type PROXYSETTINGS struct {
	Tag            string `json:"tag"`
	TransportLayer bool   `json:"transportLayer"`
}

type OUTBOUNDSETTING struct {
	Vnext     []vnext  `json:"vnext,omitempty"`
	Address   string   `json:"address,omitempty"`
	Network   string   `json:"network,omitempty"`
	Port      uint16   `json:"port,omitempty"`
	UserLevel int      `json:"userLevel,omitempty"`
	Servers   []Server `json:"servers,omitempty"`
}

type vnext struct {
	Address string `json:"address"`
	Port    uint16 `json:"port"`
	Users   []USER `json:"users"`
}

type USER struct {
	AlterId    int    `json:"alterId"`
	Id         string `json:"id"`
	Security   string `json:"security"`
	Encryption string `json:"encryption"`
	Flow       string `json:"flow,omitempty"`
}

type StreamSettingsObject struct {
	Network             string               `json:"network,omitempty"`
	Security            string               `json:"security,omitempty"`
	TlsSettings         *tlsSettings         `json:"tlssettings,omitempty"`
	TcpSettings         *tcpSettings         `json:"tcpSettings,omitempty"`
	KcpSettings         *kcpSettings         `json:"kcpSettings,omitempty"`
	WsSettings          *wsSettings          `json:"wsSettings,omitempty"`
	HttpSettings        *httpSettings        `json:"httpSettings,omitempty"`
	QuicSettings        *quicSettings        `json:"quicSettings,omitempty"`
	DsSettings          *dsSettings          `json:"dsSettings,omitempty"`
	GrpcSettings        *grpcSettings        `json:"grpcSettings,omitempty"`
	HttpupgradeSettings *httpupgradeSettings `json:"httpupgradeSettings,omitempty"`
	SplithttpSettings   *splithttpSettings   `json:"splithttpSettings,omitempty"`
	Sockopt             *SOCKOPT             `json:"sockopt,omitempty"`
	RealitySettings     *REALITYSETTINGS     `json:"realitySettings,omitempty"`
}

type REALITYSETTINGS struct {
	Fingerprint string `json:"fingerprint,omitempty"`
	PublicKey   string `json:"publicKey,omitempty"`
	ServerName  string `json:"serverName,omitempty"`
	ShortId     string `json:"shortId,omitempty"`
	SpiderX     string `json:"spiderX,omitempty"`
}

type SOCKOPT struct {
	Mark                 int    `json:"mark,omitempty"`
	TcpMaxSeg            int    `json:"tcpMaxSeg,omitempty"`
	TcpFastOpen          bool   `json:"tcpFastOpen,omitempty"`
	Tproxy               string `json:"tproxy,omitempty"`
	DomainStrategy       string `json:"domainStrategy,omitempty"`
	DialerProxy          string `json:"dialerProxy,omitempty"`
	AcceptProxyProtocol  bool   `json:"acceptProxyProtocol,omitempty"`
	TcpKeepAliveInterval int    `json:"tcpKeepAliveInterval,omitempty"`
	TcpKeepAliveIdle     int    `json:"tcpKeepAliveIdle,omitempty"`
	TcpUserTimeout       int    `json:"tcpUserTimeout,omitempty"`
	TcpCongestion        string `json:"tcpCongestion,omitempty"`
	Interface            string `json:"interface,omitempty"`
	V6Only               bool   `json:"v6only,omitempty"`
	TcpWindowClamp       int    `json:"tcpWindowClamp,omitempty"`
	TcpMptcp             bool   `json:"tcpMptcp,omitempty"`
	TcpNoDelay           bool   `json:"tcpNoDelay,omitempty"`
}

type tlsSettings struct {
	ServerName                       string        `json:"serverName,omitempty"`
	RejectUnknownSni                 bool          `json:"rejectUnknownSni,omitempty"`
	AllowInsecure                    bool          `json:"allowInsecure,omitempty"`
	Alpn                             []string      `json:"alpn,omitempty"`
	MinVersion                       string        `json:"minVersion,omitempty"`
	MaxVersion                       string        `json:"maxVersion,omitempty"`
	CipherSuites                     string        `json:"cipherSuites,omitempty"`
	Certificates                     []interface{} `json:"certificates,omitempty"`
	DisableSystemRoot                bool          `json:"disableSystemRoot,omitempty"`
	EnableSessionResumption          bool          `json:"enableSessionResumption,omitempty"`
	Fingerprint                      string        `json:"fingerprint,omitempty"`
	PinnedPeerCertificateChainSha256 []string      `json:"pinnedPeerCertificateChainSha256,omitempty"`
	MasterKeyLog                     string        `json:"masterKeyLog,omitempty"`
}

type tcpSettings struct {
	AcceptProxyProtocol bool              `json:"acceptProxyProtocol"`
	Header              map[string]string `json:"header"`
}

type kcpSettings struct {
	Mtu              int  `json:"mtu"`
	Tti              int  `json:"tti"`
	UplinkCapacity   int  `json:"uplinkCapacity"`
	DownlinkCapacity int  `json:"downlinkCapacity"`
	Congestion       bool `json:"congestion"`
	ReadBufferSize   int  `json:"readBufferSize"`
	WriteBufferSize  int  `json:"writeBufferSize"`
	Header           struct {
		Type string `json:"type"`
	} `json:"header"`
	Seed string `json:"seed"`
}

type wsSettings struct {
	AcceptProxyProtocol bool              `json:"acceptProxyProtocol"`
	Path                string            `json:"path"`
	Host                string            `json:"host"`
	Headers             map[string]string `json:"headers"`
}

type httpSettings struct {
	Host               []string `json:"host"`
	Path               string   `json:"path"`
	ReadIdleTimeout    int      `json:"read_idle_timeout"`
	HealthCheckTimeout int      `json:"health_check_timeout"`
	Method             string   `json:"method"`
	Headers            struct {
		Header []string `json:"Header"`
	} `json:"headers"`
}

type quicSettings struct {
	Security string `json:"security"`
	Key      string `json:"key"`
	Header   struct {
		Type string `json:"type"`
	} `json:"header"`
}

type dsSettings struct {
	Path     string `json:"path"`
	Abstract bool   `json:"abstract"`
	Padding  bool   `json:"padding"`
}

type grpcSettings struct {
	ServiceName         string `json:"serviceName,omitempty"`
	MultiMode           bool   `json:"multiMode,omitempty"`
	IdleTimeout         int    `json:"idle_timeout,omitempty"`
	HealthCheckTimeout  int    `json:"health_check_timeout,omitempty"`
	PermitWithoutStream bool   `json:"permit_without_stream,omitempty"`
	InitialWindowsSize  int    `json:"initial_windows_size,omitempty"`
}

type httpupgradeSettings struct {
	AcceptProxyProtocol bool   `json:"acceptProxyProtocol"`
	Path                string `json:"path"`
	Host                string `json:"host"`
	Headers             struct {
		Key string `json:"key"`
	} `json:"headers"`
}

type splithttpSettings struct {
	Path    string `json:"path"`
	Host    string `json:"host"`
	Headers struct {
		Key string `json:"key"`
	} `json:"headers"`
}

type VmessConfig struct {
	Add  string `json:"add"`
	Aid  string `json:"aid"`
	Alpn string `json:"alpn"`
	Fp   string `json:"fp"`
	Host string `json:"host"`
	ID   string `json:"id"`
	Net  string `json:"net"`
	Path string `json:"path"`
	Port uint16 `json:"port"`
	Scy  string `json:"scy"`
	Sni  string `json:"sni"`
	Tls  string `json:"tls"`
	Type string `json:"type"`
	V    string `json:"v"`
	Ps   string `json:"ps"`
}

type rawConfig struct {
	protocol    string
	ip          string
	host        string
	port        uint16
	path        string
	net         string
	id          string
	Type        string
	fp          string
	security    string
	sni         string
	tls         string
	enc         string
	serviceName string
	fingerprint string
	publicKey   string
	serverName  string
	shortId     string
	spiderX     string
	flow        string
	alpn        string
	method      string
	password    string
}

type SECURYTY struct {
	t               string
	securitySetting interface{}
}
