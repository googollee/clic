package sources

var Default = []Source{
	Flag(FlagPrefix(""), FlagSplitter(".")),
	File(FilePathFlag("config", ""), FileFormat(JSON{})),
	Env(EnvPrefix("CLIC"), EnvSplitter("_")),
}
