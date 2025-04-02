package main

type Config struct {
	PineType     string   `yaml:"type"`
	Addr         string   `yaml:"addr"`
	ExcludeFiles []string `yaml:"exclude_files"`
	ExcludeDirs  []string `yaml:"exclude_dirs"`
	IncludeDirs  []string `yaml:"include_dirs"`
	TempDirName  string   `yaml:"temp_dir"`
}
