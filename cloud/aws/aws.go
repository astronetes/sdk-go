package aws

type Client interface {
	Acm() ACMClient
	Route53() Route53Client
}

type client struct {
	acmClient     ACMClient
	route53Client Route53Client
}

func (c *client) Acm() ACMClient {
	return c.acmClient
}

func (c *client) Route53() Route53Client {
	return c.route53Client
}
