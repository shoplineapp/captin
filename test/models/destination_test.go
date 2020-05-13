package models_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	. "github.com/shoplineapp/captin/models"
)

func TestGetCallbackURL(t *testing.T) {
	var config Configuration
	var subject Destination

	config = Configuration{Name: "site_a", CallbackURL: "http://site-a.com/callback"}
	subject = Destination{Config: config}
	assert.Equal(t, subject.GetCallbackURL(), config.CallbackURL)

	overriden := "http://google.com"
	os.Setenv("HOOK_SITE_B_CALLBACK_URL", overriden)
	config = Configuration{Name: "site_b", CallbackURL: "http://site-b.com/callback"}
	subject = Destination{Config: config}
	assert.Equal(t, subject.GetCallbackURL(), overriden)
}
