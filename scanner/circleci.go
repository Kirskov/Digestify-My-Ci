package scanner

import (
	"net/http"
	"sync"
)

type circleciResolver struct {
	mu     sync.Mutex
	docker *dockerResolver
}

func newCircleCIResolver(registryToken string) *circleciResolver {
	return &circleciResolver{
		docker: newDockerResolver(registryToken),
	}
}

func (r *circleciResolver) Name() string { return "CircleCI" }

const circleciConfig = ".circleci/config.yml"

func (r *circleciResolver) IsMatch(relPath string) bool {
	return relPath == circleciConfig || relPath == ".circleci/config.yaml"
}

// Resolve pins Docker image: tags found in CircleCI config files.
// CircleCI orbs use semver and have no SHA pinning API, so they are left as-is.
func (r *circleciResolver) Resolve(content string, _, pinImages bool) (string, error) {
	if !pinImages {
		return content, nil
	}
	return r.docker.resolveImages(content), nil
}

// setClient allows tests to inject a fake HTTP client into the docker resolver.
func (r *circleciResolver) setClient(c *http.Client) {
	r.docker.client = c
}
