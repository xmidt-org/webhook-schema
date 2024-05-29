package webhook

import (
	"time"

	"github.com/xmidt-org/urlegit"
)

type Validator interface {
	ValidateOneEvent() error
	ValidateEventRegex() error
	ValidateDeviceId() error
	ValidateUntil(time.Duration, time.Duration, func() time.Time) error
	ValidateNoUntil() error
	ValidateDuration(time.Duration) error
	ValidateFailureURL(*urlegit.Checker) error
	ValidateReceiverURL(*urlegit.Checker) error
	ValidateAltURL(*urlegit.Checker) error
	SetNowFunc(func() time.Time)
}

type ValidatorConfig struct {
	URL URLVConfig
	TTL TTLVConfig
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
	// errFailedToBuildValidators    = errors.New("failed to build validators")
	// errFailedToBuildValidURLFuncs = errors.New("failed to build ValidURLFuncs")
)

// BuildURLChecker translates the configuration into url Checker to be run on the webhook.
func buildURLChecker(config ValidatorConfig) (*urlegit.Checker, error) {
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

// BuildValidators translates the configuration into a list of validators to be run on the
// webhook.
func BuildValidators(config ValidatorConfig) ([]Option, error) {
	var opts []Option

	checker, err := buildURLChecker(config)
	if err != nil {
		return nil, err
	}
	opts = append(opts,
		AtLeastOneEvent(),
		EventRegexMustCompile(),
		DeviceIDRegexMustCompile(),
		ValidateRegistrationDuration(config.TTL.Max),
		ProvideReceiverURLValidator(checker),
		ProvideFailureURLValidator(checker),
		ProvideAlternativeURLValidator(checker),
	)
	return opts, nil
}