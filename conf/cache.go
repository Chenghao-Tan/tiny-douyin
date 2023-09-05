package conf

type Cache struct {
	SyncInterval        int     `yaml:"syncInterval"`
	MaxRWTime           int     `yaml:"maxRWTime"`
	CacheExpiration     int     `yaml:"cacheExpiration"`
	EmptyExpiration     int     `yaml:"emptyExpiration"`
	DistrustProbability float32 `yaml:"distrustProbability"`
}
