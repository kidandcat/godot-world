extends MeshInstance

export var speed = 4
var steps = []
var moving = false

func _ready():
	add_to_group("Player")
	#state_machine = $AnimationTree.get("parameters/playback")

func _process(delta):
	if !moving and steps.size() > 0:
		moving = true
		$Tween.interpolate_property(self, "global_transform:origin", global_transform.origin, steps[0], 0.5, Tween.TRANS_LINEAR, Tween.EASE_IN_OUT)
		$Tween.start()
		steps.remove(0)

func moveTo(pos:Vector3):
	steps.append(pos)


func _on_Tween_tween_all_completed():
	moving = false
