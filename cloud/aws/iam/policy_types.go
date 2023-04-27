package iam

type CreatePolicyRequest struct {
	Name        string
	Description string
	Document    PolicyDocument
}

type PolicyDocument struct {
	Version   string
	Statement []PolicyStatement
}

type PolicyStatement struct {
	Effect    string
	Action    []string
	Principal map[string]interface{} `json:",omitempty"`
	Resource  *string                `json:",omitempty"`
	Condition *PolicyCondition       `json:",omitempty"`
}

type PolicyCondition struct {
	StringLike map[string][]string `json:",omitempty"`
}
