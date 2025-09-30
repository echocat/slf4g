package value

import (
	"os"
	"testing"

	nlevel "github.com/echocat/slf4g/native/level"

	"github.com/echocat/slf4g/internal/test/assert"
	"github.com/echocat/slf4g/native/consumer"
)

func Test_NewProvider(t *testing.T) {
	givenTarget := &mockProviderTarget{
		LevelTarget: &mockLevelTarget{},
		ConsumerTarget: &mockConsumerTarget{
			consumer: consumer.NewWriter(os.Stderr),
		},
	}

	actual := NewProvider(givenTarget)

	assert.ToBeNotNil(t, actual)
	assert.ToBeSame(t, givenTarget, actual.Level.Target)
	assert.ToBeSame(t, givenTarget.GetConsumer(), actual.Consumer.Formatter.Target)
}

func Test_NewProvider_customize(t *testing.T) {
	givenTarget := &mockProviderTarget{
		LevelTarget: &mockLevelTarget{},
		ConsumerTarget: &mockConsumerTarget{
			consumer: consumer.NewWriter(os.Stderr),
		},
	}

	givenNames := nlevel.NewNames()
	actual := NewProvider(givenTarget, func(provider *Provider) {
		provider.Level.Names = givenNames
	})

	assert.ToBeNotNil(t, actual)
	assert.ToBeSame(t, givenTarget, actual.Level.Target)
	assert.ToBeSame(t, givenNames, actual.Level.Names)
	assert.ToBeSame(t, givenTarget.GetConsumer(), actual.Consumer.Formatter.Target)
}

type mockProviderTarget struct {
	LevelTarget
	ConsumerTarget
}
