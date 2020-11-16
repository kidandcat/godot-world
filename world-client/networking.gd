extends Node

export var websocket_url = "ws://localhost:5000"
var _client = WebSocketClient.new()

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
	match msg[0]:
		"newmesh":
			var data = msg[1].split(",")
			get_node("/root/World/GridMap").set_cell_item(int(data[1]), 0, int(data[2]), int(data[0]), 0)

func _process(delta):
	_client.poll()

func newMesh(type:int, x:float, z:float):
	var err = _client.get_peer(1).put_packet(("create_mesh:"+String(type)+","+String(x)+","+String(z)).to_utf8())
	print("packet sent ", err)
