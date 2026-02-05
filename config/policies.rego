package sapliy.authz

import rego.v1

# Default deny
default allow := false

# Allow if role has wildcard permission
allow if {
	some role in input.roles
	roles[role].permissions[_] == "*"
}

# Allow if role has specific permission
allow if {
	some role in input.roles
	some permission in roles[role].permissions
	permission == input.action
}

# Example roles data (can be provided externally or defined here for defaults)
# In a real system, this 'data' object would often be loaded from JSON or a DB.
# For this phase, we'll embed the defaults if not provided in the data store.
roles := {
	"admin": {"permissions": ["*"]},
	"finance": {"permissions": [
		"payment.create",
		"refund.create",
		"flow.deploy",
		"flow.deploy.live"
	]},
	"developer": {"permissions": [
		"zone.create",
		"flow.deploy",
		"key.create"
	]},
	"viewer": {"permissions": []}
}
