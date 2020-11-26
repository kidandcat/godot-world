extends GridMap

onready var player  = get_node("../Player")
onready var camera = get_node("../CameraContainer/Camera")
onready var ray = camera.get_node("Ray")
onready var probe: GIProbe = get_node("../GIProbe")
onready var UI = get_node("../UI")

export var verticalLevel: int = 0
export var heightLevelThreshold = 20 # how many vertical levels is a floor up

var shouldBakeLight = true
var hit
var blinkOff = true
var selectedPos: Vector3
var selectedType: int = -1
var selectedRot: int
var shouldClearLastMesh = false

func _ready():
	set_process_input(true)
	UI.editMode(false)

func _unhandled_input(event):
	if event is InputEventMouseButton and event.pressed:
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
			print("event", event, event.button_index)
		
			if event.button_index == BUTTON_LEFT:
				if hit.size() != 0:
					Networking.sendMove(coords.x, coords.y, coords.z)
					
			if event.button_index == BUTTON_RIGHT:
				if hit.size() != 0:
					UI.editMode(true)
					var t: int
					for i in range(coords.y, coords.y - heightLevelThreshold, -1):
						t = get_cell_item(coords.x, i, coords.z)
						if t > -1:
							coords = Vector3(coords.x, i, coords.z)
							break
					if t > -1:
						clear_blink()
						selectedPos = coords
						selectedType = get_cell_item(coords.x, coords.y, coords.z)
						selectedRot = get_cell_item_orientation(coords.x, coords.y, coords.z)
						UI.on_mesh_selected()

func updateSelectedValues():
	selectedType = get_cell_item(selectedPos.x, selectedPos.y, selectedPos.z)
	selectedRot = get_cell_item_orientation(selectedPos.x, selectedPos.y, selectedPos.z)
	if selectedType < 0:
		selectedType = 0
		shouldClearLastMesh = true
	else:
		shouldClearLastMesh = false

func stopEditing():
	UI.editMode(false)
	clear_blink()
	selectedType = -1

func startEditing():
	UI.editMode(true)
	var coords = world_to_map(player.global_transform.origin)
	var t: int
	for i in range(coords.y, coords.y - heightLevelThreshold, -1):
		t = get_cell_item(coords.x, i, coords.z)
		if t > -1:
			coords = Vector3(coords.x, i, coords.z)
			break
	if t > -1:
		clear_blink()
		selectedPos = coords
		selectedType = get_cell_item(coords.x, coords.y, coords.z)
		selectedRot = get_cell_item_orientation(coords.x, coords.y, coords.z)
		UI.on_mesh_selected()
	if selectedType == -1:
		selectedType = 0

func clear_blink():
	var t = -1
	if !shouldClearLastMesh:
		t = selectedType
		blinkOff = false
	else:
		blinkOff = true
	if selectedType > -1:
		set_cell_item(selectedPos.x, selectedPos.y, selectedPos.z, t, selectedRot)

func saveMeshSelected():
	print("newMesh", selectedType, selectedPos.x, selectedPos.z, Networking.rotation_to_string(selectedRot), selectedPos.y)
	Networking.newMesh(selectedType, selectedPos.x, selectedPos.z, Networking.rotation_to_string(selectedRot), selectedPos.y)
	stopEditing()
	
func deleteMeshSelected():
	Networking.deleteMesh(selectedPos.x, selectedPos.z, selectedPos.y)
	stopEditing()

func createMesh(x: int, y: int, z: int, mesh: int, rotation: int):
	set_cell_item(x, y, z, mesh, rotation)
	shouldBakeLight = true
	
func deleteMesh(x: int, y: int, z: int):
	set_cell_item(x, y, z, -1)
	shouldBakeLight = true

func _on_LightBake_timeout():
	if shouldBakeLight:
		probe.bake()
		shouldBakeLight = false

func _on_Blink_timeout():
	if selectedType > -1:
		if blinkOff:
			set_cell_item(selectedPos.x, selectedPos.y, selectedPos.z, selectedType, selectedRot)
			blinkOff = false
		else:
			set_cell_item(selectedPos.x, selectedPos.y, selectedPos.z, -1) # clear cell
			blinkOff = true
