package webhook

import (
	"errors"
	"fmt"
	"regexp"
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

type Option interface {
	fmt.Stringer
	Validate(Validator) error
}

// Validate is a method on Registration that validates the registration
// against a list of options.
func Validate(v Validator, opts []Option) error {
	var errs error
	for _, opt := range opts {
		if opt != nil {
			if err := opt.Validate(v); err != nil {
				errs = errors.Join(errs, err)
			}
		}
	}
	return errs
}

func (v1 *RegistrationV1) ValidateOneEvent() error {
	if len(v1.Events) == 0 {
		return fmt.Errorf("%w: cannot have zero events", ErrInvalidInput)
	}
	return nil
}

func (v1 *RegistrationV1) ValidateEventRegex() error {
	var errs error
	for _, e := range v1.Events {
		_, err := regexp.Compile(e)
		if err != nil {
			errs = errors.Join(errs, fmt.Errorf("%w: unable to compile matching", ErrInvalidInput))
		}
	}
	return errs
}

func (v1 *RegistrationV1) ValidateDeviceId() error {
	var errs error
	for _, e := range v1.Matcher.DeviceID {
		_, err := regexp.Compile(e)
		if err != nil {
			errs = errors.Join(errs, fmt.Errorf("%w: unable to compile matching", ErrInvalidInput))
		}
	}
	return errs
}

func (v1 *RegistrationV1) ValidateDuration(ttl time.Duration) error {
	var errs error
	if ttl <= 0 {
		ttl = time.Duration(0)
	}

	if ttl != 0 && ttl < time.Duration(v1.Duration) {
		errs = errors.Join(errs, fmt.Errorf("%w: the registration is for too long", ErrInvalidInput))
	}

	if v1.Until.IsZero() && v1.Duration == 0 {
		errs = errors.Join(errs, fmt.Errorf("%w: either Duration or Until must be set", ErrInvalidInput))
	}

	if !v1.Until.IsZero() && v1.Duration != 0 {
		errs = errors.Join(errs, fmt.Errorf("%w: only one of Duration or Until may be set", ErrInvalidInput))
	}

	if !v1.Until.IsZero() {
		nowFunc := time.Now
		if v1.nowFunc != nil {
			nowFunc = v1.nowFunc
		}

		now := nowFunc()
		if ttl != 0 && v1.Until.After(now.Add(ttl)) {
			errs = errors.Join(errs, fmt.Errorf("%w: the registration is for too long", ErrInvalidInput))
		}

		if v1.Until.Before(now) {
			errs = errors.Join(errs, fmt.Errorf("%w: the registration has already expired", ErrInvalidInput))
		}
	}

	return errs
}

func (v1 *RegistrationV1) ValidateFailureURL(c *urlegit.Checker) error {
	if v1.FailureURL != "" {
		if err := c.Text(v1.FailureURL); err != nil {
			return fmt.Errorf("%w: failure url is invalid", err)
		}
	}
	return nil
}

func (v1 *RegistrationV1) ValidateReceiverURL(c *urlegit.Checker) error {
	if v1.Config.ReceiverURL != "" {
		if err := c.Text(v1.Config.ReceiverURL); err != nil {
			return fmt.Errorf("%w: receiver url is invalid", ErrInvalidInput)
		}
	}
	return nil
}

func (v1 *RegistrationV1) ValidateAltURL(c *urlegit.Checker) error {
	var errs error
	for _, url := range v1.Config.AlternativeURLs {
		if err := c.Text(url); err != nil {
			errs = errors.Join(errs, fmt.Errorf("%w: alternative url is invalid: %v", ErrInvalidInput, url))
		}
	}
	return errs
}

func (v1 *RegistrationV1) ValidateNoUntil() error {
	if !v1.Until.IsZero() {
		return fmt.Errorf("%w: Until is not allowed", ErrInvalidInput)
	}
	return nil
}

func (v1 *RegistrationV1) ValidateUntil(jitter time.Duration, maxTTL time.Duration, now func() time.Time) error {
	if now == nil {
		now = time.Now
	}
	if maxTTL < 0 {
		return ErrInvalidInput
	} else if jitter < 0 {
		return ErrInvalidInput
	}

	if v1.Until.IsZero() {
		return nil
	}
	limit := (now().Add(maxTTL)).Add(jitter)
	proposed := (v1.Until)
	if proposed.After(limit) {
		return fmt.Errorf("%w: %v after %v",
			ErrInvalidInput, proposed.String(), limit.String())
	}
	return nil

}

func (v1 *RegistrationV1) SetNowFunc(now func() time.Time) {
	v1.nowFunc = now
}

func (v2 *RegistrationV2) ValidateOneEvent() error {
	// if len(v2.) == 0 {
	// 	return fmt.Errorf("%w: cannot have zero events", ErrInvalidInput)
	// }
	return nil
}

func (v2 *RegistrationV2) ValidateEventRegex() error {
	// var errs error
	// for _, e := range v1.Events {
	// 	_, err := regexp.Compile(e)
	// 	if err != nil {
	// 		errs = errors.Join(errs, fmt.Errorf("%w: unable to compile matching", ErrInvalidInput))
	// 	}
	// }
	return nil
}

func (v2 *RegistrationV2) ValidateDeviceId() error {
	// var errs error
	// for _, e := range v2.Matcher {
	// 	_, err := regexp.Compile(e)
	// 	if err != nil {
	// 		errs = errors.Join(errs, fmt.Errorf("%w: unable to compile matching", ErrInvalidInput))
	// 	}
	// }
	return nil
}

func (v2 *RegistrationV2) ValidateFailureURL(c *urlegit.Checker) error {
	if v2.FailureURL != "" {
		if err := c.Text(v2.FailureURL); err != nil {
			return fmt.Errorf("%w: failure url is invalid", err)
		}
	}
	return nil
}

func (v2 *RegistrationV2) ValidateReceiverURL(c *urlegit.Checker) error {
	// if v2.Config.ReceiverURL != "" {
	// 	if err := c.Text(v1.Config.ReceiverURL); err != nil {
	// 		return fmt.Errorf("%w: receiver url is invalid", ErrInvalidInput)
	// 	}
	// }
	return nil
}

func (v2 *RegistrationV2) ValidateAltURL(c *urlegit.Checker) error {
	// var errs error
	// for _, webhook := range v2.Webhooks{
	// 	for _, url := range webhook.ReceiverURLs {
	// 		if err := c.Text(url); err != nil {
	// 			errs = errors.Join(errs, fmt.Errorf("%w: url is invalid", ErrInvalidInput))
	// 		}
	// 	}
	// }

	// return errs
	return nil
}

func (v2 *RegistrationV2) SetNowFunc(now func() time.Time) {

}

func (v2 *RegistrationV2) ValidateDuration(ttl time.Duration) error {
	// var errs error
	// if ttl <= 0 {
	// 	ttl = time.Duration(0)
	// }

	// if ttl != 0 && ttl < time.Duration(v1.Duration) {
	// 	errs = errors.Join(errs, fmt.Errorf("%w: the registration is for too long", ErrInvalidInput))
	// }

	// if v1.Until.IsZero() && v1.Duration == 0 {
	// 	errs = errors.Join(errs, fmt.Errorf("%w: either Duration or Until must be set", ErrInvalidInput))
	// }

	// if !v1.Until.IsZero() && v1.Duration != 0 {
	// 	errs = errors.Join(errs, fmt.Errorf("%w: only one of Duration or Until may be set", ErrInvalidInput))
	// }

	// if !v1.Until.IsZero() {
	// 	nowFunc := time.Now
	// 	if v1.nowFunc != nil {
	// 		nowFunc = v1.nowFunc
	// 	}

	// 	now := nowFunc()
	// 	if ttl != 0 && v1.Until.After(now.Add(ttl)) {
	// 		errs = errors.Join(errs, fmt.Errorf("%w: the registration is for too long", ErrInvalidInput))
	// 	}

	// 	if v1.Until.Before(now) {
	// 		errs = errors.Join(errs, fmt.Errorf("%w: the registration has already expired", ErrInvalidInput))
	// 	}
	// }

	// return errs
	return nil
}

func (v2 *RegistrationV2) ValidateNoUntil() error {
	// if !v1.Until.IsZero() {
	// 	return fmt.Errorf("%w: Until is not allowed", ErrInvalidInput)
	// }
	return nil
}

func (v2 *RegistrationV2) ValidateUntil(jitter time.Duration, maxTTL time.Duration, now func() time.Time) error {
	// if now == nil {
	// 	now = time.Now
	// }
	// if maxTTL < 0 {
	// 	return ErrInvalidInput
	// } else if jitter < 0 {
	// 	return ErrInvalidInput
	// }

	// if v1.Until.IsZero() {
	// 	return nil
	// }
	// limit := (now().Add(maxTTL)).Add(jitter)
	// proposed := (v1.Until)
	// if proposed.After(limit) {
	// 	return fmt.Errorf("%w: %v after %v",
	// 		ErrInvalidInput, proposed.String(), limit.String())
	// }
	return nil

}
