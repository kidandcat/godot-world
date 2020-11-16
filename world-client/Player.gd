extends MeshInstance


export var gravity = Vector3.DOWN * 30
export var speed = 4
export var jump_speed = 8

var velocity = Vector3.ZERO
var running = false
var rotation_lerp = 0.0
var rotation_speed = 1.0
var state_machine

func _ready():
	add_to_group("Player")
	#state_machine = $AnimationTree.get("parameters/playback")
