package syncx

type Group struct {
	tasks []task
}

func (g *Group) Add(execute func() error, interrupt func(error)) {
	g.tasks = append(g.tasks, task{execute, interrupt})
}

func (g *Group) Run() error {
	if len(g.tasks) == 0 {
		return nil
	}

	errors := make(chan error, len(g.tasks))
	for _, a := range g.tasks {
		go func(a task) {
			errors <- a.execute()
		}(a)
	}

	err := <-errors

	for _, a := range g.tasks {
		a.interrupt(err)
	}

	for i := 1; i < cap(errors); i++ {
		<-errors
	}

	return err
}

type task struct {
	execute   func() error
	interrupt func(error)
}
