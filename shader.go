package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/go-gl/gl/all-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

// Shader is a compiled shader program contains vertex and fragment shaders.
type Shader struct {
	ID             uint32 // the program ID
	VertexSource   string // vertex shader source code
	FragmentSource string // fragment shader source code
}

// NewShader creates a shader program, it reads shader source from shader files.
func NewShader(vertexFile, fragmentFile string) (*Shader, error) {
	vertexSource, err := ioutil.ReadFile(vertexFile)
	if err != nil {
		return nil, err
	}
	fragmentSource, err := ioutil.ReadFile(fragmentFile)
	if err != nil {
		return nil, err
	}
	program, err := newProgram(string(vertexSource), string(fragmentSource))
	if err != nil {
		fmt.Printf("ERR: %v\n", err)
		return nil, err
	}
	return &Shader{
		ID:             program,
		VertexSource:   string(vertexSource),
		FragmentSource: string(fragmentSource),
	}, nil
}

// Use activates the shader
func (s *Shader) Use() {
	gl.UseProgram(s.ID)
}

// SetUniformName sets uniform by name.
// The number of values must be 1, 2, 3 or 4.
// the type of values must be int32, uint32, float32 or float64
func (s *Shader) SetUniformName(name string, v ...interface{}) error {
	if !strings.HasSuffix(name, "\x00") {
		name += "\x00"
	}
	location := gl.GetUniformLocation(s.ID, gl.Str(name))
	return s.SetUniform(location, v...)
}

// SetUniform sets uniform by location.
// The number of values must be 1, 2, 3 or 4.
// the type of values must be int, int32, uint32, float32 or float64
func (s *Shader) SetUniform(uniform int32, v ...interface{}) error {
	if uniform < 0 {
		return errors.New("invalid uniform")
	}
	if len(v) == 0 {
		return errors.New("empty value")
	}

	switch (v[0]).(type) {
	case int:
		vv := make([]int32, len(v))
		for i := range v {
			tv, ok := v[i].(int)
			if !ok {
				return fmt.Errorf("invalid type of v[%d], must be int", i)
			}
			vv[i] = int32(tv)
		}
		s.SetUniformi(uniform, vv...)
	case int32:
		vv := make([]int32, len(v))
		for i := range v {
			tv, ok := v[i].(int32)
			if !ok {
				return fmt.Errorf("invalid type of v[%d], must be int32", i)
			}
			vv[i] = tv
		}
		s.SetUniformi(uniform, vv...)
	case uint32:
		vv := make([]uint32, len(v))
		for i := range v {
			tv, ok := v[i].(uint32)
			if !ok {
				return fmt.Errorf("invalid type of v[%d], must be uint32", i)
			}
			vv[i] = tv
		}
		s.SetUniformui(uniform, vv...)
	case float32:
		vv := make([]float32, len(v))
		for i := range v {
			tv, ok := v[i].(float32)
			if !ok {
				return fmt.Errorf("invalid type of v[%d], must be float32", i)
			}
			vv[i] = tv
		}
		s.SetUniformf(uniform, vv...)
	case float64:
		vv := make([]float64, len(v))
		for i := range v {
			tv, ok := v[i].(float64)
			if !ok {
				return fmt.Errorf("invalid type of v[%d], must be float64", i)
			}
			vv[i] = tv
		}
		s.SetUniformd(uniform, vv...)
	default:
		return errors.New("unsupported value type")
	}

	return nil
}

// SetUniformi sets uniform with int32 values.
// The number of values must be 1, 2, 3 or 4.
func (s *Shader) SetUniformi(uniform int32, v ...int32) error {
	if uniform < 0 {
		return errors.New("invalid uniform")
	}

	if n := len(v); n == 1 {
		gl.Uniform1i(uniform, v[0])
	} else if n == 2 {
		gl.Uniform2i(uniform, v[0], v[1])
	} else if n == 3 {
		gl.Uniform3i(uniform, v[0], v[1], v[2])
	} else if n >= 4 {
		gl.Uniform4i(uniform, v[0], v[1], v[2], v[3])
	} else {
		return errors.New("empty value")
	}

	return nil
}

// SetUniformui sets uniform with uint32 values.
// The number of values must be 1, 2, 3 or 4.
func (s *Shader) SetUniformui(uniform int32, v ...uint32) error {
	if uniform < 0 {
		return errors.New("invalid uniform")
	}

	if n := len(v); n == 1 {
		gl.Uniform1ui(uniform, v[0])
	} else if n == 2 {
		gl.Uniform2ui(uniform, v[0], v[1])
	} else if n == 3 {
		gl.Uniform3ui(uniform, v[0], v[1], v[2])
	} else if n >= 4 {
		gl.Uniform4ui(uniform, v[0], v[1], v[2], v[3])
	} else {
		return errors.New("empty value")
	}

	return nil
}

// SetUniformf sets uniform with float32 values.
// The number of values must be 1, 2, 3 or 4.
func (s *Shader) SetUniformf(uniform int32, v ...float32) error {
	if uniform < 0 {
		return errors.New("invalid uniform")
	}

	if n := len(v); n == 1 {
		gl.Uniform1f(uniform, v[0])
	} else if n == 2 {
		gl.Uniform2f(uniform, v[0], v[1])
	} else if n == 3 {
		gl.Uniform3f(uniform, v[0], v[1], v[2])
	} else if n >= 4 {
		gl.Uniform4f(uniform, v[0], v[1], v[2], v[3])
	} else {
		return errors.New("empty value")
	}

	return nil
}

// SetUniformd sets uniform with float64 values.
// The number of values must be 1, 2, 3 or 4.
func (s *Shader) SetUniformd(uniform int32, v ...float64) error {
	if uniform < 0 {
		return errors.New("invalid uniform")
	}

	if n := len(v); n == 1 {
		gl.Uniform1d(uniform, v[0])
	} else if n == 2 {
		gl.Uniform2d(uniform, v[0], v[1])
	} else if n == 3 {
		gl.Uniform3d(uniform, v[0], v[1], v[2])
	} else if n >= 4 {
		gl.Uniform4d(uniform, v[0], v[1], v[2], v[3])
	} else {
		return errors.New("empty value")
	}

	return nil
}

func (s *Shader) SetUniformMatrixName(name string, transpose bool, mat interface{}) error {
	if !strings.HasSuffix(name, "\x00") {
		name += "\x00"
	}
	location := gl.GetUniformLocation(s.ID, gl.Str(name))
	return s.SetUniformMatrix(location, transpose, mat)
}

func (s *Shader) SetUniformMatrix(uniform int32, traspose bool, mat interface{}) error {
	switch v := mat.(type) {
	case mgl32.Mat3:
		gl.UniformMatrix2fv(uniform, 1, traspose, &v[0])
	case mgl32.Mat4:
		gl.UniformMatrix4fv(uniform, 1, traspose, &v[0])
	default:
		return errors.New("unsupported matrix")
	}
	return nil
}

func newProgram(vertexShaderSource, fragmentShaderSource string) (uint32, error) {
	vertexShader, err := compileShader(vertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		return 0, err
	}

	fragmentShader, err := compileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		return 0, err
	}

	program := gl.CreateProgram()
	gl.AttachShader(program, vertexShader)
	gl.AttachShader(program, fragmentShader)
	gl.LinkProgram(program)

	var status int32
	gl.GetProgramiv(program, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &logLength)
		logs := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(program, logLength, nil, gl.Str(logs))

		return 0, fmt.Errorf("failed to link program: %v", logs)
	}

	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)

	return program, nil
}

func compileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)
	csources, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)
		logs := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(logs))

		return 0, fmt.Errorf("failed to compile %v : %v", source, logs)
	}

	return shader, nil
}
