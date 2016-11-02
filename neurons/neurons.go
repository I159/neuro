package perceptron

import (
	"encoding/binary"
	"errors"
	//	"fmt"
	"math"
	"math/rand"
	"os"
)

const SCALING_BASE = 0.7
const WEIGHTS_LIM = 0.5

/*TODO: Include Perceive method to a common Neuroner interface. All the neurons
* should be able to perceive .*/
type Neuroner interface {
	GenRandWeights()
	NguyenWiderow()
}

type Neuron struct {
	Weights []float64
}

func (neuron *Neuron) NguyenWidrow(picture_size float64, items_num float64) {
	scaling_factor := SCALING_BASE * math.Pow(picture_size, 1.0/items_num)
	quad_sum := float64(0)
	for _, j := range neuron.Weights {
		quad_sum += math.Pow(j, 2)
	}
	quad_sum = math.Sqrt(quad_sum)
	for i, v := range neuron.Weights {
		neuron.Weights[i] = (scaling_factor * v) / quad_sum
	}
}

func (neuron *Neuron) GenRandWeights(picture_size float64, items_num float64) {
	weights := []float64{}
	r := rand.New(rand.NewSource(99))
	for i := 0.0; i < picture_size; i++ {
		weights = append(weights, math.Mod(r.NormFloat64(), WEIGHTS_LIM))
	}
	neuron.Weights = weights
	neuron.NguyenWidrow(picture_size, items_num)
}

type InputNeuron struct {
	*Neuron
	ImageVector []float64
}

/* TODO: read about panic. */
func check(e error) {
	if e != nil {
		panic(e)
	}
}

/* TODO: Create a ImagesFile type and bind all the file specific methods to the type .*/
func IsValidBinFile(offset int, file *os.File) (int, bool) {
	magic_number_bin := make([]byte, 4)
	new_offset, err := file.ReadAt(magic_number_bin, int64(offset))
	check(err)
	magic_number := binary.BigEndian.Uint32(magic_number_bin)
	offset += new_offset
	if magic_number != 2051 {
		return offset, false
	}
	return offset, true
}

func GetDataLength(offset int, file *os.File) (int, uint32) {
	bin_data_type := make([]byte, 4)
	new_offset, err := file.ReadAt(bin_data_type, int64(offset))
	check(err)
	return offset + new_offset, binary.BigEndian.Uint32(bin_data_type)
}

func GetImageSize(offset int, file *os.File) (int, uint32, uint32) {
	x_dim := make([]byte, 4)
	y_dim := make([]byte, 4)

	new_offset, err := file.ReadAt(x_dim, int64(offset))
	check(err)
	offset += new_offset
	x_dim_size := binary.BigEndian.Uint32(x_dim)

	new_offset, err = file.ReadAt(y_dim, int64(offset))
	check(err)
	offset += new_offset
	y_dim_size := binary.BigEndian.Uint32(y_dim)

	return offset, x_dim_size, y_dim_size
}

func GetImage(offset int, step uint32, file *os.File) (int, []byte) {
	image := make([]byte, step)
	new_offset, err := file.ReadAt(image, int64(offset))
	check(err)
	return offset + new_offset, image
}

func (i *InputNeuron) Perceive(file_path string) (error, *[][]byte) {
	var invalid_file_error error
	images := new([][]byte)

	f, err := os.Open(file_path)
	check(err)
	offset := 0

	offset, is_valid := IsValidBinFile(offset, f)

	if is_valid == false {
		invalid_file_error = errors.New("Invalid mnist file.")
		return invalid_file_error, images
	}

	offset, data_length := GetDataLength(offset, f)
	offset, x_size, y_size := GetImageSize(offset, f)
	new_offset := offset
	for i := 0; i < int(data_length); i++ {
		offset, image := GetImage(new_offset, x_size*y_size, f)
		*images = append(*images, image)
		new_offset = offset
	}
	return invalid_file_error, images
}

func NewNeuron(picture_size, items_num float64) *Neuron {
	neuron := new(Neuron)
	neuron.GenRandWeights(picture_size, items_num)
	return neuron
}

func NewInputNeuron(picture_size, items_num float64) *InputNeuron {
	neuron := InputNeuron{NewNeuron(picture_size, items_num), *new([]float64)}
	return &neuron
}
