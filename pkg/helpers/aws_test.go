package helpers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetAPIGatewayManagementAPIClient(t *testing.T) {
	api := GetAPIGatewayManagementAPIClient()
	assert.NotNil(t, api)
}
