// SPDX-FileCopyrightText: 2025 Comcast Cable Communications Management, LLC
// SPDX-License-Identifier: Apache-2.0

package option

import (
	"time"

	"github.com/xmidt-org/urlegit"
)

// URLChecker provides options for validating the stream's URL and TTL
// related fields.
type URLChecker struct {
	URL    URLConfig    `json:"url"`
	TTL    TTLConfig    `json:"ttl"`
	IP     IPConfig     `json:"ip"`
	Domain DomainConfig `json:"domain"`
}

type IPConfig struct {
	Allow            bool     `json:"allow"`
	ForbiddenSubnets []string `json:"forbidden_subnets"`
}

type DomainConfig struct {
	AllowSpecialUseDomains bool     `json:"allow_special_use_domains"`
	ForbiddenDomains       []string `json:"forbidden_domains"`
}

type URLConfig struct {
	Schemes       []string `json:"schemes"`
	AllowLoopback bool     `json:"allow_loop_back"`
}

type TTLConfig struct {
	Max    time.Duration `json:"max"`
	Jitter time.Duration `json:"jitter"`
}

// Build translates the configuration into url Checker to be run on the stream.
func (opts *URLChecker) Build() (*urlegit.Checker, error) {
	var o []urlegit.Option
	if len(opts.URL.Schemes) > 0 {
		o = append(o, urlegit.OnlyAllowSchemes(opts.URL.Schemes...))
	}
	if !opts.URL.AllowLoopback {
		o = append(o, urlegit.ForbidLoopback())
	}
	if !opts.IP.Allow {
		o = append(o, urlegit.ForbidAnyIPs())
	}
	if len(opts.IP.ForbiddenSubnets) > 0 {
		o = append(o, urlegit.ForbidSubnets(opts.IP.ForbiddenSubnets))
	}
	if !opts.Domain.AllowSpecialUseDomains {
		o = append(o, urlegit.ForbidSpecialUseDomains())
	}
	if len(opts.Domain.ForbiddenDomains) > 0 {
		o = append(o, urlegit.ForbidDomainNames(opts.Domain.ForbiddenDomains...))
	}
	return urlegit.New(o...)
}
