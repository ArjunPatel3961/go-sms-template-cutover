package main

import (
	"fmt"
	"strings"
)

type Release struct{ Name, Signature, Template string }

func ValidateRelease(r Release) error {
	for label, value := range map[string]string{"name": r.Name, "signature": r.Signature, "template": r.Template} {
		if strings.TrimSpace(value) == "" {
			return fmt.Errorf("%s is required", label)
		}
	}
	return nil
}

func PublishRelease(c *InfraiClient, r Release) error {
	if err := ValidateRelease(r); err != nil {
		return err
	}
	ns := "sms-cutover-" + strings.ToLower(strings.ReplaceAll(r.Name, " ", "-"))
	if err := c.CreateSignature(ns + "-signature", r.Signature); err != nil {
		return err
	}
	return c.CreateTemplate(ns + "-template", ns + "-template", r.Template)
}
