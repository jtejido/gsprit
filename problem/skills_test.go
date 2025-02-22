package problem

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWhenSkillsAdded_TheyShouldBeInSkillSet(t *testing.T) {
	skills := NewSkillsBuilder().AddSkill("skill1").AddSkill("skill2").Build()
	assert.True(t, skills.Contains("skill1"))
	assert.True(t, skills.Contains("skill2"))
}

func TestWhenSkillsAddedCaseInsensitive_TheyShouldBeInSkillSet(t *testing.T) {
	skills := NewSkillsBuilder().AddSkill("skill1").AddSkill("skill2").Build()
	assert.True(t, skills.Contains("skilL1"))
	assert.True(t, skills.Contains("skIll2"))
}

func TestWhenSkillsAddedCaseInsensitive2_TheyShouldBeInSkillSet(t *testing.T) {
	skills := NewSkillsBuilder().AddSkill("Skill1").AddSkill("skill2").Build()
	assert.True(t, skills.Contains("skilL1"))
	assert.True(t, skills.Contains("skIll2"))
}

func TestWhenSkillsAddedThroughAddAll_TheyShouldBeInSkillSet(t *testing.T) {
	skillSet := []string{"skill1", "skill2"}
	skills := NewSkillsBuilder().AddAllSkills(skillSet).Build()
	assert.True(t, skills.Contains("skill1"))
	assert.True(t, skills.Contains("skill2"))
}

func TestWhenSkillsAddedThroughAddAllCaseInsensitive_TheyShouldBeInSkillSet(t *testing.T) {
	skillSet := []string{"skill1", "skill2"}
	skills := NewSkillsBuilder().AddAllSkills(skillSet).Build()
	assert.True(t, skills.Contains("skilL1"))
	assert.True(t, skills.Contains("skill2"))
}

func TestWhenSkillsAddedThroughAddAllCaseInsensitive2_TheyShouldBeInSkillSet(t *testing.T) {
	skillSet := []string{"skill1", "Skill2"}
	skills := NewSkillsBuilder().AddAllSkills(skillSet).Build()
	assert.True(t, skills.Contains("skill1"))
	assert.True(t, skills.Contains("skill2"))
}

func TestWhenSkillsAddedPrecedingWhitespaceShouldNotMatter(t *testing.T) {
	skillSet := []string{" skill1", "Skill2"}
	skills := NewSkillsBuilder().AddAllSkills(skillSet).Build()
	assert.True(t, skills.Contains("skill1"))
	assert.True(t, skills.Contains("skill2"))
}

func TestWhenSkillsAddedTrailingWhitespaceShouldNotMatter(t *testing.T) {
	skillSet := []string{"skill1 ", "Skill2"}
	skills := NewSkillsBuilder().AddAllSkills(skillSet).Build()
	assert.True(t, skills.Contains("skill1"))
	assert.True(t, skills.Contains("skill2"))
}

func TestWhenSkillsAddedTrailingWhitespaceShouldNotMatter2(t *testing.T) {
	skills := NewSkillsBuilder().AddSkill("skill1 ").Build()
	assert.True(t, skills.Contains("skill1"))
}
