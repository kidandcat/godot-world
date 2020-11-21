extends Control

onready var map: GridMap = get_node("../GridMap")
onready var PreviewerClass = preload("res://Previewer.tscn")
onready var rotUp = $Options/Up
onready var rotDown = $Options/Down
onready var rotLeft = $Options/Left
onready var rotRight = $Options/Right
onready var verticalLevel = $Options/VLLabel
onready var model = $Options/Model
onready var previewPos = $Options/Previewer/MeshPreview/MeshPreviewPosition
onready var meshesPreviewList = $Center/Panel/ScrollContainer/UIMeshesPreviewList

# Called when the node enters the scene tree for the first time.
func _ready():
	addMeshLibraryPreviews()
	verticalLevel.text = "Vertical level " + String(map.verticalLevel)
	updateMesh()
	match map.meshRotation:
		"up":
			rotUp.pressed = true
		"down":
			rotDown.pressed = true
		"left":
			rotLeft.pressed = true
		"right":
			rotRight.pressed = true

func update_rotation(rot: String):
	map.meshRotation = rot

func _on_Up_pressed():
	update_rotation("up")
	rotDown.pressed = false
	rotLeft.pressed = false
	rotRight.pressed = false

func _on_Down_pressed():
	update_rotation("down")
	rotUp.pressed = false
	rotLeft.pressed = false
	rotRight.pressed = false

func _on_Left_pressed():
	update_rotation("left")
	rotUp.pressed = false
	rotDown.pressed = false
	rotRight.pressed = false

func _on_Right_pressed():
	update_rotation("right")
	rotUp.pressed = false
	rotDown.pressed = false
	rotLeft.pressed = false

func _on_VerticalLevel_text_changed(new_text):
	map.verticalLevel = int(new_text)

func _on_VLUp_pressed():
	map.verticalLevel += 1
	verticalLevel.text = "Vertical level " + String(map.verticalLevel)

func _on_VLDown_pressed():
	map.verticalLevel -= 1
	verticalLevel.text = "Vertical level " + String(map.verticalLevel)

func updateMesh():
	var m: Mesh = map.mesh_library.get_item_mesh(map.meshType)
	var meshI = MeshInstance.new()
	meshI.mesh = m
	#meshI.layers = 2
	previewPos.remove_child(previewPos.get_child(0))
	previewPos.add_child(meshI)

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
		map.meshType = meshId
		updateMesh()
