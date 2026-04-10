package providers

var bitbucketConfigs = []string{
	"bitbucket-pipelines.yml",
	"bitbucket-pipelines.yaml",
}

func NewBitbucketResolver() *imageOnlyResolver {
	return &imageOnlyResolver{
		providerName: "Bitbucket Pipelines",
		matcher:      func(p string) bool { return matchesAny(slashBase(p), bitbucketConfigs) },
		docker:       newDockerResolver(""),
	}
}
