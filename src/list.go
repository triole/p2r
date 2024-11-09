package main

func list(steps tSyncSteps) {
	for _, step := range steps {
		fol1 := parsePath(step.Local)
		listPath(fol1)
		fol2 := parsePath(step.Remote)
		listPath(fol2)
	}
}

func listPath(fol tPath) {
	if fol.IsLocal {
		runCmd([]string{"ls", "-lah", fol.Path})
	} else {
		runCmd([]string{"ssh", fol.Machine, "ls", "-lah", fol.Path})
	}
}
