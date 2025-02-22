package problem

import (
	"fmt"
	"strings"
)

type SkillsBuilder struct {
	skills map[string]struct{}
}

// NewBuilder returns a new instance of the Skills builder.
func NewSkillsBuilder() *SkillsBuilder {
	return &SkillsBuilder{
		skills: make(map[string]struct{}),
	}
}

// AddSkill adds a skill to the set, transforming it to lowercase.
func (b *SkillsBuilder) AddSkill(skill string) *SkillsBuilder {
	if skill = strings.TrimSpace(strings.ToLower(skill)); skill != "" {
		b.skills[skill] = struct{}{}
	}
	return b
}

// AddAllSkills adds multiple skills to the set.
func (b *SkillsBuilder) AddAllSkills(skills []string) *SkillsBuilder {
	for _, skill := range skills {
		b.AddSkill(skill)
	}
	return b
}

// Build constructs a Skills instance.
func (b *SkillsBuilder) Build() *Skills {
	return NewSkillsFromBuilder(b)
}

type Skills struct {
	skills map[string]struct{}
}

func NewSkillsFromBuilder(builder *SkillsBuilder) *Skills {
	return &Skills{
		skills: builder.skills,
	}
}

// NewSkills creates a new Skills container
func NewSkills() *Skills {
	return &Skills{skills: make(map[string]struct{})}
}

// Values returns all skills in an unmodifiable slice
func (s *Skills) Values() []string {
	var skillList []string
	for skill := range s.skills {
		skillList = append(skillList, skill)
	}
	return skillList
}

// ContainsSkill checks if a skill is in the container (case-insensitive)
func (s *Skills) Contains(skill string) bool {
	_, exists := s.skills[strings.ToLower(strings.TrimSpace(skill))]
	return exists
}

// String returns a string representation of the skills
func (s *Skills) String() string {
	return fmt.Sprintf("[%s]", strings.Join(s.Values(), ", "))
}

// Equals checks if two Skills containers are equal
func (s *Skills) Equals(other *Skills) bool {
	if len(s.skills) != len(other.skills) {
		return false
	}
	for skill := range s.skills {
		if !other.Contains(skill) {
			return false
		}
	}
	return true
}
