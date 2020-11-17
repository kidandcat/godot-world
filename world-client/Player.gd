extends MeshInstance

export var speed = 4
var steps = []
var stepsCoords = []
var moving = false

func _ready():
	add_to_group("Player")

func _process(delta):
	if !moving and steps.size() > 0:
		moving = true
		$Tween.interpolate_property(self, "global_transform:origin", global_transform.origin, steps[0], 0.5, Tween.TRANS_LINEAR, Tween.EASE_IN_OUT)
		$Tween.start()
		Networking.notifyMovement(stepsCoords[0][0],stepsCoords[0][1],stepsCoords[0][2])
		steps.remove(0)
		stepsCoords.remove(0)
		
func clear_movement_buffer():
	steps = []
	stepsCoords = []

func moveTo(pos:Vector3, x:int,y:int,z:int):
	steps.append(pos)
	stepsCoords.append([x,y,z])

func _on_Tween_tween_all_completed():
	moving = false
