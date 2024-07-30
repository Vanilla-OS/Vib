package recipe

import "list"

_noCommonRecipeNameId: true & id != name

for i, stg in stages {
	_uniqueStageIds: {"\(stg.id)": i}
	_uniqueRecipeId: stg.id != id
	if stg.copy != _|_ {
		_stageIds: [for stg in stages {stg.id}]
		for cp in stg.copy {
			if cp.from != _|_ {
				_noFromOnInitialStage: true & i != 0
				_noFromSameStage:      true & stg.id != cp.from
				_noFromNextStage:      true & _uniqueStageIds[cp.from] < i
				_validFrom:            true & list.Contains(_stageIds, cp.from)
			}
		}
	}

	if stg.modules != _|_ {
		for i2, mod in stg.modules {
			_uniqueRecipeModuleNames: {"\(mod.name)": i}
			_uniqueStageModuleNames: {"\(mod.name)": i2}
			_uniqueRecipeName: true & mod.name != name
		}
	}

	if stg.args != _|_ {
		_noEmptyArgs: true & len([for arg in stg.args {arg}]) > 0
	}

	if stg.expose != _|_ {
		_noEmptyExpose: true & len([for expose in stg.expose {expose}]) > 0
	}

	if stg.env != _|_ {
		_noEmptyEnv: true & len([for env in stg.env {env}]) > 0
	}

	if stg.labels != _|_ {
		_noEmptyLabels: true & len([for labels in stg.labels {labels}]) > 0
	}
}

id:   #string
name: #string
stages: [...#Stage] & list.MinItems(1)

#string: string & !="" & !=null

#Add: close({
	srcdst!: [#string]: #string
	workdir?: #string
	_atLeastOneAddPath: true & list.MinItems([for k in srcdst {k}], 1)
})

#Copy: close({
	srcdst!: [#string]: #string
	from?:    #string
	workdir?: #string
	_atLeastOneCopyPath: true & list.MinItems([for k in srcdst {k}], 1)
})

#Run: close({
	commands!: [...#string] & list.MinItems(1)
	workdir?: #string
})

#Cmd: close({
	exec!: [...#string] & list.MinItems(1)
	workdir?: #string
})

#Entrypoint: close({
	exec!: [...#string] & list.MinItems(1)
	workdir?: #string
})

#ModuleTypes:
	"apt" |
	"dnf" |
	"cmake" |
	"dpkg-buildpackage" |
	"dpkg" |
	"go" |
	"make" |
	"meson" |
	"shell" |
	"includes"

#Source: {
	type!: "tar" | "file" | "git"
	url!:  #string
	if type == "tar" {
		checksum?: #string
	}
	if type == "file" {
		checksum?: #string
	}
	if type == "git" {
		{
			branch!: #string
			commit!: #string
		} |
		{
			tag!: #string
		}
	}
}

#InstFile: #string & =~".+.inst"
#DebFile:  #string & =~".+.deb"
#ModFile:  #string & =~"^modules\/.+.yml"

#AptModuleOpts: close({
	"noRecommends"?:       bool
	"installSuggestions"?: bool
	"fixMissing"?:         bool
	"fixBroken"?:          bool
})

#AptModule: close({
	source!:
		{
			"packages"!: [...#string] & list.MinItems(1)
		} |
		{
			"paths"!: [...#InstFile] & list.MinItems(1)
		}
	options?: #AptModuleOpts

	if options != _|_ {
		_optCheck: true & len(options) > 0
	}

})

#DnfModule: close({
	source!:
		{
			"packages"!: [...#string] & list.MinItems(1)
		} |
		{
			"paths"!: [...#InstFile] & list.MinItems(1)
		}
})

#CmakeModule: close({
	source!:     #Source
	buildflags?: #string
})

#DpkgBuildPackageModule: close({
	source!: #Source
})

#DpkgModule: close({
	source!:
	{
		"paths"!: [...#DebFile] & list.MinItems(1)
	}
})

#GoModule: close({
	source!:     #Source
	buildflags?: #string
})

#MakeModule: close({
	source!:     #Source
	buildflags?: #string
})

#MesonModule: close({
	source!:     #Source
	buildflags?: #string
})

#ShellModule: close({
	commands!: [...#string] & list.MinItems(1)
	workdir?: #string
})

#IncludesModule: close({
	includes!: [... #ModFile] & list.MinItems(1)
})

#Module: close({
	name!: #string
	type!: #ModuleTypes
	if type != _|_ {
		if type == "apt" {#AptModule}
		if type == "dnf" {#DnfModule}
		if type == "cmake" {#CmakeModule}
		if type == "dpkg-buildpackage" {#DpkgBuildPackageModule}
		if type == "dpkg" {#DpkgModule}
		if type == "go" {#GoModule}
		if type == "make" {#MakeModule}
		if type == "meson" {#MesonModule}
		if type == "shell" {#ShellModule}
		if type == "includes" {#IncludesModule}
	}
})

#Stage: close({
	id!:   #string
	base!: #string
	labels?: [#string]: #string
	singlelayer?: bool
	adds?: [...#Add] & list.MinItems(1)
	copy?: [...#Copy] & list.MinItems(1)
	args?: [#string]:   #string
	env?: [#string]:    #string
	expose?: [#string]: #string
	runs?:       #Run
	cmd?:        #Cmd
	entrypoint?: #Entrypoint
	modules?: [...#Module] & list.MinItems(1)
})

#Recipe: close({
	id!:     id
	name!:   name
	stages!: stages
})

#Recipe

// INPUT RECIPE BELOW

id:   "my-image-id"
name: "my-image"
stages: [
	{
		id:   "build"
		base: "node:lts"
		entrypoint: {
			workdir: "/app"
			exec: ["npm", "run", "build"]
		}
	},
]
