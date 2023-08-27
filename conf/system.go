package conf

type System struct {
	ListenAddress string `yaml:"listenAddress"`
	ListenPort    string `yaml:"listenPort"`
	AutoTLS       string `yaml:"autoTLS"`
	FFmpeg        string `yaml:"ffmpeg"`
	TempDir       string `yaml:"tempDir"`
	AutoLogout    int    `yaml:"autoLogout"`
	Capacity      int    `yaml:"capacity"`
	Recover       int    `yaml:"recover"`
}
