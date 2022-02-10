package greet

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_service_Hello(t *testing.T) {
	helloParam := &HelloParam{"Tom"}

	helloResult := fmt.Sprintf("Hello, %s!", helloParam.Name)

	newService := New()
	result, err := newService.Hello(context.TODO(), helloParam)

	assert.NoError(t, err)
	assert.Equal(t, helloResult, result)
}
