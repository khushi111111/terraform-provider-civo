package civo

import (
	"context"
	"log"
	"strings"

	"github.com/civo/civogo"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// Data source to get from the api a specific instance
// using the id or the hostname
func dataSourceKubernetesCluster() *schema.Resource {
	return &schema.Resource{
		Description: strings.Join([]string{
			"Provides a Civo Kubernetes cluster data source.",
			"Note: This data source returns a single Kubernetes cluster. When specifying a name, an error will be raised if more than one Kubernetes cluster found.",
		}, "\n\n"),
		ReadContext: dataSourceKubernetesClusterRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.NoZeroValues,
				ExactlyOneOf: []string{"id", "name"},
			},
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.NoZeroValues,
				ExactlyOneOf: []string{"id", "name"},
				Description:  "The name of the Kubernetes Cluster",
			},
			"region": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.NoZeroValues,
				Description:  "The region where cluster is running",
			},
			// computed attributes
			"num_target_nodes": {
				Type:        schema.TypeInt,
				Deprecated:  "This field is deprecated and will be removed in a future version of the provider",
				Computed:    true,
				Description: "The size of the Kubernetes cluster",
			},
			"target_nodes_size": {
				Type:        schema.TypeString,
				Deprecated:  "This field is deprecated and will be removed in a future version of the provider",
				Computed:    true,
				Description: "The size of each node",
			},
			"kubernetes_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The version of Kubernetes",
			},
			"cni": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The cni for the k3s to install (the default is `flannel`) valid options are `cilium` or `flannel`",
			},
			"tags": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A list of tags",
			},
			"applications": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A list of application installed",
			},
			"installed_applications": dataSourceApplicationSchema(),
			"pools":                  dataSourcenodePoolSchema(),
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of Kubernetes cluster",
			},
			"ready": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "If the Kubernetes cluster is ready",
			},
			"kubeconfig": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A representation of the Kubernetes cluster's kubeconfig in yaml format",
			},
			"api_endpoint": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The base URL of the API server on the Kubernetes master node",
			},
			"master_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The IP of the Kubernetes master node",
			},
			"dns_entry": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique dns entry for the cluster in this case point to the master",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date where the Kubernetes cluster was create",
			},
		},
	}
}

// schema for the node pool in the cluster
func dataSourcenodePoolSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"id": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "The ID of the pool",
				},
				"node_count": {
					Type:        schema.TypeInt,
					Computed:    true,
					Description: "The size of the pool",
				},
				"size": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "The size of each node inside the pool",
				},
				"instance_names": {
					Type:        schema.TypeSet,
					Computed:    true,
					Description: "A list of the instance in the pool",
					Elem:        &schema.Schema{Type: schema.TypeString},
				},
			},
		},
	}
}

// schema for the application in the cluster
func dataSourceApplicationSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"application": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "The name of the application",
				},
				"version": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "The version of the application",
				},
				"installed": {
					Type:        schema.TypeBool,
					Computed:    true,
					Description: "If the application is installed, this will return `true`",
				},
				"category": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "The category of the application",
				},
			},
		},
	}
}

func dataSourceKubernetesClusterRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*civogo.Client)

	// overwrite the region if is define in the datasource
	if region, ok := d.GetOk("region"); ok {
		apiClient.Region = region.(string)
	}

	var foundCluster *civogo.KubernetesCluster

	if id, ok := d.GetOk("id"); ok {
		log.Printf("[INFO] Getting the kubernetes Cluster by id")
		kubeCluster, err := apiClient.FindKubernetesCluster(id.(string))
		if err != nil {
			return diag.Errorf("[ERR] failed to retrive kubernetes cluster: %s", err)
		}
		foundCluster = kubeCluster
	} else if name, ok := d.GetOk("name"); ok {
		log.Printf("[INFO] Getting the kubernetes Cluster by name")
		kubeCluster, err := apiClient.FindKubernetesCluster(name.(string))
		if err != nil {
			return diag.Errorf("[ERR] failed to retrive kubernetes cluster: %s", err)
		}

		foundCluster = kubeCluster
	}

	d.SetId(foundCluster.ID)
	d.Set("name", foundCluster.Name)
	d.Set("num_target_nodes", foundCluster.NumTargetNode)
	d.Set("target_nodes_size", foundCluster.TargetNodeSize)
	d.Set("kubernetes_version", foundCluster.KubernetesVersion)
	d.Set("cni", foundCluster.CNIPlugin)
	d.Set("tags", foundCluster.Tags)
	d.Set("status", foundCluster.Status)
	d.Set("ready", foundCluster.Ready)
	d.Set("kubeconfig", foundCluster.KubeConfig)
	d.Set("api_endpoint", foundCluster.APIEndPoint)
	d.Set("master_ip", foundCluster.MasterIP)
	d.Set("dns_entry", foundCluster.DNSEntry)
	d.Set("created_at", foundCluster.CreatedAt.UTC().String())

	if err := d.Set("pools", flattenNodePool(foundCluster)); err != nil {
		return diag.Errorf("[ERR] error retrieving the pools for kubernetes cluster error: %#v", err)
	}

	if err := d.Set("installed_applications", flattenInstalledApplication(foundCluster.InstalledApplications)); err != nil {
		return diag.Errorf("[ERR] error retrieving the installed application for kubernetes cluster error: %#v", err)
	}

	return nil
}
