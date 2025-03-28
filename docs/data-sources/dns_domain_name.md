---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "civo_dns_domain_name Data Source - terraform-provider-civo"
subcategory: ""
description: |-
  Get information on a domain. This data source provides the name and the id.
  An error will be raised if the provided domain name is not in your Civo account.
---

# civo_dns_domain_name (Data Source)

Get information on a domain. This data source provides the name and the id.

An error will be raised if the provided domain name is not in your Civo account.

## Example Usage

```terraform
data "civo_dns_domain_name" "domain" {
    name = "domain.com"
}

output "domain_output" {
  value = data.civo_dns_domain_name.domain.name
}

output "domain_id_output" {
  value = data.civo_dns_domain_name.domain.id
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- **id** (String) The ID of this resource.
- **name** (String) The name of the domain


