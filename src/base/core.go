package base

import (
	"fmt"
	"log"
	"runtime"
	"../fileUtil"
	"strings"

	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/gl/v4.1-core/gl"
)
func init() {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()
}


type CoreEngine struct {
	Window *glfw.Window
	Program uint32
	IsRunning bool
}

func (e CoreEngine) isRunning()(bool){
	return e.IsRunning
}

func CreateCoreEngine(config fileUtil.Configuration) (e CoreEngine){
	initGLFW()
	e.Window = CreateWindow(config.WINDOW_WIDTH, config.WINDOW_HEIGHT, config.NAME)
	e.Window.MakeContextCurrent()
	
	e.Window.SetKeyCallback(onKey)
	initGl()
	e.intiProgram(config)
	e.IsRunning = false;
	return
}

func (e CoreEngine)Start(){

	if e.isRunning() {
		return
	}

	e.run()
}


func (e CoreEngine)Stop(){
	if !e.isRunning() {
		return
	}

	e.IsRunning = false
}


func (e CoreEngine)run(){
	e.IsRunning = true
	for e.isRunning() == true {
		if e.Window.ShouldClose(){
			e.Stop()
			break;
		}
		
		e.render()
	}
	e.cleanUp()
}

func (e CoreEngine)render(){

	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	gl.UseProgram(e.Program)

	window.SwapBuffers()
	glfw.PollEvents()
}

func (e CoreEngine)cleanUp(){
	defer glfw.Terminate()
}

func newProgram(vertexShaderSource, fragmentShaderSource string) (uint32, error) {
	vertexShader, err := CompileShader(vertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		return 0, err
	}

	fragmentShader, err := CompileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
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

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(program, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to link program: %v", log)
	}

	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)

	return program, nil
}

func (e CoreEngine)intiProgram(config fileUtil.Configuration){
	var err error;
	e.Program, err= newProgram(fileUtil.LoadShader(config.SHADER.VERTEX), fileUtil.LoadShader(config.SHADER.FRAGMENT))
	if err != nil {
		panic(err)
	}
}

func initGl() {
	fmt.Println("Init GL")
	if err := gl.Init(); err != nil {
		panic(err)
	}

	version := gl.GoStr(gl.GetString(gl.VERSION))
	fmt.Println("OpenGL version", version)
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)
	gl.ClearColor(1.0, 1.0, 1.0, 1.0)
}
func initGLFW(){
	fmt.Println("Init GLFW")
	if err := glfw.Init(); err != nil {
		log.Fatalln("failed to initialize glfw:", err)
	}
}