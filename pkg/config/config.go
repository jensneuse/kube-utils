package config

import "flag"

var (
	NAMESPACE = *flag.String("namespace", "jens-neuse", "override the default namespace")
)
