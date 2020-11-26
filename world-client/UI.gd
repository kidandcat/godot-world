extends Control

onready var map: GridMap = get_node("../GridMap")
onready var PreviewerClass = preload("res://Previewer.tscn")
onready var rotUp = $Options/VBoxContainer/HBoxContainer/Up
onready var rotDown = $Options/VBoxContainer/HBoxContainer3/Down
onready var rotLeft = $Options/VBoxContainer/HBoxContainer2/Left
onready var rotRight = $Options/VBoxContainer/HBoxContainer2/Right
onready var verticalLevel = $Options/VLLabel
onready var model = $Options/Model
onready var meshesPreviewList = $Center/Panel/ScrollContainer/UIMeshesPreviewList
onready var editSwitch = $Left/EditSwitch

# Called when the node enters the scene tree for the first time.
func _ready():
	addMeshLibraryPreviews()
	verticalLevel.text = "Vertical level"

func on_mesh_selected():
	verticalLevel.text = "Vertical level " + String(map.selectedPos.y)

func update_rotation(rot: String):
	map.selectedRot = Networking.rotation_to_int(rot)
	print("update_rotation ",rot, " -> ", map.selectedRot)

func _on_Up_pressed():
	map.clear_blink()
	map.selectedPos.z -= 1
	map.updateSelectedValues()

func _on_Down_pressed():
	map.clear_blink()
	map.selectedPos.z += 1
	map.updateSelectedValues()

func _on_Left_pressed():
	map.clear_blink()
	map.selectedPos.x -= 1
	map.updateSelectedValues()

func _on_Right_pressed():
	map.clear_blink()
	map.selectedPos.x += 1
	map.updateSelectedValues()

func _on_VLUp_pressed():
	map.clear_blink()
	map.selectedPos.y += 1
	verticalLevel.text = "Vertical level " + String(map.selectedPos.y)
	map.updateSelectedValues()

func _on_VLDown_pressed():
	map.clear_blink()
	map.selectedPos.y -= 1
	verticalLevel.text = "Vertical level " + String(map.selectedPos.y)
	map.updateSelectedValues()

func addMeshLibraryPreviews():
	var IDsList = map.mesh_library.get_item_list()
	for id in IDsList:
		addMeshPreview(id)

func addMeshPreview(id: int):
	var m: Mesh = map.mesh_library.get_item_mesh(id)
	var meshI = MeshInstance.new()
	meshI.mesh = m
	var previewer = PreviewerClass.instance()
	previewer.get_node("MeshPreview/MeshPreviewPosition").add_child(meshI)
	meshesPreviewList.add_child(previewer)
	previewer.connect("gui_input", self, "_on_mesh_selected", [id])

func _on_mesh_selected(event: InputEvent, meshId: int):
	if event is InputEventMouseButton and event.button_index == BUTTON_LEFT and event.pressed:
		map.clear_blink()
		map.selectedType = meshId

func _on_RotationLeft_pressed():
	match map.selectedRot:
		0: update_rotation("left")
		16: update_rotation("down")
		10: update_rotation("right")
		22: update_rotation("up")

func _on_RotationRight_pressed():
	match map.selectedRot:
		0: update_rotation("right")
		16: update_rotation("up")
		10: update_rotation("left")
		22: update_rotation("down")

func editMode(mode: bool):
	if mode:
		$Options.show()
		$Center.show()
		editSwitch.pressed = true
	else:
		$Options.hide()
		$Center.hide()
		editSwitch.pressed = false

func _on_EditSwitch_toggled(button_pressed):
	if button_pressed:
		map.startEditing()
	else:
		map.stopEditing()

func _on_SaveMesh_pressed():
	map.saveMeshSelected()
	map.stopEditing()

func _on_DeleteMesh_pressed():
	map.deleteMeshSelected()
	map.stopEditing()
