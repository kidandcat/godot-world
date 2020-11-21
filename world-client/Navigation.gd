extends GridMap

onready var player  = get_node("../Player")
onready var camera = get_node("../CameraContainer/Camera")
onready var ray = camera.get_node("Ray")

export var meshType: int = 0
export var meshRotation: String = "down"
export var verticalLevel: int = 0

var hit
var prevPos: Vector3
var prevType: int
var prevRot: int
var prevVerticalLevel: int

func _ready():
	set_process_input(true)

func _physics_process(_delta):
	var mouse = get_viewport().get_mouse_position()
	var ray_origin = camera.project_ray_origin(mouse)
	var ray_direction = camera.project_ray_normal(mouse)
	var from = ray_origin
	var to = ray_origin + ray_direction * 1000.0
	var space_state = get_world().get_direct_space_state()
	hit = space_state.intersect_ray(from, to)
	if hit.size() != 0:
		var target = hit.position
		var coords = world_to_map(target)
		coords.y = verticalLevel
		if coords == prevPos:
			return
		if prevPos:
			set_cell_item(prevPos.x, prevVerticalLevel, prevPos.z, prevType, prevRot)
		prevType = get_cell_item(coords.x, verticalLevel, coords.z)
		prevRot = get_cell_item_orientation(coords.x, verticalLevel, coords.z)
		prevVerticalLevel = verticalLevel
		prevPos = coords
		set_cell_item(coords.x, verticalLevel, coords.z, meshType, Networking.rotation_to_int(meshRotation))

func _unhandled_input(event):
	if event is InputEventMouseButton and event.button_index == BUTTON_LEFT and event.pressed:
		if hit.size() != 0:
			var target = hit.position
			var coords = world_to_map(target)
			Networking.sendMove(coords.x, coords.y, coords.z)
			
	if event is InputEventMouseButton and event.button_index == BUTTON_RIGHT and event.pressed:
		if hit.size() != 0:
			var target = hit.position
			var coords = world_to_map(target)
			prevType = meshType
			Networking.newMesh(meshType, coords.x, coords.z, meshRotation, verticalLevel)
