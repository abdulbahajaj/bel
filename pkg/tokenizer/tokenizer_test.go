package tokenizer_test

import (
  "testing"

  "github.com/stretchr/testify/assert"

  "github.com/abdulbahajaj/bel/pkg/tokenizer"
)

func TestTesting(t *testing.T) {
  if !tokenizer.Testing() {
    t.Errorf("Testing failed")
  }
}
