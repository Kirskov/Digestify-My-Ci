package scanner

import "net/http"

type bitbucketResolver struct {
	docker *dockerResolver
}

func newBitbucketResolver() *bitbucketResolver {
	return &bitbucketResolver{
		docker: newDockerResolver(""),
	}
}

func (r *bitbucketResolver) Name() string { return "Bitbucket Pipelines" }

const bitbucketPipelinesBase = "bitbucket-pipelines"

func (r *bitbucketResolver) IsMatch(relPath string) bool {
	return relPath == bitbucketPipelinesBase+".yml" || relPath == bitbucketPipelinesBase+".yaml"
}

// Resolve pins Docker image: tags found in Bitbucket Pipelines config files.
// Bitbucket Pipes use semver versioning with no SHA pinning API, so they are left as-is.
func (r *bitbucketResolver) Resolve(content string, _, pinImages bool) (string, error) {
	if !pinImages {
		return content, nil
	}
	return r.docker.resolveImages(content), nil
}

// setClient allows tests to inject a fake HTTP client into the docker resolver.
func (r *bitbucketResolver) setClient(c *http.Client) {
	r.docker.client = c
}
