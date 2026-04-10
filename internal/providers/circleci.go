package providers

var circleciConfigs = []string{
	".circleci/config.yml",
	".circleci/config.yaml",
}

func NewCircleCIResolver(registryToken string) *imageOnlyResolver {
	return &imageOnlyResolver{
		providerName: "CircleCI",
		matcher:      func(p string) bool { return matchesAny(p, circleciConfigs) },
		docker:       newDockerResolver(registryToken),
	}
}
