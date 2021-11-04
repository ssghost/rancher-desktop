package mungers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/sirupsen/logrus"

	"github.com/rancher-sandbox/rancher-desktop/src/wsl-helper/pkg/dockerproxy"
	"github.com/rancher-sandbox/rancher-desktop/src/wsl-helper/pkg/dockerproxy/models"
	"github.com/rancher-sandbox/rancher-desktop/src/wsl-helper/pkg/dockerproxy/platform"
)

type containersCreateBody struct {
	models.ContainerConfig
	HostConfig       models.HostConfig
	NetworkingConfig models.NetworkingConfig
}

// munge POST /containers/create to use WSL paths
func mungeContainersCreate(req *http.Request, contextValue *dockerproxy.RequestContextValue, templates map[string]string) error {
	body := containersCreateBody{}
	err := readRequestBodyJSON(req, &body)
	if err != nil {
		return err
	}
	logrus.WithField("body", fmt.Sprintf("%+v", body)).Debug("read body")

	modified := false
	for bindIndex, bind := range body.HostConfig.Binds {
		logrus.WithField(fmt.Sprintf("bind %d", bindIndex), bind).Debug("got bind")
		host, container, options, isPath := platform.ParseBindString(bind)
		if isPath {
			translated, err := platform.TranslatePathFromClient(host)
			if err != nil {
				return fmt.Errorf("could not translate mount path %s: %w", host, err)
			}
			host = translated
			modified = true
		}
		if options == "" {
			body.HostConfig.Binds[bindIndex] = fmt.Sprintf("%s:%s", host, container)
		} else {
			body.HostConfig.Binds[bindIndex] = fmt.Sprintf("%s:%s:%s", host, container, options)
		}
	}
	if !modified {
		return nil
	}

	buf, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("could not re-marshel parameters: %w", err)
	}
	req.Body = io.NopCloser(bytes.NewBuffer(buf))
	req.ContentLength = int64(len(buf))
	req.Header.Set("Content-Length", fmt.Sprintf("%d", len(buf)))

	return nil
}

func init() {
	dockerproxy.RegisterRequestMunger(http.MethodPost, "/containers/create", mungeContainersCreate)
}
