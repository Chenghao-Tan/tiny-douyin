package conf

type Log struct {
	Path       string `yaml:"path"`       // 输出路径
	Level      string `yaml:"level"`      // 输出级别
	Prefix     string `yaml:"prefix"`     // 日志前缀
	ShowLine   bool   `yaml:"showLine"`   // 显示行号
	MaxSize    int    `yaml:"maxSize"`    // 单个日志文件最大大小
	MaxBackups int    `yaml:"maxBackups"` // 最多保留数量
	MaxAge     int    `yaml:"maxAge"`     // 最多保留天数
	Compress   bool   `yaml:"compress"`   // 是否gzip压缩
}
