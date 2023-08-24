package conf

type System struct {
	HttpPort   string `yaml:"httpPort"`
	AutoTLS    string `yaml:"autoTLS"`
	FFmpeg     string `yaml:"ffmpeg"`
	TempDir    string `yaml:"tempDir"`
	AutoLogout int    `yaml:"autoLogout"`
	Capacity   int    `yaml:"capacity"`
	Recover    int    `yaml:"recover"`
}
