package providers

const (
	bitbucketPipelinesYML  = "bitbucket-pipelines.yml"
	bitbucketPipelinesYAML = "bitbucket-pipelines.yaml"
)

func NewBitbucketResolver() *imageOnlyResolver {
	return &imageOnlyResolver{
		providerName: "Bitbucket Pipelines",
		matcher: func(p string) bool {
			return matchesAny(slashBase(p),
				bitbucketPipelinesYML,
				bitbucketPipelinesYAML,
			)
		},
		docker: newDockerResolver(""),
	}
}
