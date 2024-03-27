package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
	"k8s.io/apimachinery/pkg/api/errors"
)

func (api *ApiImpl) CreateRegistrarCluster(c echo.Context) error {
	logger := api.logger
	if !api.rbac.IsAllowed("TODO username", "create", "managementcluster") {
		return c.JSON(http.StatusForbidden, nil)
	}

	if c.Request().Body == nil {
		logger.Error(fmt.Errorf("create cluster request error"), "request body is nil")
		return c.JSON(http.StatusBadRequest, nil)
	}

	req, err := io.ReadAll(c.Request().Body)
	defer c.Request().Body.Close()
	if err != nil {
		logger.Error(err, "failed to read request body on create cluster")
		return c.JSON(http.StatusInternalServerError, nil)
	}

	cluster := &RegistrarCluster{}
	if err := json.Unmarshal(req, cluster); err != nil {
		logger.Error(err, "failed to unmarshal request body on create cluster")
		return c.JSON(http.StatusBadRequest, nil)
	}

	if cluster.Kubeconfig == nil {
		return c.JSON(http.StatusBadRequest, fmt.Errorf("kubeconfig field is required"))
	}

	if cluster.Name == "" {
		return c.JSON(http.StatusBadRequest, fmt.Errorf("name field is required"))
	}

	regCluster, err := api.manager.GetClusterRegistrar().Create(&ClusterRegistrar{
		Name:       cluster.Name,
		Region:     cluster.Region,
		Kubeconfig: *cluster.Kubeconfig,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, nil)
	}

	return c.JSON(http.StatusCreated, CreateRegistrarClusterResponse{
		JSON201: &RegistrarCluster{
			Id:         &regCluster.Id,
			Name:       regCluster.Name,
			Kubeconfig: &regCluster.Kubeconfig,
			Region:     regCluster.Region,
		},
	})
}

func (api *ApiImpl) ListRegistrarClusters(c echo.Context) error {
	if !api.rbac.IsAllowed("TODO username", "list", "managementcluster") {
		return c.JSON(http.StatusForbidden, nil)
	}

	clusters, err := api.manager.GetClusterRegistrar().List()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, nil)
	}

	resp := ListRegistrarClustersResponse{
		JSON200: &struct {
			Clusters []RegistrarCluster "json:\"clusters\""
			Warnings *[]string          "json:\"warnings,omitempty\""
		}{
			Clusters: []RegistrarCluster{},
			Warnings: nil,
		},
	}

	for _, cluster := range clusters {
		resp.JSON200.Clusters = append(resp.JSON200.Clusters, RegistrarCluster{
			Id:         &cluster.Id,
			Name:       cluster.Name,
			Kubeconfig: &cluster.Kubeconfig,
			Region:     cluster.Region,
		})
	}

	return c.JSON(http.StatusOK, resp)
}

func (api *ApiImpl) GetRegistrarCluster(c echo.Context, region string) error {
	if !api.rbac.IsAllowed("TODO username", "get", "managementcluster") {
		return c.JSON(http.StatusForbidden, nil)
	}

	cluster, err := api.manager.GetClusterRegistrar().Get(region)
	if err != nil {
		if errors.IsNotFound(err) {
			return c.JSON(http.StatusNotFound, nil)
		}

		return c.JSON(http.StatusInternalServerError, nil)
	}

	return c.JSON(http.StatusOK, GetRegistrarClusterResponse{
		JSON200: &RegistrarCluster{
			Id:         &cluster.Id,
			Name:       cluster.Name,
			Kubeconfig: &cluster.Kubeconfig,
			Region:     cluster.Region,
		},
	})
}

func (api *ApiImpl) DeleteRegistrarCluster(c echo.Context, id string) error {
	if !api.rbac.IsAllowed("TODO username", "delete", "managementcluster") {
		return c.JSON(http.StatusForbidden, nil)
	}

	if err := api.manager.GetClusterRegistrar().Delete(id); err != nil {
		if errors.IsNotFound(err) {
			return c.JSON(http.StatusNotFound, nil)
		}

		return c.JSON(http.StatusInternalServerError, nil)
	}

	return c.JSON(http.StatusNoContent, nil)
}