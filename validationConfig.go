package webhook

import (
	"time"

	"github.com/xmidt-org/urlegit"
)

type ValidatorConfig struct {
	URL     URLVConfig
	TTL     TTLVConfig
	Options OptionsConfig
}

type URLVConfig struct {
	HTTPSOnly            bool
	AllowLoopback        bool
	AllowIP              bool
	AllowSpecialUseHosts bool
	AllowSpecialUseIPs   bool
	InvalidHosts         []string
	InvalidSubnets       []string
}

type TTLVConfig struct {
	Max    time.Duration
	Jitter time.Duration
	Now    func() time.Time
}

type OptionsConfig struct {
	AtLeastOneEvent                bool
	EventRegexMustCompile          bool
	DeviceIDRegexMustCompile       bool
	ValidateRegistrationDuration   bool
	ProvideReceiverURLValidator    bool
	ProvideFailureURLValidator     bool
	ProvideAlternativeURLValidator bool
}

// BuildURLChecker translates the configuration into url Checker to be run on the registration.
var (
	SpecialUseIPs = []string{
		"0.0.0.0/8",          //local ipv4
		"fe80::/10",          //local ipv6
		"255.255.255.255/32", //broadcast to neighbors
		"2001::/32",          //ipv6 TEREDO prefix
		"2001:5::/32",        //EID space for lisp
		"2002::/16",          //ipv6 6to4
		"fc00::/7",           //ipv6 unique local
		"192.0.0.0/24",       //ipv4 IANA
		"2001:0000::/23",     //ipv6 IANA
		"224.0.0.1/32",       //ipv4 multicast
	}
	SpecialUseHosts = []string{
		".example.",
		".invalid.",
		".test.",
		"localhost",
	}
)

// BuildURLChecker translates the configuration into url Checker to be run on the webhook.
func BuildURLChecker(config ValidatorConfig) (*urlegit.Checker, error) {
	var o []urlegit.Option
	if config.URL.HTTPSOnly {
		o = append(o, urlegit.OnlyAllowSchemes("https"))
	}
	if !config.URL.AllowLoopback {
		o = append(o, urlegit.ForbidLoopback())
	}
	if !config.URL.AllowIP {
		o = append(o, urlegit.ForbidAnyIPs())
	}
	if !config.URL.AllowSpecialUseHosts {
		o = append(o, urlegit.ForbidSpecialUseDomains())
	}
	if !config.URL.AllowSpecialUseIPs {
		o = append(o, urlegit.ForbidSubnets(SpecialUseIPs))
	}
	checker, err := urlegit.New(o...)
	if err != nil {
		return nil, err
	}
	return checker, nil
}

//BuildOptions translates the configuration into a list of options to be used to validate the registration
func BuildOptions(config ValidatorConfig, checker *urlegit.Checker) []Option {
	var opts []Option
	if config.Options.AtLeastOneEvent {
		opts = append(opts, AtLeastOneEvent())
	}
	if config.Options.EventRegexMustCompile {
		opts = append(opts, EventRegexMustCompile())
	}
	if config.Options.DeviceIDRegexMustCompile {
		opts = append(opts, DeviceIDRegexMustCompile())
	}
	if config.Options.ProvideReceiverURLValidator {
		opts = append(opts, ProvideReceiverURLValidator(checker))
	}
	if config.Options.ProvideFailureURLValidator {
		opts = append(opts, ProvideFailureURLValidator(checker))
	}
	if config.Options.ProvideAlternativeURLValidator {
		opts = append(opts, ProvideAlternativeURLValidator(checker))
	}
	return opts
}
