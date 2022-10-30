package stepper

type Stepper interface {
	StartStep(step string) Step

	StartTemporaryStep(step string) Step
	StartTemporaryStepWithoutNewLine(step string) Step

	StopCurrentStep()
}
