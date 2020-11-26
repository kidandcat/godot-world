extends Node

export var websocket_url = "wss://131d454e7693.ngrok.io"
var _client = WebSocketClient.new()
onready var gridMap: GridMap  = get_node("/root/World/GridMap")
onready var player  = get_node("/root/World/Player")

func _ready():
	print("Networking Ready")
	_client.connect("connection_closed", self, "_closed")
	_client.connect("connection_error", self, "_closed")
	_client.connect("connection_established", self, "_connected")
	_client.connect("data_received", self, "_on_data")
	
	var err = _client.connect_to_url(websocket_url)
	if err != OK:
		print("Unable to connect")
		set_process(false)

func _closed(was_clean = false):
	print("Closed, clean: ", was_clean)
	set_process(false)

func _connected(proto = ""):
	print("Connected with protocol: ", proto)
	var err = _client.get_peer(1).put_packet("login:jairo,get".to_utf8())
	print("packet sent ", err)
	_client.get_peer(1).put_packet("world_around:0,0".to_utf8())

func _on_data():
	var msg = _client.get_peer(1).get_packet().get_string_from_utf8().split(":")
	print("<- ", msg)
	match msg[0]:
		"player":
			var data = msg[1].split(",")
			var pos = gridMap.map_to_world(int(data[0]), 0, int(data[1])) # x,y,z
			player.global_transform.origin = pos # teleport player to current position
		"newmesh":
			var data = msg[1].split(",")
			var rotation: int = rotation_to_int(data[4])
			gridMap.createMesh(int(data[1]), int(data[3]), int(data[2]), int(data[0]), rotation)
			if gridMap.selectedPos == Vector3(int(data[1]), int(data[3]), int(data[2])):
				# stop editing if the new mesh created is the one we are editing
				gridMap.selectedType = -1 # unselect editing mesh
		"deletemesh":
			var data = msg[1].split(",")
			gridMap.deleteMesh(int(data[0]), int(data[2]), int(data[1]))
		"move":
			var data = msg[1].split(",")
			var pos = gridMap.map_to_world(int(data[0]), 0, int(data[1])) # x,y,z
			player.moveTo(pos, int(data[0]), 0, int(data[1]))

func _process(delta):
	_client.poll()

func newMesh(type:int, x:float, z:float, rotation:String, verticalLevel:int):
	var msg = "create_mesh:"+String(type)+","+String(x)+","+String(z)+","+String(verticalLevel)+","+rotation
	_client.get_peer(1).put_packet(msg.to_utf8())
	print("-> ", msg)

func deleteMesh(x:float, z:float, verticalLevel:int):
	var msg = "delete_mesh:"+String(x)+","+String(z)+","+String(verticalLevel)
	_client.get_peer(1).put_packet(msg.to_utf8())
	print("-> ", msg)

func sendMove(x:int, y:int, z:int):
	var msg = "walk_to:"+String(x)+","+String(z)
	player.clear_movement_buffer()
	_client.get_peer(1).put_packet(msg.to_utf8())
	print("-> ", msg)

func notifyMovement(x:int, y:int, z:int):
	var msg = "notify_movement:"+String(x)+","+String(z)
	_client.get_peer(1).put_packet(msg.to_utf8())
	print("-> ", msg)

func rotation_to_int(rot:String) -> int:
	match rot:
		"down": return 0
		"right": return 16
		"up": return 10
		"left": return 22
	return 0

func rotation_to_string(rot:int) -> String:
	match rot:
		0: return "down"
		16: return "right"
		10: return "up"
		22: return "left"
	return "down"
