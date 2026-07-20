package llm

import (
	"context"
	"fmt"

	"nextcode.io/fantasy"
)

// TaskType represents different types of code review tasks
type TaskType string

const (
	TaskSecurity     TaskType = "security"
	TaskPerformance  TaskType = "performance"
	TaskStyle        TaskType = "style"
	TaskDocumentation TaskType = "documentation"
	TaskTesting      TaskType = "testing"
	TaskFix          TaskType = "fix"
)

// ModelRouter intelligently routes tasks to different LLM models
type ModelRouter struct {
	defaultModel fantasy.LanguageModel
	models       map[string]fantasy.LanguageModel
	routing      map[TaskType]string // Task -> Model name
	costPerToken map[string]float32  // Model -> cost per 1K tokens
}

// NewModelRouter creates a new model router
func NewModelRouter(defaultModel fantasy.LanguageModel) *ModelRouter {
	return &ModelRouter{
		defaultModel: defaultModel,
		models:       make(map[string]fantasy.LanguageModel),
		routing:      make(map[TaskType]string),
		costPerToken: make(map[string]float32),
	}
}

// RegisterModel registers a new LLM model
func (r *ModelRouter) RegisterModel(name string, model fantasy.LanguageModel, costPer1KTokens float32) {
	r.models[name] = model
	r.costPerToken[name] = costPer1KTokens
}

// SetTaskRouting sets the model to use for a specific task
func (r *ModelRouter) SetTaskRouting(taskType TaskType, modelName string) error {
	if _, exists := r.models[modelName]; !exists {
		return fmt.Errorf("model %s not registered", modelName)
	}
	r.routing[taskType] = modelName
	return nil
}

// GetModelForTask returns the best model for a task
func (r *ModelRouter) GetModelForTask(taskType TaskType) fantasy.LanguageModel {
	if modelName, exists := r.routing[taskType]; exists {
		if model, ok := r.models[modelName]; ok {
			return model
		}
	}

	// Default models if not configured
	switch taskType {
	case TaskSecurity:
		// Security needs accuracy - use best model
		if model, ok := r.models["claude-opus"]; ok {
			return model
		}
	case TaskPerformance:
		// Performance analysis - balanced
		if model, ok := r.models["gpt-4"]; ok {
			return model
		}
	case TaskStyle:
		// Style - can use faster/cheaper model
		if model, ok := r.models["gpt-3.5-turbo"]; ok {
			return model
		}
	case TaskTesting, TaskDocumentation:
		// Generation tasks - quality matters
		if model, ok := r.models["gpt-4o"]; ok {
			return model
		}
	}

	return r.defaultModel
}

// EstimateCost estimates the cost of a task
func (r *ModelRouter) EstimateCost(taskType TaskType, tokensUsed int) float32 {
	model := r.GetModelForTask(taskType)

	// Find which model was used
	for modelName, m := range r.models {
		if m == model {
			costPer1K := r.costPerToken[modelName]
			return costPer1K * float32(tokensUsed) / 1000.0
		}
	}

	return 0.0
}

// ListAvailableModels returns all registered models
func (r *ModelRouter) ListAvailableModels() map[string]float32 {
	return r.costPerToken
}

// GetRoutingConfig returns the current routing configuration
func (r *ModelRouter) GetRoutingConfig() map[TaskType]string {
	config := make(map[TaskType]string)
	for task, model := range r.routing {
		config[task] = model
	}
	return config
}

// OptimizeForCost reconfigures routing to minimize costs
func (r *ModelRouter) OptimizeForCost() {
	// Route less critical tasks to cheaper models
	r.routing[TaskStyle] = "gpt-3.5-turbo"       // Cheapest
	r.routing[TaskPerformance] = "gpt-4"         // Medium
	r.routing[TaskDocumentation] = "gpt-3.5-turbo" // Cheap (good enough)
	r.routing[TaskTesting] = "gpt-4"             // Quality needed
	r.routing[TaskSecurity] = "claude-opus"      // Best accuracy
}

// OptimizeForQuality reconfigures routing to maximize quality
func (r *ModelRouter) OptimizeForQuality() {
	// Route all tasks to best models regardless of cost
	r.routing[TaskSecurity] = "claude-opus"
	r.routing[TaskPerformance] = "gpt-4"
	r.routing[TaskStyle] = "gpt-4"
	r.routing[TaskTesting] = "gpt-4o"
	r.routing[TaskDocumentation] = "gpt-4o"
}

// OptimizeForSpeed reconfigures routing for fastest responses
func (r *ModelRouter) OptimizeForSpeed() {
	// Route to fastest models
	r.routing[TaskStyle] = "gpt-3.5-turbo"
	r.routing[TaskPerformance] = "gpt-3.5-turbo"
	r.routing[TaskSecurity] = "gpt-4"
	r.routing[TaskTesting] = "gpt-4"
	r.routing[TaskDocumentation] = "gpt-3.5-turbo"
}

// CallModel calls the appropriate model for a task
func (r *ModelRouter) CallModel(ctx context.Context, taskType TaskType, messages []fantasy.Message, maxTokens int) (*fantasy.CompleteResponse, error) {
	model := r.GetModelForTask(taskType)

	response, err := model.Complete(ctx, messages, &fantasy.CompleteOptions{
		MaxTokens:   maxTokens,
		Temperature: r.getTemperatureForTask(taskType),
	})

	return response, err
}

// getTemperatureForTask returns appropriate temperature for task type
func (r *ModelRouter) getTemperatureForTask(taskType TaskType) float32 {
	switch taskType {
	case TaskSecurity:
		return 0.1 // Very conservative - accuracy critical
	case TaskPerformance:
		return 0.3 // Conservative - accuracy important
	case TaskStyle:
		return 0.5 // Balanced
	case TaskDocumentation:
		return 0.7 // More creative
	case TaskTesting:
		return 0.4 // Coverage-focused
	case TaskFix:
		return 0.3 // Conservative - correctness critical
	default:
		return 0.5
	}
}

// CostOptimizationStrategy defines how to optimize routing
type CostOptimizationStrategy string

const (
	StrategyCost    CostOptimizationStrategy = "cost"
	StrategyQuality CostOptimizationStrategy = "quality"
	StrategySpeed   CostOptimizationStrategy = "speed"
	StrategyBalanced CostOptimizationStrategy = "balanced"
)

// ApplyStrategy applies a cost optimization strategy
func (r *ModelRouter) ApplyStrategy(strategy CostOptimizationStrategy) {
	switch strategy {
	case StrategyCost:
		r.OptimizeForCost()
	case StrategyQuality:
		r.OptimizeForQuality()
	case StrategySpeed:
		r.OptimizeForSpeed()
	case StrategyBalanced:
		// Default/balanced routing
	}
}
