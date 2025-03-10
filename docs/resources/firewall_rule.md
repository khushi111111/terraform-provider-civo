---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "civo_firewall_rule Resource - terraform-provider-civo"
subcategory: ""
description: |-
  Provides a Civo firewall rule resource. This can be used to create, modify, and delete firewalls rules. This resource don't have an update option because Civo backend doesn't support it at this moment. In that case, we use ForceNew for all object in the resource.
---

# civo_firewall_rule (Resource)

Provides a Civo firewall rule resource. This can be used to create, modify, and delete firewalls rules. This resource don't have an update option because Civo backend doesn't support it at this moment. In that case, we use `ForceNew` for all object in the resource.

## Example Usage

```terraform
# Query small instance size
data "civo_instances_size" "small" {
    filter {
        key = "name"
        values = ["g3.small"]
        match_by = "re"
    }

    filter {
        key = "type"
        values = ["instance"]
    }

}

# Query instance disk image
data "civo_disk_image" "debian" {
   filter {
        key = "name"
        values = ["debian-10"]
   }
}

# Create a new instance
resource "civo_instance" "foo" {
    hostname = "foo.com"
    size = element(data.civo_instances_size.small.sizes, 0).name
    disk_image = element(data.civo_disk_image.debian.diskimages, 0).id
}

# Create a network
resource "civo_network" "custom_net" {
    label = "my-custom-network"
}

# Create a firewall
resource "civo_firewall" "custom_firewall" {
  name = "my-custom-firewall"
  network_id = civo_network.custom_net.id
}

# Create a firewall rule and only allow
# connections from instance we created above
resource "civo_firewall_rule" "custom_port" {
    firewall_id = civo_firewall.custom_firewall.id
    protocol = "tcp"
    start_port = "3000"
    end_port = "3000"
    cidr = [format("%s/%s",civo_instance.foo.public_ip,"32")]
    direction = "ingress"
    label = "custom-application"
    depends_on = [civo_firewall.custom_firewall]
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- **action** (String) the action of the rule can be allow or deny
- **cidr** (Set of String) The CIDR notation of the other end to affect, or a valid network CIDR (e.g. 0.0.0.0/0 to open for everyone or 1.2.3.4/32 to open just for a specific IP address)
- **direction** (String) The direction of the rule can be ingress or egress
- **firewall_id** (String) The Firewall ID

### Optional

- **end_port** (String) The end of the port range (this is optional, by default it will only apply to the single port listed in start_port)
- **id** (String) The ID of this resource.
- **label** (String) A string that will be the displayed name/reference for this rule
- **protocol** (String) The protocol choice from `tcp`, `udp` or `icmp` (the default if unspecified is `tcp`)
- **region** (String) The region for this rule
- **start_port** (String) The start of the port range to configure for this rule (or the single port if required)

## Import

Import is supported using the following syntax:

```shell
# using firewall_id:firewall_rule_id
terraform import civo_firewall_rule.http b8ecd2ab-2267-4a5e-8692-cbf1d32583e3:4b0022ee-00b2-4f81-a40d-b4f8728923a7
```
