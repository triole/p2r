package main

func list(steps tSyncSteps) {
	for _, step := range steps {
		listPath(step.Set.Local)
		listPath(step.Set.Remote)
	}
}

func listPath(fol tPath) {
	if fol.IsLocal {
		runCmd([]string{"ls", "-lah", fol.FullPath})
	} else {
		runCmd([]string{"ssh", fol.Machine, "ls", "-lah", fol.Path})
	}
}
