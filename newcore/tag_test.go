package newcore

import (
	"testing"
)

func TestNewTagSet(t *testing.T) {
	tag := AddTags.Copy()
	t.Log(tag)
}

func TestTagSetMergeNil(t *testing.T) {
	tag := AddTags.Copy().Merge(nil)
	t.Log(tag)
}
