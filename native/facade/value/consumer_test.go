package value

import (
	"os"
	"testing"

	"github.com/echocat/slf4g/internal/test/assert"
	"github.com/echocat/slf4g/native/consumer"
)

func Test_NewConsumer_setDefaultConsumer(t *testing.T) {
	givenTarget := &mockConsumerTarget{}

	instance := NewConsumer(givenTarget)

	assert.ToBeNotNil(t, instance)
	assert.ToBeSame(t, consumer.Default, givenTarget.consumer)
	assert.ToBeSame(t, givenTarget.consumer, instance.Formatter.Target)
}

func Test_NewConsumer_setFallbackConsumer(t *testing.T) {
	old := consumer.Default
	defer func() {
		consumer.Default = old
	}()
	consumer.Default = nil

	givenTarget := &mockConsumerTarget{}

	instance := NewConsumer(givenTarget)

	assert.ToBeNotNil(t, instance)
	assert.ToBeEqual(t, consumer.NewWriter(os.Stderr), givenTarget.consumer)
	assert.ToBeSame(t, givenTarget.consumer, instance.Formatter.Target)
}

func Test_NewConsumer_usingExisting(t *testing.T) {
	givenConsumer := consumer.NewWriter(os.Stderr)
	givenTarget := &mockConsumerTarget{
		consumer: givenConsumer,
	}

	instance := NewConsumer(givenTarget)

	assert.ToBeNotNil(t, instance)
	assert.ToBeSame(t, givenConsumer, givenTarget.consumer)
	assert.ToBeSame(t, givenTarget.consumer, instance.Formatter.Target)
}

func Test_NewConsumer_customize(t *testing.T) {
	givenTarget := &mockConsumerTarget{}
	instance := NewConsumer(givenTarget, func(c *Consumer) {
		c.Formatter.Codec = NoopFormatterCodec()
	})

	assert.ToBeNotNil(t, instance)
	assert.ToBeSame(t, NoopFormatterCodec(), instance.Formatter.Codec)
}

func Test_NewConsumer_panicsOnIncompatible(t *testing.T) {
	givenTarget := &mockConsumerTarget{
		consumer: consumer.NewRecorder(),
	}

	assert.Execution(t, func() {
		NewConsumer(givenTarget)
	}).WillPanicWith("^\\*consumer\\.Recorder does not implement formatter\\.MutableAware$")
}

type mockConsumerTarget struct {
	consumer consumer.Consumer
}

func (instance *mockConsumerTarget) GetConsumer() consumer.Consumer {
	return instance.consumer
}

func (instance *mockConsumerTarget) SetConsumer(v consumer.Consumer) {
	instance.consumer = v
}
