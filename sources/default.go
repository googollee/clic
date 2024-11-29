package sources

var Default = []Source{
	File(FilePathFlag("config", ""), FileFormat(JSON{})),
	Flag(FlagPrefix(""), FlagSplitter(".")),
	Env(EnvPrefix("CLIC"), EnvSplitter("_")),
}
