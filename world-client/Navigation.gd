extends Navigation

onready var robot  = get_node("../Player")
onready var camera = get_node("../CameraContainer/Camera")
onready var ray = camera.get_node("Ray")

var hit
var prevPos
var prevType

func _ready():
	set_process_input(true)

func _physics_process(delta):
	var mouse = get_viewport().get_mouse_position()
	var ray_origin = camera.project_ray_origin(mouse)
	var ray_direction = camera.project_ray_normal(mouse)
	var from = ray_origin
	var to = ray_origin + ray_direction * 1000.0
	var space_state = get_world().get_direct_space_state()
	hit = space_state.intersect_ray(from, to)
	if hit.size() != 0:
		if prevPos:
			$NavigationMeshInstance/GridMap.set_cell_item(prevPos.x, 0, prevPos.z, prevType, 0)
		var target = hit.position
		var coords = $NavigationMeshInstance/GridMap.world_to_map(target)
		prevType = $NavigationMeshInstance/GridMap.get_cell_item(coords.x, 0, coords.z)
		prevPos = coords
		$NavigationMeshInstance/GridMap.set_cell_item(coords.x, 0, coords.z, 4, 0)

func _unhandled_input(event):
	if event is InputEventMouseButton and event.button_index == BUTTON_LEFT and event.pressed:
		if hit.size() != 0:
			var target = hit.position
			var coords = $NavigationMeshInstance/GridMap.world_to_map(target)
			$NavigationMeshInstance/GridMap.set_cell_item(coords.x, 0, coords.z, 4, 0)
			Networking.newMesh(4, coords.x, coords.z)
			
	if event is InputEventMouseButton and event.button_index == BUTTON_RIGHT and event.pressed:
		if hit.size() != 0:
			var target = hit.position
			var coords = $NavigationMeshInstance/GridMap.world_to_map(target)
			$NavigationMeshInstance/GridMap.set_cell_item(coords.x, 0, coords.z, 0, 0)
			Networking.newMesh(0, coords.x, coords.z)
			

