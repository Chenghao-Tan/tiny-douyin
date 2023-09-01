package conf

type Cache struct {
	SyncInterval        int     `yaml:"syncInterval"`
	MaxWriteTime        int     `yaml:"maxWriteTime"`
	CacheExpiration     int     `yaml:"cacheExpiration"`
	EmptyExpiration     int     `yaml:"emptyExpiration"`
	DistrustProbability float32 `yaml:"distrustProbability"`
}
