package providers

import "strings"

// composeConfigs lists exact base file names for Compose spec files.
// Files matching docker-compose.*.yml/yaml (e.g. docker-compose.prod.yml) are
// also matched via the HasPrefix check below.
var composeConfigs = []string{
	"docker-compose.yml",
	"docker-compose.yaml",
	"compose.yml",
	"compose.yaml",
}

// NewComposeResolver returns a provider that pins Docker image: tags in
// docker-compose and Compose spec files.
func NewComposeResolver() *imageOnlyResolver {
	return &imageOnlyResolver{
		providerName: "Docker Compose",
		matcher: func(p string) bool {
			base := slashBase(p)
			return matchesAny(base, composeConfigs) ||
				(strings.HasPrefix(base, "docker-compose.") && isYAML(base))
		},
		docker: newDockerResolver(""),
	}
}
