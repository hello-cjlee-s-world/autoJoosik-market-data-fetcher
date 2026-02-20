package strategyengine

type GateEvaluator interface {
	Name() string
	Evaluate(ctx EvalContext) (bool, string)
}

type FactorEvaluator interface {
	Name() string
	Evaluate(ctx EvalContext) (float64, string)
}

type EvaluatorRegistry struct {
	gates   map[string]GateEvaluator
	factors map[string]FactorEvaluator
}

func NewEvaluatorRegistry() *EvaluatorRegistry {
	return &EvaluatorRegistry{
		gates:   make(map[string]GateEvaluator),
		factors: make(map[string]FactorEvaluator),
	}
}

func (r *EvaluatorRegistry) RegisterGate(g GateEvaluator) {
	r.gates[g.Name()] = g
}

func (r *EvaluatorRegistry) RegisterFactor(f FactorEvaluator) {
	r.factors[f.Name()] = f
}

func (r *EvaluatorRegistry) EvaluateEntry(ctx EvalContext) EntryDecision {
	decision := EntryDecision{Allow: true}
	for _, gateName := range ctx.Config.RequiredGates {
		gate, ok := r.gates[gateName]
		if !ok {
			decision.Allow = false
			decision.Reason = append(decision.Reason, "missing gate:"+gateName)
			return decision
		}
		pass, reason := gate.Evaluate(ctx)
		if !pass {
			decision.Allow = false
			decision.Reason = append(decision.Reason, gateName+":"+reason)
			return decision
		}
	}

	totalWeight := 0.0
	for name, weight := range ctx.Config.Weights {
		factor, ok := r.factors[name]
		if !ok {
			continue
		}
		score, reason := factor.Evaluate(ctx)
		totalWeight += weight
		decision.Score += score * weight
		if reason != "" {
			decision.Reason = append(decision.Reason, name+":"+reason)
		}
	}
	if totalWeight > 0 {
		decision.Score /= totalWeight
	}
	decision.Allow = decision.Allow && decision.Score >= ctx.Config.EntryThreshold
	return decision
}
