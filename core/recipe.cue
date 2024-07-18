package recipe

import "list"

_uniqueStageIds: {for i, stg in stages {
	"\(stg.id)": i
}}

_duplicateIds: [for stg in stages if stg.id == id {stg.id}]
_uniqueRecipeId: true & len(_duplicateIds) == 0

for stg in stages {
	if stg.copy != _|_ {
		_ids: [for stg in stages {stg.id}]
		for cp in stg.copy {
			if cp.from != _|_ {
				_validFrom: true & list.Contains(_ids, cp.from)
			}
		}
	}
}

id: #string
stages: [...#Stage] & list.MinItems(1)

#string: string & !="" & !=null

#Add: close({
	srcdst!: [#string]: #string
	workdir?: #string
})

#Copy: close({
	srcdst!: [#string]: #string
	from?:    #string
	workdir?: #string
})

#Run: close({
	commands!: [#string]
	workdir?: #string
})

#Cmd: close({
	exec!: [#string]
	workdir?: #string
})

#Entrypoint: close({
	exec!: [#string]
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
	"shell"

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
			"paths"!: [...#InstFile]
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
			"paths"!: [...#InstFile]
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
		"paths"!: [...#DebFile]
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
	commands!: [... #string]
	workdir?: #string
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
	name!:   #string
	stages!: stages
})
