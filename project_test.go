package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckAffected(t *testing.T) {
	assert := assert.New(t)

	changedFiles := []string{
		"project1/app.py",
		"project2/README.md",
		"",
		"README.md",
		"project4/source/test.swift",
		"project5/source/project1/main.swift",
		"project6/main.swift",
		"project7/source/main.swift",
		"project7/tests/test.swift",
		"project8/project8Tests/main.swift",
	}

	p1 := Project{Label: "project1", Path: []string{"project1/"}, Skip: []string{}}
	assert.Equal(true, p1.checkAffected(changedFiles))

	p2 := Project{Label: "project2", Path: []string{"project2"}, Skip: []string{"somelabel"}}
	assert.Equal(true, p2.checkAffected(changedFiles))

	p3 := Project{Label: "project3", Path: []string{"project3/", "project2/foo/"}, Skip: []string{"project1"}}
	assert.Equal(false, p3.checkAffected(changedFiles))

	// test no changes
	assert.Equal(false, p3.checkAffected([]string{}))

	p4 := Project{Label: "project4", Path: []string{"project4/"}, ExcludePath: []string{"project4/source"}}
	assert.Equal(false, p4.checkAffected(changedFiles))

	p5 := Project{Label: "project5", Path: []string{"project5/**/*.h"}}
	assert.Equal(false, p5.checkAffected(changedFiles))

	p5 = Project{Label: "project5", Path: []string{"project5/**/*.swift"}}
	assert.Equal(true, p5.checkAffected(changedFiles))

	p5 = Project{Label: "project5", Path: []string{"project5/**/*.swift"}, ExcludePath: []string{"project5/source/project1/main.swift"}}
	assert.Equal(false, p5.checkAffected(changedFiles))

	p6 := Project{Label: "project6", Path: []string{"project6/"}, ExcludePath: []string{"project6/**/*.swift"}}
	assert.Equal(false, p6.checkAffected(changedFiles))

	p6 = Project{Label: "project6", Path: []string{"project6/"}, ExcludePath: []string{"project6/tests/**/*.swift"}}
	assert.Equal(true, p6.checkAffected(changedFiles))

	p7 := Project{Label: "project7", Path: []string{"project7/"}, ExcludePath: []string{"project7/tests"}}
	assert.Equal(true, p7.checkAffected(changedFiles))

	p8 := Project{Label: "project8", Path: []string{"project8/project8"}}
	assert.Equal(false, p8.checkAffected(changedFiles))

	p8 = Project{Label: "project8", Path: []string{"project8/project8/"}}
	assert.Equal(false, p8.checkAffected(changedFiles))
}
