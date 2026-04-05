package scanner

const bitbucketPipelinesBase = "bitbucket-pipelines"

func newBitbucketResolver() *imageOnlyResolver {
	return &imageOnlyResolver{
		providerName: "Bitbucket Pipelines",
		matcher: func(p string) bool {
			return p == bitbucketPipelinesBase+".yml" || p == bitbucketPipelinesBase+".yaml"
		},
		docker: newDockerResolver(""),
	}
}
