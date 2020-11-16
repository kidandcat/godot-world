extends GridMap

onready var player  = get_node("../Player")
onready var camera = get_node("../CameraContainer/Camera")
onready var ray = camera.get_node("Ray")

export var selectedMeshType: int = 0

var hit
var prevPos
var prevType

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
		if coords == prevPos:
			return
		if prevPos:
			set_cell_item(prevPos.x, 0, prevPos.z, prevType, 0)
		prevType = get_cell_item(coords.x, 0, coords.z)
		prevPos = coords
		set_cell_item(coords.x, 0, coords.z, selectedMeshType, 0)

func _unhandled_input(event):
	if event is InputEventMouseButton and event.button_index == BUTTON_LEFT and event.pressed:
		if hit.size() != 0:
			var target = hit.position
			var coords = world_to_map(target)
			var pos = map_to_world(coords.x, coords.y, coords.z)
			player.moveTo(pos)
			
	if event is InputEventMouseButton and event.button_index == BUTTON_RIGHT and event.pressed:
		if hit.size() != 0:
			var target = hit.position
			var coords = world_to_map(target)
			set_cell_item(coords.x, 0, coords.z, selectedMeshType, 0)
			prevType = selectedMeshType
			Networking.newMesh(selectedMeshType, coords.x, coords.z)
			

