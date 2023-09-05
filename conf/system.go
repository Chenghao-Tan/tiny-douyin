package conf

type System struct {
	TrustedProxy  string `yaml:"trustedProxy"`
	ListenAddress string `yaml:"listenAddress"`
	ListenPort    string `yaml:"listenPort"`
	AutoTLS       string `yaml:"autoTLS"`
	FFmpeg        string `yaml:"ffmpeg"`
	TempDir       string `yaml:"tempDir"`
	AutoLogout    int    `yaml:"autoLogout"`
	RateLimit     int    `yaml:"rateLimit"`
}
