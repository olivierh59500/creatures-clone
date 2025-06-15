package creature

import (
	"math"
	"math/rand"
)

// Brain represents the creature's neural network
type Brain struct {
	// Network structure
	inputSize  int
	hiddenSize []int
	outputSize int

	// Network weights
	weights [][]float64
	biases  [][]float64

	// Activation values
	activations [][]float64

	// Learning parameters
	learningRate float64
	momentum     float64

	// Previous weight changes for momentum
	prevWeightChanges [][]float64
	prevBiasChanges   [][]float64

	// Output values
	output []float64
}

// NewBrain creates a new neural network brain
func NewBrain() *Brain {
	inputSize := 32             // Vision(20) + Internal(7) + Touch(4) + Time(1)
	hiddenSize := []int{20, 20} // Two hidden layers
	outputSize := OutputMax

	b := &Brain{
		inputSize:    inputSize,
		hiddenSize:   hiddenSize,
		outputSize:   outputSize,
		learningRate: 0.1,
		momentum:     0.9,
		output:       make([]float64, outputSize),
	}

	b.initializeNetwork()
	return b
}

// initializeNetwork sets up the network structure
func (b *Brain) initializeNetwork() {
	// Calculate layer sizes
	layerSizes := []int{b.inputSize}
	layerSizes = append(layerSizes, b.hiddenSize...)
	layerSizes = append(layerSizes, b.outputSize)

	// Initialize weights and biases
	b.weights = make([][]float64, len(layerSizes)-1)
	b.biases = make([][]float64, len(layerSizes)-1)
	b.prevWeightChanges = make([][]float64, len(layerSizes)-1)
	b.prevBiasChanges = make([][]float64, len(layerSizes)-1)
	b.activations = make([][]float64, len(layerSizes))

	// Initialize each layer
	for i := 0; i < len(layerSizes); i++ {
		b.activations[i] = make([]float64, layerSizes[i])

		if i < len(layerSizes)-1 {
			// Xavier initialization for weights
			numWeights := layerSizes[i] * layerSizes[i+1]
			b.weights[i] = make([]float64, numWeights)
			b.prevWeightChanges[i] = make([]float64, numWeights)

			// Initialize weights with Xavier initialization
			scale := math.Sqrt(2.0 / float64(layerSizes[i]+layerSizes[i+1]))
			for j := range b.weights[i] {
				b.weights[i][j] = (rand.Float64()*2 - 1) * scale
			}

			// Initialize biases
			b.biases[i] = make([]float64, layerSizes[i+1])
			b.prevBiasChanges[i] = make([]float64, layerSizes[i+1])
			for j := range b.biases[i] {
				b.biases[i][j] = (rand.Float64()*2 - 1) * 0.1
			}
		}
	}
}

// Process runs the neural network forward pass
func (b *Brain) Process(input []float64) {
	// Validate input size
	if len(input) != b.inputSize {
		// Pad or truncate as needed
		if len(input) < b.inputSize {
			padded := make([]float64, b.inputSize)
			copy(padded, input)
			input = padded
		} else {
			input = input[:b.inputSize]
		}
	}

	// Set input layer
	copy(b.activations[0], input)

	// Forward propagation through each layer
	for layer := 0; layer < len(b.weights); layer++ {
		currentLayerSize := len(b.activations[layer])
		nextLayerSize := len(b.activations[layer+1])

		// Calculate activations for next layer
		for j := 0; j < nextLayerSize; j++ {
			sum := b.biases[layer][j]

			// Sum weighted inputs
			for i := 0; i < currentLayerSize; i++ {
				weightIndex := i*nextLayerSize + j
				sum += b.activations[layer][i] * b.weights[layer][weightIndex]
			}

			// Apply activation function (sigmoid)
			b.activations[layer+1][j] = sigmoid(sum)
		}
	}

	// Copy output
	copy(b.output, b.activations[len(b.activations)-1])
}

// Reinforce applies reinforcement learning
func (b *Brain) Reinforce(reward float64) {
	// Simple reinforcement: strengthen connections that led to reward
	// This is a simplified version - a full implementation would use
	// temporal difference learning or policy gradients

	learningFactor := b.learningRate * reward

	// Update weights based on recent activations
	for layer := 0; layer < len(b.weights); layer++ {
		currentLayerSize := len(b.activations[layer])
		nextLayerSize := len(b.activations[layer+1])

		for j := 0; j < nextLayerSize; j++ {
			// Update bias
			biasChange := learningFactor * b.activations[layer+1][j]
			b.biases[layer][j] += biasChange + b.momentum*b.prevBiasChanges[layer][j]
			b.prevBiasChanges[layer][j] = biasChange

			// Update weights
			for i := 0; i < currentLayerSize; i++ {
				weightIndex := i*nextLayerSize + j
				weightChange := learningFactor * b.activations[layer][i] * b.activations[layer+1][j]
				b.weights[layer][weightIndex] += weightChange + b.momentum*b.prevWeightChanges[layer][weightIndex]
				b.prevWeightChanges[layer][weightIndex] = weightChange
			}
		}
	}
}

// Learn performs supervised learning with target outputs
func (b *Brain) Learn(input []float64, target []float64) {
	// Process input first
	b.Process(input)

	// Calculate errors using backpropagation
	layerCount := len(b.activations)
	errors := make([][]float64, layerCount)
	for i := range errors {
		errors[i] = make([]float64, len(b.activations[i]))
	}

	// Calculate output layer errors
	outputLayer := layerCount - 1
	for i := 0; i < b.outputSize; i++ {
		output := b.activations[outputLayer][i]
		errors[outputLayer][i] = (target[i] - output) * sigmoidDerivative(output)
	}

	// Backpropagate errors
	for layer := outputLayer - 1; layer > 0; layer-- {
		currentLayerSize := len(b.activations[layer])
		nextLayerSize := len(b.activations[layer+1])

		for i := 0; i < currentLayerSize; i++ {
			sum := 0.0
			for j := 0; j < nextLayerSize; j++ {
				weightIndex := i*nextLayerSize + j
				sum += errors[layer+1][j] * b.weights[layer][weightIndex]
			}
			errors[layer][i] = sum * sigmoidDerivative(b.activations[layer][i])
		}
	}

	// Update weights and biases
	for layer := 0; layer < len(b.weights); layer++ {
		currentLayerSize := len(b.activations[layer])
		nextLayerSize := len(b.activations[layer+1])

		for j := 0; j < nextLayerSize; j++ {
			// Update bias
			biasChange := b.learningRate * errors[layer+1][j]
			b.biases[layer][j] += biasChange + b.momentum*b.prevBiasChanges[layer][j]
			b.prevBiasChanges[layer][j] = biasChange

			// Update weights
			for i := 0; i < currentLayerSize; i++ {
				weightIndex := i*nextLayerSize + j
				weightChange := b.learningRate * errors[layer+1][j] * b.activations[layer][i]
				b.weights[layer][weightIndex] += weightChange + b.momentum*b.prevWeightChanges[layer][weightIndex]
				b.prevWeightChanges[layer][weightIndex] = weightChange
			}
		}
	}
}

// GetOutput returns the current output values
func (b *Brain) GetOutput() []float64 {
	return b.output
}

// GetWeights returns a copy of all weights for genetic inheritance
func (b *Brain) GetWeights() [][]float64 {
	weightsCopy := make([][]float64, len(b.weights))
	for i := range b.weights {
		weightsCopy[i] = make([]float64, len(b.weights[i]))
		copy(weightsCopy[i], b.weights[i])
	}
	return weightsCopy
}

// SetWeights sets the network weights (used in breeding)
func (b *Brain) SetWeights(weights [][]float64) {
	if len(weights) != len(b.weights) {
		return
	}

	for i := range weights {
		if len(weights[i]) == len(b.weights[i]) {
			copy(b.weights[i], weights[i])
		}
	}
}

// Mutate randomly modifies some weights
func (b *Brain) Mutate(mutationRate float64) {
	for layer := range b.weights {
		// Mutate weights
		for i := range b.weights[layer] {
			if rand.Float64() < mutationRate {
				// Add gaussian noise
				b.weights[layer][i] += (rand.Float64()*2 - 1) * 0.1
			}
		}

		// Mutate biases
		for i := range b.biases[layer] {
			if rand.Float64() < mutationRate {
				b.biases[layer][i] += (rand.Float64()*2 - 1) * 0.1
			}
		}
	}
}

// sigmoid activation function
func sigmoid(x float64) float64 {
	return 1.0 / (1.0 + math.Exp(-x))
}

// sigmoidDerivative calculates the derivative of sigmoid
func sigmoidDerivative(x float64) float64 {
	return x * (1.0 - x)
}

// Save serializes the brain to a byte array
func (b *Brain) Save() []byte {
	// In a full implementation, this would serialize the network
	// For now, return empty
	return []byte{}
}

// Load deserializes the brain from a byte array
func (b *Brain) Load(data []byte) {
	// In a full implementation, this would deserialize the network
}
