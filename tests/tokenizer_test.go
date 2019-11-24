package tokenizer_test

import (
  "testing"

  "github.com/stretchr/testify/assert"

  "github.com/abdulbahajaj/brutus/pkg/tokenizer"
)

func TestTesting(t *testing.T) {
  assert.True(t, tokenizer.Testing())
}
