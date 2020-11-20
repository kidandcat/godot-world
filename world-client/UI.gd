extends CanvasLayer

onready var map: GridMap = get_node("../GridMap")
onready var rotUp = $HBoxContainer/Options/Up
onready var rotDown = $HBoxContainer/Options/Down
onready var rotLeft = $HBoxContainer/Options/Left
onready var rotRight = $HBoxContainer/Options/Right
onready var verticalLevel = $HBoxContainer/Options/VLLabel
onready var model = $HBoxContainer/Options/Model
onready var previewPos = $HBoxContainer/Options/ViewportContainer/MeshPreview/MeshPreviewPosition

# Called when the node enters the scene tree for the first time.
func _ready():
	verticalLevel.text = "Vertical level " + String(map.verticalLevel)
	model.text = String(map.meshType)
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

func _on_Model_text_changed(new_text):
	map.meshType = int(new_text)
	updateMesh()

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
	meshI.layers = 2
	previewPos.remove_child(previewPos.get_child(0))
	previewPos.add_child(meshI)
