package resources

// DefaultShortOutputTemplate defines default template for the list output
var DefaultShortOutputTemplate = "{{range .Tags}}{{if eq (.Key | String) \"Name\"}}{{.Value | String | printf \"%-60s\"}}{{end}}{{end}}   {{range .Tags}}{{if eq (.Key | String) \"environment\"}}{{.Value | String | printf \"%-10s\"}}{{end}}{{end}}   {{range .Tags}}{{if eq (.Key | String) \"project\"}}{{.Value | String | printf \"%-10s\"}}{{end}}{{end}} {{.Placement.AvailabilityZone}}    {{.PrivateIpAddress| String | printf \"%-20s\"}} {{.PublicIpAddress| String |printf \"%-20s\"}}\n"

// DefaultLongOutputTemplate defines default template for the detailed output
var DefaultLongOutputTemplate = `
{{range .Tags}}{{if eq (.Key | String) "Name"}}{{.Value}}{{end}}{{end}}
  ami_id              {{.ImageId}}
  az                  {{.Placement.AvailabilityZone}}
  dns_name            {{.PublicDnsName}}
  id                  {{.InstanceId}}
  instance_type       {{.InstanceType}}
  key_name            {{.KeyName}}
  private_dns_name    {{.PrivateDnsName}}
  private_ip          {{.PrivateIpAddress}}
  public_ip           {{.PublicIpAddress}}
  tags
    {{range .Tags}}{{.Key | String | printf "%-10s"}}        {{.Value}}
    {{end}}
`
