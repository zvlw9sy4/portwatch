// Package filter implements port ignore rules for portwatch.
//
// Rules are loaded from a YAML file and applied to scanner results
// before diffing, allowing operators to suppress known ports (e.g.
// SSH on 22, HTTP on 80) from triggering alerts.
//
// Example YAML:
//
//	ignore:
//	  - port: 22
//	    protocol: tcp
//	    comment: SSH always open
//	  - port: 80
//	    protocol: tcp
package filter
