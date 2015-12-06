package service
import (
	"testing"
	"errors"
	"github.com/stretchr/testify/assert"
)

func TestAggregateError(t *testing.T) {
	agg := NewAggregateError()
	agg.AddError(errors.New("blah"))
	assert.False(t, agg.Empty())
	expected := "blah\n"
	assert.Equal(t, expected, agg.Error())
	assert.Equal(t, expected, agg.String())
}
